/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/vmware/octant/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Annotations_Marshal(t *testing.T) {
	cases := []struct {
		name         string
		input        *component.Annotations
		expectedPath string
	}{
		{
			name: "in general",
			input: component.NewAnnotations(map[string]string{
				"foo": "bar",
			}),
			expectedPath: "annotations.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.input)
			require.NoError(t, err)

			expected, err := ioutil.ReadFile(filepath.Join("testdata", tc.expectedPath))
			require.NoError(t, err)

			assert.JSONEq(t, string(expected), string(got))
		})
	}
}

func Test_Annotations_IsEmpty(t *testing.T) {
	cases := []struct {
		name    string
		input   *component.Annotations
		isEmpty bool
	}{
		{
			name:    "empty (nil)",
			input:   component.NewAnnotations(nil),
			isEmpty: true,
		},
		{
			name:    "empty",
			input:   component.NewAnnotations(map[string]string{}),
			isEmpty: true,
		},
		{
			name: "not empty",
			input: component.NewAnnotations(map[string]string{
				"foo": "bar",
			}),
			isEmpty: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isEmpty, tc.input.IsEmpty())
		})
	}
}

func Test_Annotations_GetMetadata(t *testing.T) {
	input := component.NewAnnotations(map[string]string{
		"foo": "bar",
	})

	assert.Equal(t, "annotations", input.GetMetadata().Type)
}
