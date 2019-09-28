/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"context"

	"github.com/kubenext/lissio/internal/controllers"
	"github.com/kubenext/lissio/pkg/action"
)

const (
	RequestPerformAction = "performAction"
)

// ActionRequestManager manages action requests. Action requests allow a generic interface
// for supporting dynamic requests from clients.
type ActionRequestManager struct {
}

var _ StateManager = (*ActionRequestManager)(nil)

// NewActionRequestManager creates an instance of ActionRequestManager.
func NewActionRequestManager() *ActionRequestManager {
	return &ActionRequestManager{}
}

func (a ActionRequestManager) Start(ctx context.Context, state controllers.State, s LissioClient) {
}

// Handlers returns the handlers this manager supports.
func (a *ActionRequestManager) Handlers() []controllers.ClientRequestHandler {
	return []controllers.ClientRequestHandler{
		{
			RequestType: RequestPerformAction,
			Handler:     a.PerformAction,
		},
	}
}

// PerformAction is a handler than runs an action.
func (a *ActionRequestManager) PerformAction(state controllers.State, payload action.Payload) error {
	ctx := context.TODO()

	actionName, err := payload.String("action")
	if err != nil {
		// TODO: alert the user this action doesn't exist
		return nil
	}

	if err := state.Dispatch(ctx, actionName, payload); err != nil {
		return err
	}

	return nil
}
