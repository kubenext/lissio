/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/kubenext/lissio/internal/config"
	"github.com/kubenext/lissio/internal/controllers"
	"github.com/kubenext/lissio/pkg/action"
)

//go:generate mockgen -destination=./fake/mock_state_manager.go -package=fake github.com/kubenext/lissio/internal/api StateManager
//go:generate mockgen -destination=./fake/mock_lissio_client.go -package=fake github.com/kubenext/lissio/internal/api LissioClient

var (
	reContentPathNamespace = regexp.MustCompile(`^/namespace/(?P<namespace>[^/]+)/?`)
)

// StateManager manages states for WebsocketState.
type StateManager interface {
	Handlers() []controllers.ClientRequestHandler
	Start(ctx context.Context, state controllers.State, s LissioClient)
}

func defaultStateManagers(clientID string, dashConfig config.Dash) []StateManager {
	logger := dashConfig.Logger().With("client-id", clientID)

	return []StateManager{
		NewContentManager(dashConfig.ModuleManager(), logger),
		NewFilterManager(),
		NewNavigationManager(dashConfig),
		NewNamespacesManager(dashConfig),
		NewContextManager(dashConfig),
		NewActionRequestManager(),
	}
}

// LissioClient is an LissioClient.
type LissioClient interface {
	Send(event controllers.Event)
	ID() string
}

type atomicString struct {
	mu sync.RWMutex
	s  string
}

func newStringValue(initial string) *atomicString {
	return &atomicString{
		mu: sync.RWMutex{},
		s:  initial,
	}
}

func (s *atomicString) get() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.s
}

func (s *atomicString) set(v string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.s = v
}

// WebsocketStateOption is an option for configuring WebsocketState.
type WebsocketStateOption func(w *WebsocketState)

// WebsocketStateManagers configures WebsocketState's state managers.
func WebsocketStateManagers(managers []StateManager) WebsocketStateOption {
	return func(w *WebsocketState) {
		w.managers = managers
	}
}

// WebsocketState manages state for a websocket client.
type WebsocketState struct {
	dashConfig         config.Dash
	wsClient           LissioClient
	contentPath        *atomicString
	namespace          *atomicString
	filters            []controllers.Filter
	contentPathUpdates map[string]controllers.ContentPathUpdateFunc
	namespaceUpdates   map[string]controllers.NamespaceUpdateFunc

	mu               sync.RWMutex
	managers         []StateManager
	actionDispatcher ActionDispatcher

	startCtx           context.Context
	managersCancelFunc context.CancelFunc
}

var _ controllers.State = (*WebsocketState)(nil)

// NewWebsocketState creates an instance of WebsocketState.
func NewWebsocketState(dashConfig config.Dash, actionDispatcher ActionDispatcher, wsClient LissioClient, options ...WebsocketStateOption) *WebsocketState {
	defaultNamespace := dashConfig.DefaultNamespace()

	w := &WebsocketState{
		dashConfig:         dashConfig,
		wsClient:           wsClient,
		contentPathUpdates: make(map[string]controllers.ContentPathUpdateFunc),
		namespaceUpdates:   make(map[string]controllers.NamespaceUpdateFunc),
		namespace:          newStringValue(defaultNamespace),
		contentPath:        newStringValue(""),
		filters:            make([]controllers.Filter, 0),
		actionDispatcher:   actionDispatcher,
	}

	for _, option := range options {
		option(w)
	}

	if len(w.managers) < 1 {
		w.managers = defaultStateManagers(wsClient.ID(), dashConfig)
	}

	return w
}

// Start starts WebsocketState by starting all associated StateManagers.
func (c *WebsocketState) Start(ctx context.Context) {
	for i := range c.managers {
		go c.managers[i].Start(ctx, c, c.wsClient)
	}
}

// Handlers returns all the handlers for WebsocketState.
func (c *WebsocketState) Handlers() []controllers.ClientRequestHandler {
	var handlers []controllers.ClientRequestHandler

	for _, manager := range c.managers {
		handlers = append(handlers, manager.Handlers()...)
	}

	return handlers
}

// Dispatch dispatches a message.
func (c *WebsocketState) Dispatch(ctx context.Context, actionName string, payload action.Payload) error {
	return c.actionDispatcher.Dispatch(ctx, c, actionName, payload)
}

// SetContentPath sets the content path.
func (c *WebsocketState) SetContentPath(contentPath string) {
	if contentPath == "" {
		contentPath = path.Join("overview", "namespace", c.namespace.get())
	} else if c.contentPath.get() == contentPath {
		return
	}

	c.dashConfig.Logger().With(
		"contentPath", contentPath).
		Debugf("setting content path")

	c.contentPath.set(contentPath)

	m, ok := c.dashConfig.ModuleManager().ModuleForContentPath(contentPath)
	if !ok {
		c.dashConfig.Logger().
			With("contentPath", contentPath).
			Warnf("unable to find module for content path")
	} else {
		modulePath := strings.TrimPrefix(contentPath, m.Name())
		match := reContentPathNamespace.FindStringSubmatch(modulePath)
		result := make(map[string]string)
		if len(match) > 0 {
			for i, name := range reContentPathNamespace.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}

			if result["namespace"] != "" {
				c.SetNamespace(result["namespace"])
			}
		}
	}

	for _, fn := range c.contentPathUpdates {
		fn(contentPath)
	}

}

