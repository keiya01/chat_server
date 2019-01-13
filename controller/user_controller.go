package controller

import (
	"encoding/json"
	"github.com/keiya01/chat_room/model"
	"github.com/keiya01/chat_room/validation"

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

	resp := model.Response{}

	user := model.User{}
	user.ID = userID
	if err := s.Select("name, email, description").FindOne(&user); err != nil {
		resp.Error = model.NewError("ログインに失敗しました")
		resp.Data = map[string]string{}
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

	json.NewEncoder(w).Encode(resp)
}

func (u *UserController) Create(w http.ResponseWriter, r *http.Request) {
	s := service.NewService()
	defer s.Close()

	var params model.User
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		panic(err)
	}

	var resp model.Response

	if err, ok := validation.Validate(params); !ok {
		if ok := validation.Empty(params.Name); !ok {
			resp.Error = model.NewError("名前を入力してください")
		} else {
			resp.Error = model.NewError(err)
		}
		resp.Data = map[string]string{}

		json.NewEncoder(w).Encode(resp)
		return
	}

	encryptedPassword := auth.EncryptPassword(params.Password)

	user := model.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: encryptedPassword,
	}

	if err := s.Create(&user); err != nil {
		resp.Error = model.NewError("データを保存できませんでした")
		resp.Data = map[string]string{}

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
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		panic(err)
	}

	resp := model.Response{}

	if err, ok := validation.Validate(params); !ok {

		resp.Error = model.NewError(err)
		resp.Data = map[string]string{}

		json.NewEncoder(w).Encode(resp)
		return
	}

	user := model.User{}
	if err := s.FindOne(&user, "email = ?", params.Email); err != nil {
		resp.Error = model.NewError("メールアドレスが正しくありません")
		resp.Data = map[string]string{}

		json.NewEncoder(w).Encode(resp)
		return
	}

	isAuth := auth.ComparePassword(params.Password, user.Password)
	if !isAuth {
		resp.Error = model.NewError("パスワードが正しくありません")
		resp.Data = map[string]string{}

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

	param := model.User{}
	json.NewDecoder(r.Body).Decode(&param)
	params := map[string]interface{}{
		"Name":        param.Name,
		"Email":       param.Email,
		"Description": param.Description,
	}

	resp := model.Response{}

	user := model.User{}
	if err := s.Update(&user, params); err != nil {
		resp.Error = model.NewError("データを更新できませんでした")
		resp.Data = map[string]string{}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Data = user
	resp.Message = "データを更新しました"

	json.NewEncoder(w).Encode(resp)

}
