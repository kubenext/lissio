/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_RoleListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	now := testutil.Time()

	role := testutil.CreateRole("pod-reader")
	role.CreationTimestamp = metav1.Time{Time: now}
	roleList := &rbacv1.RoleList{
		Items: []rbacv1.Role{
			*role,
		},
	}

	tpo.PathForObject(role, role.Name, "/role")

	ctx := context.Background()
	observed, err := RoleListHandler(ctx, roleList, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Age")
	expected := component.NewTable("Roles", "We couldn't find any roles!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", role.Name, "/role"),
		"Age":  component.NewTimestamp(role.CreationTimestamp.Time),
	})

	component.AssertEqual(t, expected, observed)
}

func Test_RoleConfiguration(t *testing.T) {
	role := testutil.CreateRole("role")

	cases := []struct {
		name     string
		role     *rbacv1.Role
		isErr    bool
		expected *component.Summary
	}{
		{
			name: "general",
			role: role,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Name",
					Content: component.NewText("role"),
				},
			}...),
		},
		{
			name:  "role is nil",
			role:  nil,
			isErr: true,
		},
	}

	for _, tc := range cases {
		controller := gomock.NewController(t)
		defer controller.Finish()

		tpo := newTestPrinterOptions(controller)
		printOptions := tpo.ToOptions()

		rc := NewRoleConfiguration(tc.role)

		summary, err := rc.Create(printOptions)
		if tc.isErr {
			require.Error(t, err)
			return
		}
		require.NoError(t, err)

		component.AssertEqual(t, tc.expected, summary)
	}
}

func Test_createRolePolicyRulesView(t *testing.T) {
	role := testutil.CreateRole("role")
	// TODO: (GuessWhoSamFoo) Test more complex rules
	role.Rules = []rbacv1.PolicyRule{
		{
			Resources:       []string{""},
			NonResourceURLs: []string{"/healthz"},
			ResourceNames:   []string{""},
			Verbs:           []string{"update"},
		},
	}

	got, err := createRolePolicyRulesView(role)
	require.NoError(t, err)

	cols := component.NewTableCols("Resources", "Non-Resource URLs", "Resource Names", "Verbs")
	expected := component.NewTable("Policy Rules", "There are no policy rules!", cols)
	expected.Add([]component.TableRow{
		{
			"Resources":         component.NewText(""),
			"Non-Resource URLs": component.NewText("['/healthz']"),
			"Resource Names":    component.NewText(""),
			"Verbs":             component.NewText("['update']"),
		},
	}...)

	component.AssertEqual(t, expected, got)
}
