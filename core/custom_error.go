package core

import "fmt"

type NetworkError struct {
	Message string
	Cause   error
}

func NewNetworkError(message string, cause error) NetworkError {
	return NetworkError{
		Message: message,
		Cause:   cause,
	}
}

func (e NetworkError) Error() string {
	if e.Cause == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
}

type UnexpectedInputError struct {
	Message string
	Cause   error
}

func NewUnexpectedInputError(message string, cause error) UnexpectedInputError {
	return UnexpectedInputError{
		Message: message,
		Cause:   cause,
	}
}

func (e UnexpectedInputError) Error() string {
	if e.Cause == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
}

type InvalidUIError struct {
	Message string
	Cause   error
}

func NewInvalidUIError(message string, cause error) InvalidUIError {
	return InvalidUIError{
		Message: message,
		Cause:   cause,
	}
}

func (e InvalidUIError) Error() string {
	if e.Cause == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
}

type OtherError struct {
	Message string
	Cause   error
}

func NewOtherError(message string, cause error) OtherError {
	return OtherError{
		Message: message,
		Cause:   cause,
	}
}

func (e OtherError) Error() string {
	if e.Cause == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
}
