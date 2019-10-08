package errors

import (
	"time"
)

type InternalError interface {
	ID() string
	Error() string
	Timestamp() time.Time
}
