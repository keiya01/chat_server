package model

type User struct {
	Model
	Name        string `json:"name"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	Description string `json:"description"`
}

func (u User) SetErrorField(errorField string) (field string) {
	switch errorField {
	case "Email":
		field = "メールアドレス"
	case "Password":
		field = "パスワード"
	}

	return
}
