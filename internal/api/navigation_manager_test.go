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

	"github.com/vmware/octant/internal/api"
	"github.com/vmware/octant/internal/api/fake"
	configFake "github.com/vmware/octant/internal/config/fake"
	"github.com/vmware/octant/internal/module"
	moduleFake "github.com/vmware/octant/internal/module/fake"
	"github.com/vmware/octant/internal/octant"
	octantFake "github.com/vmware/octant/internal/octant/fake"
	"github.com/vmware/octant/pkg/navigation"
)

func TestNavigationManager_GenerateNavigation(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	state := octantFake.NewMockState(controller)
	state.EXPECT().GetContentPath().Return("/path")

	octantClient := fake.NewMockOctantClient(controller)

	sections := []navigation.Navigation{{Title: "module"}}

	octantClient.EXPECT().
		Send(api.CreateNavigationEvent(sections, "/path"))

	poller := api.NewSingleRunPoller()
	manager := api.NewNavigationManager(dashConfig,
		api.WithNavigationGeneratorPoller(poller),
		api.WithNavigationGenerator(func(ctx context.Context, state octant.State, config api.NavigationManagerConfig) ([]navigation.Navigation, error) {
			return sections, nil
		}),
	)

	ctx := context.Background()
	manager.Start(ctx, state, octantClient)
}

func TestNavigationGenerator(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(controller *gomock.Controller) (*configFake.MockDash, *octantFake.MockState)
		isErr    bool
		expected []navigation.Navigation
	}{
		{
			name: "in general",
			setup: func(controller *gomock.Controller) (*configFake.MockDash, *octantFake.MockState) {
				m := moduleFake.NewMockModule(controller)
				m.EXPECT().ContentPath().Return("/module")
				m.EXPECT().Name().Return("module").AnyTimes()
				m.EXPECT().
					Navigation(gomock.Any(), "default", "/module").
					Return([]navigation.Navigation{
						{Title: "module"},
					}, nil)

				moduleManager := moduleFake.NewMockManagerInterface(controller)
				moduleManager.EXPECT().Modules().Return([]module.Module{m})

				dashConfig := configFake.NewMockDash(controller)
				dashConfig.EXPECT().ModuleManager().Return(moduleManager)

				state := octantFake.NewMockState(controller)
				state.EXPECT().GetNamespace().Return("default")

				return dashConfig, state
			},
			expected: []navigation.Navigation{
				{Title: "module"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			require.NotNil(t, test.setup)
			dashConfig, state := test.setup(controller)

			ctx := context.Background()
			got, err := api.NavigationGenerator(ctx, state, dashConfig)
			require.NoError(t, err)

			require.Equal(t, test.expected, got)
		})
	}
}
