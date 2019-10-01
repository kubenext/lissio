/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package servicemesh

import (
	"context"
	"github.com/kubenext/lissio/pkg/store"
	"github.com/kubenext/lissio/pkg/view/component"
	"github.com/pkg/errors"
)

var (
	serviceColumns = component.NewTableCols("Name", "Health", "Config")
)

const (
	labelApp     = "app"
	labelService = "service"
)

type SummarizerConfig interface {
	ObjectStore() store.Store
}

type Summarizer interface {
	Summarize(ctx context.Context, namespace string, config SummarizerConfig) (*component.Table, error)
}

type summarizer struct {
}

func (s *summarizer) Summarize(ctx context.Context, namespace string, config SummarizerConfig) (*component.Table, error) {
	if config == nil {
		return nil, errors.Errorf("config is nil")
	}

	table := component.NewTable("Services", "services", serviceColumns)

	services, err := listServices(ctx, config.ObjectStore(), namespace)
	if err != nil {
		return nil, errors.Errorf("get list nil")
	}

	for i := range services {
		table.Add(component.TableRow{
			"Name":   component.NewText(services[i].Name),
			"Health": component.NewText("true"),
			"Config": component.NewText("true"),
		})
	}
	return table, nil
}

type service struct {
	Name string
}

func listServices(ctx context.Context, objectStore store.Store, namespace string) ([]service, error) {
	key := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Service",
	}

	services, _, err := objectStore.List(ctx, key)
	if err != nil {
		return nil, err
	}

	list := make([]service, 0)

	for i := range services.Items {
		s := service{Name: services.Items[i].GetName()}
		list = append(list, s)
	}

	return list, nil
}
