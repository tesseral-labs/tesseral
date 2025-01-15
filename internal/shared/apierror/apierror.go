package apierror

type APIError struct {
	Message     string
	SourceError error
}

func New(message string, sourceError error) *APIError {
	return &APIError{
		Message:     message,
		SourceError: sourceError,
	}
}

func (a *APIError) Error() string {
	return a.Message
}

func (a *APIError) Unwrap() error {
	return a.SourceError
}
