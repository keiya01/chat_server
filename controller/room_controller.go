package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/keiya01/chat_room/http/request"
	"github.com/keiya01/chat_room/model"
	"github.com/keiya01/chat_room/service"
)

type RoomController struct{}

func NewRoomController() *RoomController {
	return &RoomController{}
}

func (r *RoomController) Index(w http.ResponseWriter, req *http.Request) {
	s := service.NewService()
	defer s.Close()

	initPage, nextPage := SetNextPage(req)

	var resp model.Response
	var room []model.Room
	if err := s.Select("name, question, is_resolved, created_at").Pagination(initPage, nextPage).FindAll(&room, "created_at desc"); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを取得できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	resp.Data = room

	json.NewEncoder(w).Encode(resp)
}

func (r *RoomController) Show(w http.ResponseWriter, req *http.Request) {
	s := service.NewService()

	param := request.GetParam(req, "room_id")
	roomID, err := strconv.Atoi(param)
	if err != nil {
		panic(err)
	}

	var resp model.Response

	var room model.Room
	room.ID = roomID
	if err := s.Select("name, question, is_resolved, created_at").FindOne(&room); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを取得できませんでした")

		json.NewEncoder(w).Encode(resp)
	}

	initPage, nextPage := SetNextPage(req)
	var chat []model.Chat
	if err := s.Select("body, room_id, created_at").Pagination(initPage, nextPage).FindAll(&chat, "created_at asc", "room_id = ?", roomID); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを取得できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	resp.Data = map[string]interface{}{
		"room": room,
		"chat": chat,
	}

	json.NewEncoder(w).Encode(resp)

}

func (r *RoomController) Create(w http.ResponseWriter, req *http.Request) {
	s := service.NewService()
	defer s.Close()

	var room model.Room
	json.NewDecoder(req.Body).Decode(&room)
	var resp model.Response
	if err := s.Create(&room); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")

		resp.Error = model.NewError("ルームの作成に失敗しました")

		json.NewEncoder(w).Encode(resp)
	}

	resp.Data = room
	resp.Message = "ルームの作成に成功しました。"

	json.NewEncoder(w).Encode(resp)

}
