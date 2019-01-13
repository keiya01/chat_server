package controller

import (
	"encoding/json"
	"log"
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

	resp := model.Response{}
	rooms := []model.Room{}
	if err := s.Select("name, question, is_resolved, created_at").Pagination(initPage, nextPage).FindAll(&rooms, "created_at desc"); err != nil {
		log.Println(err)

		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを取得できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	roomAuthors := []model.User{}
	for _, room := range rooms {
		var gotUserIDList []int
		var user model.User
		user.ID = room.UserID

		if len(gotUserIDList) > 0 {
			for _, id := range gotUserIDList {
				if id == room.UserID {
					continue
				}
			}
		}

		if err := s.Select("id, name, email, description").FindOne(&user); err != nil {
			log.Println(err)
			w.Header().Add("Content-Type", "application/json")
			resp.Error = model.NewError("チャットのユーザー情報を取得できませんでした")

			json.NewEncoder(w).Encode(resp)
			return
		}

		gotUserIDList = append(gotUserIDList, room.UserID)
		roomAuthors = append(roomAuthors, user)
	}

	resp.Data = map[string]interface{}{
		"rooms": map[string]interface{}{
			"roomData": rooms,
			"authors":  roomAuthors,
		},
	}

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
	if err := s.Select("name, user_id, question, is_resolved, created_at").FindOne(&room); err != nil {
		log.Println(err)

		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("ルームのデータを取得できませんでした")

		json.NewEncoder(w).Encode(resp)
		return
	}

	var roomAuthor model.User
	roomAuthor.ID = room.UserID
	if err := s.Select("id, name, email, question").FindOne(&roomAuthor); err != nil {
		log.Println(err)

		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("ルーム作成者のデータを取得できませんでした")

		json.NewEncoder(w).Encode(resp)
		return
	}

	initPage, nextPage := SetNextPage(req)
	var chats []model.Chat
	if err := s.Select("body, user_id, room_id, created_at").Pagination(initPage, nextPage).FindAll(&chats, "created_at asc", "room_id = ?", roomID); err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("チャットのデータを取得できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	var chatUsers []model.User
	for _, chat := range chats {
		var gotUserIDList []int
		var user model.User
		user.ID = chat.UserID

		if len(gotUserIDList) > 0 {
			for _, id := range gotUserIDList {
				if id == chat.UserID {
					continue
				}
			}
		}

		if err := s.Select("id, name, email, description").FindOne(&user); err != nil {
			log.Println(err)
			w.Header().Add("Content-Type", "application/json")
			resp.Error = model.NewError("チャットのユーザー情報を取得できませんでした")

			json.NewEncoder(w).Encode(resp)
			return
		}

		gotUserIDList = append(gotUserIDList, chat.UserID)
		chatUsers = append(chatUsers, user)
	}

	resp.Data = map[string]interface{}{
		"room": map[string]interface{}{
			"roomData": room,
			"author":   roomAuthor,
		},
		"chats": map[string]interface{}{
			"chatsData": chats,
			"chatUser":  chatUsers,
		},
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
		log.Println(err)

		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("ルームの作成に失敗しました")

		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Data = room
	resp.Message = "ルームの作成に成功しました。"

	json.NewEncoder(w).Encode(resp)

}
