package server

var (
	NilError = &StatusError{Code: StatusOK, message: "OK"}
)

// Error represents a handler error. It provides methods for a server status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int32
}

// StatusError represents an error with an associated server status code.
type StatusError struct {
	Code int32
	Err  error

	message string
}

// Allows StatusError to satisfy the error interface.
func (se *StatusError) Error() string {
	if len(se.message) > 0 {
		return se.message
	}

	return se.Err.Error()
}

// Returns our server status code.
func (se *StatusError) Status() int32 {
	return se.Code
}

func NewInternalServerError(err error) *StatusError {
	return &StatusError{
		Code: StatusInternalServerError,
		Err:  err,
	}
}

func NewBadRequestError(err error) *StatusError {
	return &StatusError{
		Code: StatusBadRequest,
		Err:  err,
	}
}
