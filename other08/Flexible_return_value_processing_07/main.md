# 灵活的返回值处理

上一篇文章是关于如何快速解析客户端传递过来的参数的，参数解析出来后就开始了我们的业务的开发流程了。

业务处理的过程 gin 并没有给出对应的设计，这给业务开发带来了很多不方便的地方，很多公司会基于 gin 做二次开发，定制契合公司基础技术建设的框架升级，关于 gin 定制框架的内容这里不再详细展开，请关注后续文章。

经过业务逻辑框架的处理，已经有了对应的处理结果了，需要结果返回给客户端了，本篇文章主要介绍 gin 是如何处理响应结果的。

仍然以原生的 net/http 简单的例子开始我们的源码分析。

```go
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World"))
    })

    if err := http.ListenAndServe(":8000", nil); err != nil {
        fmt.Println("start http server fail:", err)
    }
}
```

output:

```go
▶ curl -i -XGET 127.0.0.1:8000
HTTP/1.1 200 OK
Date: Sun, 10 Oct 2021 10:28:15 GMT
Content-Length: 11
Content-Type: text/plain; charset=utf-8

Hello World
```

可以看到调用 http.ResponseWriter.Write 即可将响应结果返回给客户端。不过也可以看出一些问题：

- 这个函数返回的值是默认的 text/plain 类型。如果想返回 application/json 就需要调用额外的设置 header 相关函数。
- 这个函数只能接受 []byte 类型变量。一般情况下，我们经过业务逻辑处理过的数据都是结构体类型的，要使用 Write，需要把结构体转换 []byte，这个就太不方便。 类似 gin 提供的参数处理，gin 同样提供了很多格式的返回值，能让我们简化返回数据的处理。


下面是 gin 提供的 echo server，无需任何处理，就能返回一个 json 类型的返回值。

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run()
}
```

output:

```go
▶ curl -i -XGET 127.0.0.1:8080/ping   
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 10 Oct 2021 05:40:21 GMT
Content-Length: 18

{"message":"pong"}
```

当然 gin 还提其他类型格式的返回值，如 xml, yaml, protobuf 等。

```go
var (
 _ Render     = JSON{}
 _ Render     = IndentedJSON{}
 _ Render     = SecureJSON{}
 _ Render     = JsonpJSON{}
 _ Render     = XML{}
 _ Render     = String{}
 _ Render     = Redirect{}
 _ Render     = Data{}
 _ Render     = HTML{}
 _ HTMLRender = HTMLDebug{}
 _ HTMLRender = HTMLProduction{}
 _ Render     = YAML{}
 _ Render     = Reader{}
 _ Render     = AsciiJSON{}
 _ Render     = ProtoBuf{}
)
```

本文仅以比较常见的 json 类型格式的返回值阐述 gin 对 ResponseWriter 的实现原理。

# 源码分析


## 1. 设置 json 的返回格式

```go
// gin/context.go:L956
func (c *Context) JSON(code int, obj interface{}) {
   c.Render(code, render.JSON{Data: obj})
}
```

初始化 render.JSON 类型变量

## 2. 通过 interface 动态转发调用真正的 json 处理函数

```go
// gin/context.go:L904
func (c *Context) Render(code int, r render.Render) {
    c.Status(code)

    if !bodyAllowedForStatus(code) {
        r.WriteContentType(c.Writer)
        c.Writer.WriteHeaderNow()
        return
    }

    if err := r.Render(c.Writer); err != nil {
        panic(err)
    }
}
```

- 设置 Http status
- 处理 Http status 为 100 - 199、204、304 的情况
- 调用真正的 json 处理函数

## 3. 组装 response 数据

```go
// gin/render/json.go:L67
func WriteJSON(w http.ResponseWriter, obj interface{}) error {
    writeContentType(w, jsonContentType)
    jsonBytes, err := json.Marshal(obj)
    if err != nil {
        return err
    }
    _, err = w.Write(jsonBytes)
    return err
}
```

设置 response Header 的 Content-Type 为 application/json; charset=utf-8
将要返回数据编码成 json 字符串
写入 gin.responseWriter

## 4. 写入真正的 http.ResponseWriter

```go
// gin/response_writer.go:L76
func (w *responseWriter) Write(data []byte) (n int, err error) {
    w.WriteHeaderNow()
    n, err = w.ResponseWriter.Write(data)
    w.size += n
    return
}
```

这里 gin 实现了 ResponseWriter interface，对原生的 response 做了一定的扩展，不过最终依然是调用 net/http 的 response.Write 完成对请求数据的最终写入。

# 总结
本篇文章主要介绍了 gin 是如何完成对数据的组装然后返回给客户端的。写到这里基本上 gin 的整个流程就梳理完成了。gin 提供的功能就这么多，第一篇源码分析文章我提到 gin 是个 httprouter 基本就是这个原因。

