package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Println("中间件开始执行了!")
		c.Set("request", "中间件")
		status := c.Writer.Status()
		fmt.Println("中间件执行完毕", status)
		t2 := time.Since(t)
		fmt.Println("time:", t2)

	}
}

func main() {
	r := gin.Default()
	r.Use(MiddleWare())
	// {}为了代码规范
	{
		r.GET("/ce", func(c *gin.Context) {
			req, exists := c.Get("request")
			if exists {
				fmt.Println(req)
			}
			c.JSON(200, gin.H{"request": req})
		})
	}

	r.Run(":8003")
}
