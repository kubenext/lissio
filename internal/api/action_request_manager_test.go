/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kubenext/lissio/internal/api"
	octantFake "github.com/kubenext/lissio/internal/controllers/fake"
	"github.com/kubenext/lissio/pkg/action"
)

func TestActionRequestManager_Handlers(t *testing.T) {
	manager := api.NewActionRequestManager()
	AssertHandlers(t, manager, []string{api.RequestPerformAction})
}

func TestActionRequestManager_PerformAction(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	state := octantFake.NewMockState(controller)

	manager := api.NewActionRequestManager()

	payload := action.CreatePayload(api.RequestPerformAction, map[string]interface{}{
		"foo": "bar",
	})

	state.EXPECT().
		Dispatch(gomock.Any(), api.RequestPerformAction, payload).
		Return(nil)

	require.NoError(t, manager.PerformAction(state, payload))
}
