package core

type BaseError struct {
	ErrorCode   string `json:"errorCode"`
	Description string `json:"description"`
}

func UnexpectedError(description string) *BaseError {
	return &BaseError{
		ErrorCode:   "UNEXPECTED_ERROR",
		Description: description,
	}
}

func ValidationError(errorList string) *BaseError {
	return &BaseError{
		ErrorCode:   "VALIDATION_ERROR",
		Description: errorList,
	}
}
