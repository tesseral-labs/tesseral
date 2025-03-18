package apierror

import (
	"connectrpc.com/connect"
	commonv1 "github.com/tesseral-labs/tesseral/internal/common/gen/tesseral/common/v1"
)

var errAlreadyExists = "already_exists"
var errFailedPrecondition = "failed_precondition"
var errInvalidArgument = "invalid_argument"
var errNotFound = "not_found"
var errPermissionDenied = "permission_denied"
var errUnauthenticated = "unauthenticated"
var errPasswordCompromised = "password_compromised"

func NewAlreadyExistsError(description string, sourceError error) error {
	apiErr := New(errAlreadyExists, sourceError)

	err := connect.NewError(connect.CodeAlreadyExists, apiErr)

	// Add details to the connect error
	if detail, detailErr := connect.NewErrorDetail(&commonv1.ErrorDetail{
		Description: description,
	}); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}

func NewFailedPreconditionError(description string, sourceError error) error {
	apiErr := New(errFailedPrecondition, sourceError)

	err := connect.NewError(connect.CodeFailedPrecondition, apiErr)

	// Add details to the connect error
	if detail, detailErr := connect.NewErrorDetail(&commonv1.ErrorDetail{
		Description: description,
	}); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}

func NewInvalidArgumentError(description string, sourceError error) error {
	apiErr := New(errInvalidArgument, sourceError)

	err := connect.NewError(connect.CodeInvalidArgument, apiErr)

	// Add details to the connect error
	if detail, detailErr := connect.NewErrorDetail(&commonv1.ErrorDetail{
		Description: description,
	}); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}

func NewNotFoundError(description string, sourceError error) error {
	apiErr := New(errNotFound, sourceError)

	err := connect.NewError(connect.CodeNotFound, apiErr)

	// Add details to the connect error
	if detail, detailErr := connect.NewErrorDetail(&commonv1.ErrorDetail{
		Description: description,
	}); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}

func NewPermissionDeniedError(description string, sourceError error) error {
	apiErr := New(errPermissionDenied, sourceError)

	err := connect.NewError(connect.CodePermissionDenied, apiErr)

	// Add details to the connect error
	if detail, detailErr := connect.NewErrorDetail(&commonv1.ErrorDetail{
		Description: description,
	}); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}

func NewUnauthenticatedError(description string, sourceError error) error {
	apiErr := New(errUnauthenticated, sourceError)

	err := connect.NewError(connect.CodeUnauthenticated, apiErr)

	// Add details to the connect error
	if detail, detailErr := connect.NewErrorDetail(&commonv1.ErrorDetail{
		Description: description,
	}); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}

func NewPasswordCompromisedError(description string, sourceError error) error {
	apiErr := New(errPasswordCompromised, sourceError)

	err := connect.NewError(connect.CodeFailedPrecondition, apiErr)

	// Add details to the connect error
	if detail, detailErr := connect.NewErrorDetail(&commonv1.ErrorDetail{
		Description: description,
	}); detailErr == nil {
		err.AddDetail(detail)
	}

	return err
}
