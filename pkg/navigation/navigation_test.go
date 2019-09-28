/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package navigation

import (
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/icon"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/store/fake"
)

func Test_NewNavigation(t *testing.T) {
	navPath := "/navPath"
	title := "title"

	nav, err := New(title, navPath)
	require.NoError(t, err)

	assert.Equal(t, navPath, nav.Path)
	assert.Equal(t, title, nav.Title)
}

func TestEntriesHelper(t *testing.T) {
	neh := EntriesHelper{}

	neh.Add("title", "suffix", icon.OverviewService, false)

	list, err := neh.Generate("/prefix")
	require.NoError(t, err)

	expected := Navigation{
		Title:    "title",
		Path:     path.Join("/prefix", "suffix"),
		IconName: fmt.Sprintf("internal:%s", icon.OverviewService),
	}

	assert.Len(t, list, 1)
	assert.Equal(t, expected.Title, list[0].Title)
	assert.Equal(t, expected.Path, list[0].Path)
	assert.Equal(t, expected.IconName, list[0].IconName)
	assert.NotEmpty(t, list[0].IconSource)
}

func TestCRDEntries_namespace_scoped(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	objectStore := fake.NewMockStore(controller)
	clusterScopedCRD := createCRD("cluster-scoped", "ClusterScoped", true)
	namespaceScopedCRD := createCRD("namespace-scoped", "NamespaceScoped", false)

	crds := testutil.ToUnstructuredList(t, clusterScopedCRD, namespaceScopedCRD)
	crdKey := store.Key{
		APIVersion: "apiextensions.k8s.io/v1beta1",
		Kind:       "CustomResourceDefinition",
	}
	objectStore.EXPECT().
		List(gomock.Any(), crdKey).
		Return(crds, false, nil).
		AnyTimes()

	clusterCR := createCR("testing", "v1", "ClusterScoped", "cluster-scoped")
	clusterCRs := testutil.ToUnstructuredList(t, clusterCR)
	namespaceCR := createCR("testing", "v1", "NamespaceScoped", "namespace-scoped")
	namespaceCRs := testutil.ToUnstructuredList(t, namespaceCR)

	crNamespaceKey := store.Key{
		Namespace:  "default",
		APIVersion: "testing/v1",
		Kind:       "NamespaceScoped",
	}
	objectStore.EXPECT().
		List(gomock.Any(), crNamespaceKey).
		Return(namespaceCRs, false, nil).
		AnyTimes()
	crClusterKey := store.Key{
		APIVersion: "testing/v1",
		Kind:       "ClusterScoped",
	}
	objectStore.EXPECT().
		List(gomock.Any(), crClusterKey).
		Return(clusterCRs, false, nil).
		AnyTimes()

	ctx := context.Background()

	namespaceGot, _, err := CRDEntries(ctx, "/prefix", "default", objectStore, false)
	require.NoError(t, err)

	namespaceExpected := []Navigation{
		createNavForCR(t, namespaceCR.GetName()),
	}

	assert.Equal(t, namespaceExpected, namespaceGot)

	clusterGot, _, err := CRDEntries(ctx, "/prefix", "default", objectStore, true)
	require.NoError(t, err)

	clusterExpected := []Navigation{
		createNavForCR(t, clusterCR.GetName()),
	}

	assert.Equal(t, clusterExpected, clusterGot)
}

func createNavForCR(t *testing.T, name string) Navigation {
	nav, err := New(name, path.Join("/prefix", name), SetNavigationIcon(icon.CustomResourceDefinition))
	require.NoError(t, err)

	return *nav
}

func createCRD(name, kind string, isClusterScoped bool) *apiextv1beta1.CustomResourceDefinition {
	scope := apiextv1beta1.ClusterScoped
	if !isClusterScoped {
		scope = apiextv1beta1.NamespaceScoped
	}

	crd := testutil.CreateCRD(name)
	crd.Spec.Scope = scope
	crd.Spec.Group = "testing"
	crd.Spec.Names = apiextv1beta1.CustomResourceDefinitionNames{
		Kind: kind,
	}
	// TODO fix this because Version is deprecated
	crd.Spec.Version = "v1"
	crd.Spec.Versions = []apiextv1beta1.CustomResourceDefinitionVersion{
		{
			Name:    "v1",
			Served:  true,
			Storage: true,
		},
	}

	return crd
}

func createCR(group, version, kind, name string) *unstructured.Unstructured {
	m := make(map[string]interface{})
	u := &unstructured.Unstructured{Object: m}

	u.SetName(name)
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	})

	return u
}
