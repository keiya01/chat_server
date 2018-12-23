package controller

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"

	"github.com/keiya01/chat_room/model"
	"github.com/keiya01/chat_room/service"
)

type ChatController struct{}

func NewChatController() *ChatController {
	return &ChatController{}
}

func (c *ChatController) Index(w http.ResponseWriter, r *http.Request) {
	s := service.NewService()
	chat := []model.Chat{}

	var resp model.Response
	if err := s.FindAll(&chat, "created_at asc"); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを取得できませんでした")

		json.NewEncoder(w).Encode(resp)

		return
	}

	resp.Data = chat

	json.NewEncoder(w).Encode(resp)

}

func (c *ChatController) Create(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		var err error

		for {
			var reply string

			if err = websocket.Message.Receive(ws, &reply); err != nil {
				fmt.Println("Can't receive")
				break
			}

			go func() {
				s := service.NewService()
				defer s.Close()

				p := model.NewChat(reply)
				s.Create(p)
			}()

			if err = websocket.Message.Send(ws, reply); err != nil {
				fmt.Println("Can't send")
				break
			}
		}
	}).ServeHTTP(w, r)
}
