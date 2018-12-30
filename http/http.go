package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"log"
	"net/http"

	"github.com/keiya01/chat_room/auth"
	"github.com/keiya01/chat_room/controller"
)

// Server Server
type Server struct {
	*chi.Mux
}

// NewServer Server構造体のコンストラクタ
func NewServer() *Server {
	return &Server{
		Mux: chi.NewRouter(),
	}
}

func (s *Server) routes() {
	cors := corsNew()
	s.Use(cors.Handler)
	s.Use(middleware.RequestID)
	s.Use(middleware.Logger)
	s.Use(middleware.URLFormat)
	s.Use(render.SetContentType(render.ContentTypeJSON))
	s.Use(auth.JWTAuthentication) // ここでJWTの認証を行う

	c := controller.NewChatController()
	r := controller.NewRoomController()
	u := controller.NewUserController()

	s.Route("/api", func(api chi.Router) {
		api.Route("/chat", func(chat chi.Router) {
			chat.Get("/{room_id}", c.Create)
		})
		api.Route("/rooms", func(rooms chi.Router) {
			rooms.Get("/", r.Index)
			rooms.Get("/{room_id}", r.Show)
			rooms.Post("/create", r.Create)
		})
		api.Route("/users", func(users chi.Router) {
			users.Post("/", u.Create)
			users.Post("/login", u.Login)
			users.Patch("/update", u.Update)
		})
	})

}

func (s *Server) Run() {
	log.Print("Server running ...")

	s.routes()

	if err := http.ListenAndServe(":8686", s); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func corsNew() *cors.Cors {
	acceptURI := "http://localhost:3000"
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{acceptURI},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return cors
}
