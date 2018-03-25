package errors

var UnexpectedErrorCode = 0
var ChangesetErrorCode = 1
var NotFoundErrorCode = 2
var DuplicateErrorCode = 3
var AssociationErrorCode = 4

type Error struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Code    int    `json:"json,omitempty`
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) UnexpectedError() bool {
	return e.Code == UnexpectedErrorCode
}

func (e Error) ChangesetError() bool {
	return e.Code == ChangesetErrorCode
}

func (e Error) NotFoundError() bool {
	return e.Code == NotFoundErrorCode
}

func (e Error) DuplicateError() bool {
	return e.Code == DuplicateErrorCode
}

func (e Error) AssociationError() bool {
	return e.Code == AssociationErrorCode
}

func New(message string, field string, code int) Error {
	return Error{message, field, code}
}

func UnexpectedError(message string) Error {
	return Error{
		Message: message,
		Code:    UnexpectedErrorCode,
	}
}

func NotFoundError(message string) Error {
	return Error{
		Message: message,
		Code:    NotFoundErrorCode,
	}
}

func ChangesetError(message string, field string) Error {
	return Error{
		Message: message,
		Field:   field,
		Code:    ChangesetErrorCode,
	}
}

func DuplicateError(message string, field string) Error {
	return Error{
		Message: message,
		Field:   field,
		Code:    DuplicateErrorCode,
	}
}

func AssociationError(message string, field string) Error {
	return Error{
		Message: message,
		Field:   field,
		Code:    AssociationErrorCode,
	}
}

type Errors []Error

func (es Errors) Error() string {
	var messages string

	for i, e := range es {
		messages += e.Error()

		if i < len(es)-1 {
			messages += ", "
		}
	}

	return messages
}
