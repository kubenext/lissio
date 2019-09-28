/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"github.com/vmware/octant/pkg/action"
)

// ClientRequestHandler is a client request.
type ClientRequestHandler struct {
	RequestType string
	Handler     func(state State, payload action.Payload) error
}
