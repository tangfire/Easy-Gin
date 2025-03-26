package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
对绑定解析到结构体上的参数，自定义验证功能
比如我们需要对URL的接受参数进行判断，判断用户名是否为root如果是root通过否则返回false
*/
type Login struct {
	User    string `uri:"user" validate:"checkName"`
	Pssword string `uri:"password"`
}

// 自定义验证函数
func checkName(fl validator.FieldLevel) bool {
	if fl.Field().String() != "root" {
		return false
	}
	return true
}
func main() {
	r := gin.Default()
	validate := validator.New()
	r.GET("/:user/:password", func(c *gin.Context) {
		var login Login
		//注册自定义函数，与struct tag关联起来
		err := validate.RegisterValidation("checkName", checkName)
		if err := c.ShouldBindUri(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = validate.Struct(login)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println(err)
			}
			return
		}
		fmt.Println("success")
	})
	r.Run()
}
