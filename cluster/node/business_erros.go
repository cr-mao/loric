package node

import (
	"github.com/cr-mao/loric/errors"
)

var (
	ErrInvalidArgument    = errors.New("ErrInvalidArgument")
	ErrNotFoundUserSource = errors.New("not found user source")
	ErrInvalidGID         = errors.New("invalid gate id")
	ErrInvalidSessionKind = errors.New("invalid session kind")
	ErrInvalidMessage     = errors.New("invalid message")
	ErrInvalidNID         = errors.New("invalid node id")
	ErrReceiveTargetEmpty = errors.New("the receive target is empty")
)
