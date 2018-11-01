package pdns

import (
	"errors"
	"fmt"
)

type failure struct {
	error   Error
	Code    string `json:"code"`
	Message string `json:"error"`
}

func (f *failure) getError() error {
	var err error

	switch f.Code {
	case "ERR_ZONE_ALREADY_EXISTS":
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
