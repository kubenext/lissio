/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kubenext/lissio/internal/api"
	"github.com/kubenext/lissio/internal/api/fake"
	"github.com/kubenext/lissio/internal/controllers"
	lissioFake "github.com/kubenext/lissio/internal/controllers/fake"
	"github.com/kubenext/lissio/internal/log"
	moduleFake "github.com/kubenext/lissio/internal/module/fake"
	"github.com/kubenext/lissio/pkg/action"
	"github.com/kubenext/lissio/pkg/view/component"
)

func TestContentManager_Handlers(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	moduleManager := moduleFake.NewMockManagerInterface(controller)

	logger := log.NopLogger()

	manager := api.NewContentManager(moduleManager, logger)
	AssertHandlers(t, manager, []string{
		api.RequestSetContentPath,
		api.RequestSetNamespace,
	})
}

func TestContentManager_GenerateContent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	params := map[string][]string{}

	moduleManager := moduleFake.NewMockManagerInterface(controller)
	state := lissioFake.NewMockState(controller)

	state.EXPECT().GetContentPath().Return("/path")
	state.EXPECT().GetNamespace().Return("default")
	state.EXPECT().GetQueryParams().Return(params)
	state.EXPECT().OnContentPathUpdate(gomock.Any()).DoAndReturn(func(fn controllers.ContentPathUpdateFunc) controllers.UpdateCancelFunc {
		fn("foo")
		return func() {}
	})
	lissioClient := fake.NewMockOctantClient(controller)

	contentResponse := component.ContentResponse{
		IconName: "fake",
	}
	contentEvent := api.CreateContentEvent(contentResponse, "default", "/path", params)
	lissioClient.EXPECT().Send(contentEvent).AnyTimes()

	logger := log.NopLogger()

	poller := api.NewSingleRunPoller()

	contentGenerator := func(ctx context.Context, state controllers.State) (component.ContentResponse, bool, error) {
		return contentResponse, false, nil
	}
	manager := api.NewContentManager(moduleManager, logger,
		api.WithContentGenerator(contentGenerator),
		api.WithContentGeneratorPoller(poller))

	ctx := context.Background()
	manager.Start(ctx, state, lissioClient)
}

func TestContentManager_SetContentPath(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	m := moduleFake.NewMockModule(controller)
	m.EXPECT().Name().Return("name").AnyTimes()

	moduleManager := moduleFake.NewMockManagerInterface(controller)

	state := lissioFake.NewMockState(controller)
	state.EXPECT().SetContentPath("/path")

	logger := log.NopLogger()

	manager := api.NewContentManager(moduleManager, logger,
		api.WithContentGeneratorPoller(api.NewSingleRunPoller()))

	payload := action.Payload{
		"contentPath": "/path",
	}

	require.NoError(t, manager.SetContentPath(state, payload))
}

func TestContentManager_SetNamespace(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	m := moduleFake.NewMockModule(controller)
	m.EXPECT().Name().Return("name").AnyTimes()

	moduleManager := moduleFake.NewMockManagerInterface(controller)

	state := lissioFake.NewMockState(controller)
	state.EXPECT().SetNamespace("kube-system")

	logger := log.NopLogger()

	manager := api.NewContentManager(moduleManager, logger,
		api.WithContentGeneratorPoller(api.NewSingleRunPoller()))

	payload := action.Payload{
		"namespace": "kube-system",
	}

	require.NoError(t, manager.SetNamespace(state, payload))
}

func TestContentManager_SetQueryParams(t *testing.T) {
	tests := []struct {
		name    string
		payload action.Payload
		setup   func(state *lissioFake.MockState)
	}{
		{
			name: "single filter",
			payload: action.Payload{
				"params": map[string]interface{}{
					"filters": "foo:bar",
				},
			},
			setup: func(state *lissioFake.MockState) {
				state.EXPECT().SetFilters([]controllers.Filter{
					{Key: "foo", Value: "bar"},
				})
			},
		},
		{
			name: "multiple filters",
			payload: action.Payload{
				"params": map[string]interface{}{
					"filters": []interface{}{
						"foo:bar",
						"baz:qux",
					},
				},
			},
			setup: func(state *lissioFake.MockState) {
				state.EXPECT().SetFilters([]controllers.Filter{
					{Key: "foo", Value: "bar"},
					{Key: "baz", Value: "qux"},
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			m := moduleFake.NewMockModule(controller)
			m.EXPECT().Name().Return("name").AnyTimes()

			moduleManager := moduleFake.NewMockManagerInterface(controller)

			state := lissioFake.NewMockState(controller)
			require.NotNil(t, test.setup)
			test.setup(state)

			logger := log.NopLogger()

			manager := api.NewContentManager(moduleManager, logger,
				api.WithContentGeneratorPoller(api.NewSingleRunPoller()))
			require.NoError(t, manager.SetQueryParams(state, test.payload))
		})
	}
}
