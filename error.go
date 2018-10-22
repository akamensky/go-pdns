package pdns

import (
	"errors"
	"fmt"
	"strings"
)

type failure struct {
	error   Error
	Message string `json:"error"`
}

func (f *failure) getError() error {
	var err error

	// Well, this is ridiculous, but since pdns API does not return any codes have to do it like this
	if strings.HasSuffix(f.Message, "already exists") {
		err = ErrZoneAlreadyExists
	}

	// If some error happened,
	// but we could not match any existing error
	// make new one out of message string
	if err == nil && f.Message != "" {
		err = errors.New(f.Message)
	}

	return err
}

type Error string

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code(), e.Message())
}

func (e Error) Code() string {
	return string(e)
}

func (e Error) Message() string {
	switch e {
	case ErrZoneAlreadyExists:
		return "Zone for this domain already exists"
	default:
		return ""
	}
}

const (
	ErrZoneAlreadyExists Error = "ErrZoneAlreadyExists"
)
