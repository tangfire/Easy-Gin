package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 重定向
func main() {
	r := gin.Default()
	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.5lmh.com")
	})
	r.Run(":8003")
}
