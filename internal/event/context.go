/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"sort"
	"time"

	"github.com/pkg/errors"

	"github.com/kubenext/lissio/internal/config"
	"github.com/kubenext/lissio/internal/controllers"
	"github.com/kubenext/lissio/internal/kubeconfig"
)

// kubeContextsResponse is a response for current kube contexts.
type kubeContextsResponse struct {
	Contexts       []kubeconfig.Context `json:"contexts"`
	CurrentContext string               `json:"currentContext"`
}

type ContextGeneratorOption func(generator *ContextsGenerator)

// ContextsGenerator generates kube contexts for the front end.
type ContextsGenerator struct {
	ConfigLoader kubeconfig.Loader
	DashConfig   config.Dash
}

var _ controllers.Generator = (*ContextsGenerator)(nil)

func NewContextsGenerator(dashConfig config.Dash, options ...ContextGeneratorOption) *ContextsGenerator {
	kcg := &ContextsGenerator{
		ConfigLoader: kubeconfig.NewFSLoader(),
		DashConfig:   dashConfig,
	}

	for _, option := range options {
		option(kcg)
	}

	return kcg
}

func (g *ContextsGenerator) Event(ctx context.Context) (controllers.Event, error) {
	configPath := g.DashConfig.KubeConfigPath()

	kubeConfig, err := g.ConfigLoader.Load(configPath)
	if err != nil {
		return controllers.Event{}, errors.Wrap(err, "unable to load kube config")
	}

	currentContext := g.DashConfig.ContextName()
	if currentContext == "" {
		currentContext = kubeConfig.CurrentContext
	}

	resp := kubeContextsResponse{
		CurrentContext: currentContext,
		Contexts:       kubeConfig.Contexts,
	}

	sort.Slice(resp.Contexts, func(i, j int) bool {
		return resp.Contexts[i].Name < resp.Contexts[j].Name
	})

	e := controllers.Event{
		Type: controllers.EventTypeKubeConfig,
		Data: resp,
	}

	return e, nil
}

func (ContextsGenerator) ScheduleDelay() time.Duration {
	return DefaultScheduleDelay
}

func (ContextsGenerator) Name() string {
	return "kubeConfig"
}
