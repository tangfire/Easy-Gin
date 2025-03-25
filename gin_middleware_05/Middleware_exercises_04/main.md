# 中间件联系

- 定义程序计时中间件，然后定义2个路由，执行函数后应该打印统计的执行时间，如下：

```go
package main

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
)

// 定义中间
func myTime(c *gin.Context) {
    start := time.Now()
    c.Next()
    // 统计时间
    since := time.Since(start)
    fmt.Println("程序用时：", since)
}

func main() {
    // 1.创建路由
    // 默认使用了2个中间件Logger(), Recovery()
    r := gin.Default()
    // 注册中间件
    r.Use(myTime)
    // {}为了代码规范
    shoppingGroup := r.Group("/shopping")
    {
        shoppingGroup.GET("/index", shopIndexHandler)
        shoppingGroup.GET("/home", shopHomeHandler)
    }
    r.Run(":8000")
}

func shopIndexHandler(c *gin.Context) {
    time.Sleep(5 * time.Second)
}

func shopHomeHandler(c *gin.Context) {
    time.Sleep(3 * time.Second)
}
```

效果演示：

![img](img.png)