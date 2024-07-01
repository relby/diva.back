package domainerrors

type BaseError struct {
	Message string
}

func (err *BaseError) Error() string {
	return err.Message
}

type NotFoundError struct {
	BaseError
}

func (err *NotFoundError) Error() string {
	return err.Message
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		BaseError: BaseError{
			Message: message,
		},
	}
}

type ValidationError struct {
	BaseError
}

func (err *ValidationError) Error() string {
	return err.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		BaseError: BaseError{
			Message: message,
		},
	}
}
