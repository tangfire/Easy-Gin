package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func main() {
	validate := validator.New()
	email := "admin#admin.com" // 测试邮箱格式错误
	// email := ""              // 测试 required 错误

	// 验证变量
	err := validate.Var(email, "required,email")
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			fmt.Printf("Error: %s\n", e.Error())
		}
		return
	}
	fmt.Println("Validation passed")
}
