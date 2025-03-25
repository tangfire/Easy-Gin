package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.MaxMultipartMemory = 8 << 20

	r.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		}
		files := form.File["files"]

		for _, file := range files {
			if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("save file err: %s", err.Error()))
				return
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded", len(files)))
	})

	r.Run(":8003")
}
