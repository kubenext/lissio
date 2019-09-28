/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/kubenext/lissio/internal/api"
	"github.com/kubenext/lissio/internal/api/fake"
	configFake "github.com/kubenext/lissio/internal/config/fake"
	"github.com/kubenext/lissio/internal/controllers"
	lissioFake "github.com/kubenext/lissio/internal/controllers/fake"
	"github.com/kubenext/lissio/internal/log"
)

func TestContextManager_Handlers(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	manager := api.NewContextManager(dashConfig)
	AssertHandlers(t, manager, []string{api.RequestSetContext})
}

func TestContext_GenerateContexts(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	state := lissioFake.NewMockState(controller)
	lissioClient := fake.NewMockLissioClient(controller)

	ev := controllers.Event{
		Type: "eventType",
	}
	lissioClient.EXPECT().Send(ev)

	logger := log.NopLogger()

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().Logger().Return(logger).AnyTimes()

	poller := api.NewSingleRunPoller()
	generatorFunc := func(ctx context.Context, state controllers.State) (controllers.Event, error) {
		return ev, nil
	}
	manager := api.NewContextManager(dashConfig,
		api.WithContextGenerator(generatorFunc),
		api.WithContextGeneratorPoll(poller))

	ctx := context.Background()
	manager.Start(ctx, state, lissioClient)
}
