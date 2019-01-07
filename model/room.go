package model

type Room struct {
	Model
	Name       string `json:"name"`
	Question   string `json:"question"`
	IsResolved bool   `json:"is_resolved"`
	UserID     int    `json:"user_id"`
}
