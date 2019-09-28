/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package module_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"

	clusterfake "github.com/vmware/octant/internal/cluster/fake"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/module/fake"
)

func TestManager(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	clusterClient := clusterfake.NewMockClientInterface(controller)

	actionRegistrar := fake.NewMockActionRegistrar(controller)

	manager, err := module.NewManager(clusterClient, "default", actionRegistrar, log.NopLogger())
	require.NoError(t, err)

	modules := manager.Modules()
	require.NoError(t, err)
	require.Len(t, modules, 0)

	m := fake.NewMockModule(controller)
	m.EXPECT().Start().Return(nil)
	m.EXPECT().Stop()
	m.EXPECT().SetNamespace("other").Return(nil)

	require.NoError(t, manager.Register(m))

	modules = manager.Modules()
	require.NoError(t, err)
	require.Len(t, modules, 1)

	manager.SetNamespace("other")
	manager.Unload()
}

func TestManager_ObjectPath(t *testing.T) {
	cases := []struct {
		name       string
		apiVersion string
		kind       string
		isErr      bool
		expected   string
	}{
		{
			name:       "exists",
			apiVersion: "group/version",
			kind:       "kind",
			expected:   "/path",
		},
		{
			name:       "does not exist",
			apiVersion: "v1",
			kind:       "kind",
			expected:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			clusterClient := clusterfake.NewMockClientInterface(controller)
			actionRegistrar := fake.NewMockActionRegistrar(controller)

			manager, err := module.NewManager(clusterClient, "default", actionRegistrar, log.NopLogger())
			require.NoError(t, err)

			modules := manager.Modules()
			require.NoError(t, err)
			require.Len(t, modules, 0)

			m := fake.NewMockModule(controller)
			m.EXPECT().Start().Return(nil)
			supportedGVK := []schema.GroupVersionKind{
				{Group: "group", Version: "version", Kind: "kind"},
			}
			m.EXPECT().SupportedGroupVersionKind().Return(supportedGVK)
			m.EXPECT().
				GroupVersionKindPath("namespace", "group/version", "kind", "name").
				Return("/path", nil).
				AnyTimes()

			require.NoError(t, manager.Register(m))

			got, err := manager.ObjectPath("namespace", tc.apiVersion, tc.kind, "name")

			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}
