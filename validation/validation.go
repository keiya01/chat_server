package validation

import (
	"gopkg.in/go-playground/validator.v9"
	"log"
)

type Validator interface {
	SetErrorField(string) string
}

func Validate(v Validator) (string, bool) {
	if errors := NewValidate(v); errors != nil {
		var errorMsg string
		for _, err := range errors {
			field := v.SetErrorField(err["field"])
			if m := SetValidateMessage(err["valid"], field); m != "" {
				errorMsg = m
				break
			}
		}

		return errorMsg, false
	}

	return "", true
}

var validate = validator.New()

func NewValidate(model interface{}) []map[string]string {
	err := validate.Struct(model)

	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}

		var errors []map[string]string
		for _, err := range err.(validator.ValidationErrors) {
			errorMap := make(map[string]string)
			errorMap["field"] = err.Field()
			errorMap["valid"] = err.Tag()
			errors = append(errors, errorMap)
		}

		return errors
	}

	return nil
}

func Empty(param string) bool {

	errs := validate.Var(param, "required")

	if errs != nil {
		log.Println("validation_error: ", errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return false
	}

	return true
}

func SetValidateMessage(validate string, field string) (errorMsg string) {
	switch validate {
	case "required":
		errorMsg = field + "を入力してください"
	case "email":
		errorMsg = field + "が正しくありません"
	}

	return
}
