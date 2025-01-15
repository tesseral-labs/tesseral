package errorcodes

import "errors"

var errAlreadyExists = errors.New("already_exists")
var errFailedPrecondition = errors.New("failed_precondition")
var errInvalidArgument = errors.New("invalid_argument")
var errNotFound = errors.New("not_found")
var errPermissionDenied = errors.New("permission_denied")
var errUnauthenticated = errors.New("unauthenticated")

func NewAlreadyExistsError() error {
	return errAlreadyExists
}

func NewFailedPreconditionError() error {
	return errFailedPrecondition
}

func NewInvalidArgumentError() error {
	return errInvalidArgument
}

func NewNotFoundError() error {
	return errNotFound
}

func NewPermissionDeniedError() error {
	return errPermissionDenied
}

func NewUnauthenticatedError() error {
	return errUnauthenticated
}
