package handle

type Error struct {
	IsError bool   `json:"error"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(text string) error {
	return &Error{
		IsError: true,
		Message: text,
	}
}
