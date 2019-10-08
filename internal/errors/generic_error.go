/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type GenericError struct {
	id        string
	timestamp time.Time
	err       error
}

func NewGenericError(err error) *GenericError {
	id, _ := uuid.NewUUID()
	return &GenericError{
		id:        id.String(),
		timestamp: time.Now(),
		err:       err,
	}
}

func (g *GenericError) ID() string {
	return g.id
}

func (g *GenericError) Error() string {
	return fmt.Sprintf("%s", g.err)
}

func (g *GenericError) Timestamp() time.Time {
	return g.timestamp
}

var _ InternalError = (*GenericError)(nil)
