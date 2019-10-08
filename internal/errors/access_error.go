/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"fmt"
	"github.com/kubenext/lissio/pkg/store"
	"time"
)

type AccessError struct {
	id        string
	key       store.Key
	timestamp time.Time
	verb      string
	err       error
}

var _ InternalError = (*AccessError)(nil)

func NewAccessError(key store.Key, verb string, err error) *AccessError {
	return &AccessError{
		id:        key.String(),
		key:       key,
		timestamp: time.Now(),
		verb:      verb,
		err:       err,
	}
}

// ID returns the error unique ID.
func (o *AccessError) ID() string {
	return o.id
}

func (o *AccessError) Timestamp() time.Time {
	return o.timestamp
}

func (o *AccessError) Error() string {
	return fmt.Sprintf("%s: %s: %s", o.verb, o.key, o.err)
}

func (o *AccessError) Key() store.Key {
	return o.key
}

func (o *AccessError) Verb() string {
	return o.verb
}
