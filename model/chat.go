package model

type Chat struct {
	Model
	Body   string `json:"body"`
	RoomID int    `json:"room_id"`
}
