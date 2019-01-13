package model

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

func NewError(errMsg string) *Error {
	return &Error{
		IsError: true,
		Message: errMsg,
	}
}
