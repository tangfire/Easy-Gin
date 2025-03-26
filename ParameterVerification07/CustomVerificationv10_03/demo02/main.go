package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func main() {
	validate := validator.New()
	type User struct {
		ID     int64  `json:"id" validate:"gt=0"`
		Name   string `json:"name" validate:"required"`
		Gender string `json:"gender" validate:"required,oneof=man woman"`
		Age    uint8  `json:"age" validate:"required,gte=0,lte=130"`
		Email  string `json:"email" validate:"required,email"`
	}
	user := &User{
		ID:     1,
		Name:   "frank",
		Gender: "boy",
		Age:    180,
		Email:  "gopher@88.com",
	}
	err := validate.Struct(user)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// output: Key: 'User.Age' Error:Field validation for 'Age' failed on the 'lte' tag
		// fmt.Println(validationErrors)
		//fmt.Println(validationErrors.Translate(trans))
		fmt.Println(validationErrors)
		return
	}
}
