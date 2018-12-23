package main

import (
	"github.com/keiya01/chat_room/http"
	"github.com/keiya01/chat_room/service/database"
	"github.com/keiya01/chat_room/service/migrate"
)

func main() {
	DBHandler := database.NewHandler()
	migrate.Set(DBHandler.DB)
	s := http.NewServer()
	s.Run()
}
