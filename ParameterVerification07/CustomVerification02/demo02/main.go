package main

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Booking contains binded and validated data.
type Booking struct {
	//定义一个预约的时间大于今天的时间
	CheckIn time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	//gtfield=CheckIn退出的时间大于预约的时间
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func bookableDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return date.After(time.Now())
}

func main() {
	route := gin.Default()
	//注册验证
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//绑定第一个参数是验证的函数第二个参数是自定义的验证函数
		v.RegisterValidation("bookabledate", bookableDate)
	}

	route.GET("/5lmh", getBookable)
	route.Run()
}

func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// curl -X GET "http://localhost:8080/5lmh?check_in=2019-11-07&check_out=2019-11-20"
// curl -X GET "http://localhost:8080/5lmh?check_in=2019-09-07&check_out=2019-11-20"
// curl -X GET "http://localhost:8080/5lmh?check_in=2019-11-07&check_out=2019-11-01"
