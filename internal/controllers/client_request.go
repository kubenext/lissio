/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package controllers

import (
	"github.com/kubenext/lissio/pkg/action"
)

// ClientRequestHandler is a client request.
type ClientRequestHandler struct {
	RequestType string
	Handler     func(state State, payload action.Payload) error
}
