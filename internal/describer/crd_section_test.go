/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/pkg/view/component"
)

func Test_crdSectionDescriber(t *testing.T) {
	csd := NewCRDSection("/path", "title")

	d1View := component.NewText("d1")
	d1 := NewStubDescriber("/d1", component.NewList("", []component.Component{d1View}))

	csd.Add("d1", d1)

	ctx := context.Background()

	view1, err := csd.Describe(ctx, "default", Options{})
	require.NoError(t, err)

	expect1 := component.ContentResponse{
		Title: component.TitleFromString("title"),
		Components: []component.Component{
			component.NewList("Custom Resources", []component.Component{d1View}),
		},
	}

	assert.Equal(t, expect1, view1)

	csd.Remove("d1")

	view2, err := csd.Describe(ctx, "default", Options{})
	require.NoError(t, err)

	expect2 := component.ContentResponse{
		Title: component.TitleFromString("title"),
		Components: []component.Component{
			component.NewList("Custom Resources", nil),
		},
	}

	assert.Equal(t, expect2, view2)

	csd.Add("d1", d1)
	require.NoError(t, csd.Reset(ctx))

	view3, err := csd.Describe(ctx, "default", Options{})
	require.NoError(t, err)

	expected3 := component.ContentResponse{
		Title: component.TitleFromString("title"),
		Components: []component.Component{
			component.NewList("Custom Resources", nil),
		},
	}

	assert.Equal(t, expected3, view3)
}
