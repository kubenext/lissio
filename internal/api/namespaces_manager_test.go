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
	clusterFake "github.com/kubenext/lissio/internal/cluster/fake"
	configFake "github.com/kubenext/lissio/internal/config/fake"
	lissioFake "github.com/kubenext/lissio/internal/controllers/fake"
)

func TestNamespacesManager_GenerateNamespaces(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	state := lissioFake.NewMockState(controller)
	lissioClient := fake.NewMockLissioClient(controller)

	namespaces := []string{"default"}

	lissioClient.EXPECT().
		Send(api.CreateNamespacesEvent(namespaces))

	poller := api.NewSingleRunPoller()
	manager := api.NewNamespacesManager(dashConfig,
		api.WithNamespacesGeneratorPoller(poller),
		api.WithNamespacesGenerator(func(ctx context.Context, config api.NamespaceManagerConfig) (strings []string, e error) {
			return namespaces, nil
		}))

	ctx := context.Background()
	manager.Start(ctx, state, lissioClient)
}

func TestNamespacesGenerator(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(controller *gomock.Controller) *configFake.MockDash
		isErr    bool
		expected []string
	}{
		{
			name: "in general",
			setup: func(controller *gomock.Controller) *configFake.MockDash {
				namespaces := []string{"ns-1"}

				namespaceClient := clusterFake.NewMockNamespaceInterface(controller)
				namespaceClient.EXPECT().Names().Return(namespaces, nil)

				clusterClient := clusterFake.NewMockClientInterface(controller)
				clusterClient.EXPECT().NamespaceClient().Return(namespaceClient, nil)

				dashConfig := configFake.NewMockDash(controller)
				dashConfig.EXPECT().ClusterClient().Return(clusterClient)

				return dashConfig
			},
			expected: []string{"ns-1"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			require.NotNil(t, test.setup)
			dashConfig := test.setup(controller)

			ctx := context.Background()
			got, err := api.NamespacesGenerator(ctx, dashConfig)
			require.NoError(t, err)

			require.Equal(t, test.expected, got)
		})
	}
}
