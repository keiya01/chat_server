package controller

import (
	"encoding/json"
	"github.com/keiya01/chat_room/model"
	"log"

	"github.com/keiya01/chat_room/auth"
	"github.com/keiya01/chat_room/service"
	"net/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) Show(w http.ResponseWriter, r *http.Request) {
	s := service.NewService()
	defer s.Close()

	userID, ok := getUserID(r)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
	}

	var resp model.Response

	var user model.User
	user.ID = userID
	if err := s.Select("name, email, description").FindOne(&user); err != nil {
		w.Header().Add("Content-Type", "application/json")

		resp.Error = model.NewError("ログインに失敗しました")
		json.NewEncoder(w).Encode(resp)
	}

	token := auth.JWTToken{
		UserID:    user.ID,
		UserEmail: user.Email,
	}
	jwtToken := token.GetJWTToken()

	resp.Token = jwtToken
	resp.Data = user

	json.NewEncoder(w).Encode(resp)
}

func (u *UserController) Create(w http.ResponseWriter, r *http.Request) {
	s := service.NewService()
	defer s.Close()

	var params model.User
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		panic(err)
	}

	encryptedPassword := auth.EncryptPassword(params.Password)

	user := model.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: encryptedPassword,
	}

	var resp model.Response
	if err := s.Create(&user); err != nil {
		log.Println(err)

		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("データを保存できませんでした")

		json.NewEncoder(w).Encode(resp)
		return
	}

	token := auth.JWTToken{
		UserID:    user.ID,
		UserEmail: user.Email,
	}
	jwtToken := token.GetJWTToken()

	resp.Token = jwtToken
	resp.Data = user
	resp.Message = "データを保存しました"

	json.NewEncoder(w).Encode(resp)

}

func (u *UserController) Login(w http.ResponseWriter, r *http.Request) {
	s := service.NewService()
	defer s.Close()

	var params model.User
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		panic(err)
	}

	var resp model.Response
	user := model.User{}
	if err := s.FindOne(&user, "email = ?", params.Email); err != nil {
		log.Println(err)

		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("ログインに失敗しました")

		json.NewEncoder(w).Encode(resp)
	}

	isAuth := auth.ComparePassword(params.Password, user.Password)
	if !isAuth {
		w.Header().Add("Content-Type", "application/json")
		resp.Error = model.NewError("ログインに失敗しました")

		json.NewEncoder(w).Encode(resp)
		return
	}

	token := auth.JWTToken{
		UserID:    user.ID,
		UserEmail: user.Email,
	}
	jwtToken := token.GetJWTToken()

	resp.Token = jwtToken
	resp.Data = user
	resp.Message = "ログインしました"

	json.NewEncoder(w).Encode(resp)

}

func (u *UserController) Update(w http.ResponseWriter, r *http.Request) {
	s := service.NewService()
	defer s.Close()

	var param model.User
	json.NewDecoder(r.Body).Decode(&param)
	params := map[string]interface{}{
		"Name":        param.Name,
		"Email":       param.Email,
		"Description": param.Description,
	}

	var resp model.Response

	var user model.User
	if err := s.Update(&user, params); err != nil {
		log.Println(err)

		w.Header().Add("Content-Type", "applicaion/json")
		resp.Error = model.NewError("データを更新できませんでした")
		json.NewEncoder(w).Encode(resp)
	}

	resp.Data = user
	resp.Message = "データを更新しました"

	json.NewEncoder(w).Encode(resp)

}