// GetContentPath returns the content path.
func (c *WebsocketState) GetContentPath() string {
	return c.contentPath.get()
}

// OnContentPathUpdate registers a function that will be called when the content path changes.
func (c *WebsocketState) OnContentPathUpdate(fn controllers.ContentPathUpdateFunc) controllers.UpdateCancelFunc {
	c.mu.Lock()
	defer c.mu.Unlock()

	id, _ := uuid.NewUUID()
	c.contentPathUpdates[id.String()] = fn

	cancelFunc := func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.contentPathUpdates, id.String())
	}

	return cancelFunc
}

// SetNamespace sets the namespace.
func (c *WebsocketState) SetNamespace(namespace string) {
	cur := c.namespace.get()
	if namespace == cur {
		return
	}

	c.dashConfig.Logger().
		With("namespace", namespace).
		Debugf("setting namespace")
	c.namespace.set(namespace)

	newPath := updateContentPathNamespace(c.contentPath.get(), namespace)
	if newPath != c.contentPath.get() {
		c.SetContentPath(newPath)
	}

	for _, fn := range c.namespaceUpdates {
		fn(namespace)
	}
}

// GetNamespace gets the namespace.
func (c *WebsocketState) GetNamespace() string {
	return c.namespace.get()
}

// OnNamespaceUpdate registers a function that will be run when the namespace changes.
func (c *WebsocketState) OnNamespaceUpdate(fn controllers.NamespaceUpdateFunc) controllers.UpdateCancelFunc {
	c.mu.Lock()
	defer c.mu.Unlock()

	id, _ := uuid.NewUUID()
	c.namespaceUpdates[id.String()] = fn

	cancelFunc := func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.namespaceUpdates, id.String())
	}

	return cancelFunc
}

// AddFilter adds a content filter.
func (c *WebsocketState) AddFilter(filter controllers.Filter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.filters {
		if c.filters[i].IsEqual(filter) {
			return
		}
	}

	c.filters = append(c.filters, filter)
}

// RemoveFilter removes a content filter.
func (c *WebsocketState) RemoveFilter(filter controllers.Filter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var newFilters []controllers.Filter

	for i := range c.filters {
		if c.filters[i].IsEqual(filter) {
			continue
		}
		newFilters = append(newFilters, c.filters[i])
	}

	c.filters = newFilters
}

// GetFilters returns all filters.
func (c *WebsocketState) GetFilters() []controllers.Filter {
	filters := make([]controllers.Filter, len(c.filters))
	copy(filters, c.filters)

	sort.Slice(filters, func(i, j int) bool {
		return filters[i].Key < filters[j].Key
	})

	return filters
}

func (c *WebsocketState) SetFilters(filters []controllers.Filter) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.filters = filters
}

// SetContext sets the Kubernetes context.
func (c *WebsocketState) SetContext(requestedContext string) {
	if err := c.dashConfig.UseContext(context.TODO(), requestedContext); err != nil {
		c.dashConfig.Logger().WithErr(err).Errorf("update context")
	}

	c.SetNamespace(c.dashConfig.DefaultNamespace())

	for _, fn := range c.contentPathUpdates {
		fn(c.GetContentPath())
	}

	c.wsClient.Send(CreateAlertUpdate(action.CreateAlert(
		action.AlertTypeInfo,
		fmt.Sprintf("Changing context to %s", requestedContext),
		action.DefaultAlertExpiration,
	)))
}

func (c *WebsocketState) GetQueryParams() map[string][]string {
	filters := c.filters

	c.wsClient.Send(CreateFiltersUpdate(filters))

	queryParams := map[string][]string{}

	var filterList []string
	for _, filter := range filters {
		filterList = append(filterList, filter.ToQueryParam())
	}
	if len(filterList) > 0 {
		queryParams["filters"] = filterList
	}

	return queryParams
}

// SendAlert sends an alert to the websocket client.
func (c *WebsocketState) SendAlert(alert action.Alert) {
	c.wsClient.Send(CreateAlertUpdate(alert))
}

func updateContentPathNamespace(in, namespace string) string {
	parts := strings.Split(in, "/")
	if in == "" {
		return ""
	}

	if len(parts) > 1 && parts[1] == "namespace" {
		parts[2] = namespace
		return path.Join(parts...)
	}
	return in
}

// CreateFiltersUpdate creates a filters update event.
func CreateFiltersUpdate(filters []controllers.Filter) controllers.Event {
	if filters == nil {
		filters = make([]controllers.Filter, 0)
	}
	return CreateEvent("filters", action.Payload{
		"filters": filters,
	})
}

// CreateAlertUpdate creates an alert update event.
func CreateAlertUpdate(alert action.Alert) controllers.Event {
	return CreateEvent(controllers.EventTypeAlert, action.Payload{
		"type":       alert.Type,
		"message":    alert.Message,
		"expiration": alert.Expiration,
	})
}
