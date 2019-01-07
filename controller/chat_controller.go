package controller

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"

	"github.com/keiya01/chat_room/http/request"
	"github.com/keiya01/chat_room/model"
	"github.com/keiya01/chat_room/service"
)

type ChatController struct{}

func NewChatController() *ChatController {
	return &ChatController{}
}

func (c *ChatController) Create(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
	}

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		var err error

		for {
			var reply string

			if err = websocket.Message.Receive(ws, &reply); err != nil {
				log.Panicln(err)
				fmt.Println("Can't receive")
				break
			}

			go func() {
				s := service.NewService()
				defer s.Close()

				param := request.GetParam(r, "room_id")
				roomID, err := strconv.Atoi(param)
				if err != nil {
					panic(err)
				}

				chat := model.Chat{
					Body:   reply,
					UserID: userID,
					RoomID: roomID,
				}
				s.Create(&chat)
			}()

			if err = websocket.Message.Send(ws, reply); err != nil {
				log.Println(err)
				fmt.Println("Can't send")
				break
			}
		}
	}).ServeHTTP(w, r)
}
