/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kubenext/lissio/pkg/action"
	"time"
)

type ActionError struct {
	id          string
	timestamp   time.Time
	payload     action.Payload
	requestType string
	err         error
}

func NewActionError(requestType string, payload action.Payload, err error) *ActionError {
	id, _ := uuid.NewUUID()
	return &ActionError{
		id:          id.String(),
		timestamp:   time.Now(),
		payload:     payload,
		requestType: requestType,
		err:         err,
	}
}

func (a *ActionError) ID() string {
	return a.id
}

func (a *ActionError) Error() string {
	return fmt.Sprintf("%s: %s", a.requestType, a.err)
}

func (a *ActionError) Timestamp() time.Time {
	return a.timestamp
}

var _ InternalError = (*ActionError)(nil)
