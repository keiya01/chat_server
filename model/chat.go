package model

type Chat struct {
	Model
	Body string `json:"body"`
}

func NewChat(body string) *Chat {
	return &Chat{
		Body: body,
	}
}
