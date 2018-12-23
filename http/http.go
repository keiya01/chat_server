package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"log"
	"net/http"

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

	c := controller.NewChatController()

	s.Route("/api", func(api chi.Router) {
		api.Route("/chat", func(chat chi.Router) {
			chat.Get("/", c.Index)
			chat.Get("/create", c.Create)
			// posts.Post("/create", p.Create)
			// posts.Put("/{id}/update", p.Update)
			// posts.Delete("/{id}/delete", p.Delete)
		})
		// api.Route("/users", func(users chi.Router) {
		// 	users.Post("/login", u.Login)
		// 	users.Post("/create", u.Create)
		// 	users.Get("/{id}", u.Show)
		// 	users.Put("/{id}/update", u.Update)
		// 	users.Delete("/{id}/delete", u.Delete)
		// })
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
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return cors
}
