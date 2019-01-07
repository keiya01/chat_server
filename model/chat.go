package model

type Chat struct {
	Model
	Body   string `json:"body"`
	UserID int    `json:"user_id"`
	RoomID int    `json:"room_id"`
}
