package errorcodes

import (
	"errors"

	"connectrpc.com/connect"
)

var errAlreadyExists = errors.New("already_exists")
var errFailedPrecondition = errors.New("failed_precondition")
var errInvalidArgument = errors.New("invalid_argument")
var errNotFound = errors.New("not_found")
var errPermissionDenied = errors.New("permission_denied")
var errUnauthenticated = errors.New("unauthenticated")

func NewAlreadyExistsError(sourceError error) error {
	detail := &connect.ErrorDetail{}

	err := connect.NewError(connect.CodeAlreadyExists, errAlreadyExists)
	err.AddDetail(detail)

	return err
}

func NewFailedPreconditionError(sourceError error) error {
	detail := &connect.ErrorDetail{}

	err := connect.NewError(connect.CodeFailedPrecondition, errFailedPrecondition)
	err.AddDetail(detail)

	return err
}

func NewInvalidArgumentError(sourceError error) error {
	detail := &connect.ErrorDetail{}

	err := connect.NewError(connect.CodeInvalidArgument, errInvalidArgument)
	err.AddDetail(detail)

	return err
}

func NewNotFoundError(sourceError error) error {
	detail := &connect.ErrorDetail{}

	err := connect.NewError(connect.CodeNotFound, errNotFound)
	err.AddDetail(detail)

	return err
}

func NewPermissionDeniedError(sourceError error) error {
	detail := &connect.ErrorDetail{}

	err := connect.NewError(connect.CodePermissionDenied, errPermissionDenied)
	err.AddDetail(detail)

	return err
}

func NewUnauthenticatedError(sourceError error) error {
	detail := &connect.ErrorDetail{}

	err := connect.NewError(connect.CodeUnauthenticated, errUnauthenticated)
	err.AddDetail(detail)

	return err
}
