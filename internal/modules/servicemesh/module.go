/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package servicemesh

import (
	"context"
	"github.com/kubenext/lissio/internal/config"
	"github.com/kubenext/lissio/internal/controllers"
	"github.com/kubenext/lissio/internal/describer"
	"github.com/kubenext/lissio/internal/generator"
	"github.com/kubenext/lissio/internal/module"
	"github.com/kubenext/lissio/pkg/navigation"
	"github.com/kubenext/lissio/pkg/view/component"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"path/filepath"
)

// Options are options for configuring Module.
type Options struct {
	DashConfig config.Dash
}

type Module struct {
	Options
	pathMatcher *describer.PathMatcher
}

var _ module.Module = (*Module)(nil)

// New creates an instance of module.
func New(ctx context.Context, options Options) *Module {
	pm := describer.NewPathMatcher("servicemesh")
	//for _, pf := range rootDescriber.PathFilters() {
	//	pm.Register(ctx, pf)
	//}
	return &Module{
		Options:     options,
		pathMatcher: pm,
	}
}

// Name is the name of the module.
func (m *Module) Name() string {
	return "servicemesh"
}

// ClientRequestHandlers are client handlers for the module.
func (m *Module) ClientRequestHandlers() []controllers.ClientRequestHandler {
	return nil
}

// Content generates content for a content path.
func (m *Module) Content(ctx context.Context, contentPath string, opts module.ContentOptions) (component.ContentResponse, error) {
	g, err := generator.NewGenerator(m.pathMatcher, m.DashConfig)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	return g.Generate(ctx, contentPath, generator.Options{})
}

func (m *Module) ContentPath() string {
	return m.Name()
}

func (m *Module) Navigation(ctx context.Context, namespace, root string) ([]navigation.Navigation, error) {
	rootPath := filepath.Join(m.ContentPath(), "namespace", namespace)

	rootNav := navigation.Navigation{
		Title: "Service Mesh",
		Path:  rootPath,
	}

	return []navigation.Navigation{rootNav}, nil
}

// SetNamespace sets the module's namespace.
func (m *Module) SetNamespace(namespace string) error {
	return nil
}

func (m *Module) Start() error {
	return nil
}

func (m *Module) Stop() {

}

func (m *Module) SetContext(ctx context.Context, contextName string) error {
	return nil
}

// Generators does nothing.
func (m *Module) Generators() []controllers.Generator {
	return nil
}

func (m *Module) SupportedGroupVersionKind() []schema.GroupVersionKind {
	return nil
}

func (m *Module) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	return "", errors.Errorf("not supported")
}

func (m *Module) AddCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

func (m *Module) RemoveCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

func (m *Module) ResetCRDs(ctx context.Context) error {
	return nil
}
