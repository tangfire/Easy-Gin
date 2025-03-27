通过 gin 的路由，已经把请求分配到具体的函数里面里面了，下面就要开始处理具体的业务逻辑了。

这里就进入 gin 封装的非常重要的的功能，对请求参数快速解析，让我们不纠结于参数的繁琐处理。当然这是对于比较标准的参数处理才可以，对于那些自定义的参数格式只能自己处理了。

# 参数风格

对于 RESTful 风格的 http 请求来说，参数的表现会有下面几种方式：

# URI参数

什么是 URI 参数？RESTful 风格的请求，某些请求的参数会通过 URI 来表现。

举个简单的例子：张三通过网上银行给李四转了 500 元，这个路由可以这么设计：

```go
xxx.com/:name/transfer/:money/to/:name

非常具体的体现：
xxx.com/zhangsan/transfer/500/to/lisi
```

当然你会说这个路由设计会比较丑陋，不过在 URI 里面增加参数有的时候是比较方便的，gin 支持这种方式获取参数。

```go
// This handler will match /user/john but will not match /user/ or /user
router.GET("/user/:name", uriFunc)
```

对于获取这种路由参数，gin 提供了两种方式去解析这种参数。

## 方式1：Param

```go
func uriFunc(c *gin.Context) {
    name := c.Param("name")
    c.String(http.StatusOK, "Hello %s", name)
}
```

## 方式2：bindUri

```go
type Person struct {
   Name string `uri:"name" binding:"required"`
}

func uriFunc(c *gin.Context) {
  var person Person
  if err := c.ShouldBindUri(&person); err != nil {
     c.JSON(400, gin.H{"msg": err.Error()})
     return
  }
  c.JSON(200, gin.H{"name": person.Name)
}
```

其实现原理很简单，就是在创建路由树的时候，将路由参数以及对应的值放入一个特定的 map 中即可。

```go
func (ps Params) Get(name string) (string, bool) {
    for _, entry := range ps {
      if entry.Key == name {
        return entry.Value, true
      }
    }
    return "", false
}
```

# QueryString Parameter

query String 即路由的 ? 之后的所带的参数，这种方式是比较常见的。

例如：/welcome?firstname=Jane&lastname=Doe

这里要注意的是，不管是 GET 还是 POST 都可以带 queryString Parameter。我曾经遇到某公司所有的参数都挂在 query string 上，这样做其实是不建议的，不过大家都这么做，只能顺其自然了。这么做的缺点很明显：

- 容易突破 URI 的长度限制，导致接口参数被截断。一般情况下服务器为了安全会对 URL 做长度限制，最大为2048
- 同时服务器也会对传输的大小也是有限制的，一般是 2k
- 当然这么做也是不安全的，都是明文的
这里就不具体罗列了，反正缺点挺多的。

这种参数也有两种获取方式：


## 方式1：Query

```go
firstname := c.DefaultQuery("firstname", "Guest")
lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
```


## 方式2：Bind

```go
type Person struct {
   FirstName  string `form:"name"`
}

func queryFunc(c *gin.Context) {
    var person Person
    if c.ShouldBindQuery(&person) == nil {
        log.Println(person.Name)
    }
}
```
实现原理：其实很简单就是将请求参数解析出来而已，利用的 net/url 的相关函数。


```go
//net/url.go:L1109
func (u *URL) Query() Values {
    v, _ := ParseQuery(u.RawQuery)
    return v
}
```

# Form

Form 一般还是更多用在跟前端的混合开发的情况下。Form 可以用于所有的方法 POST,GET,HEAD,PATCH ……

这种参数也有两种获取方式：

## 方式1：

```go
name := c.PostForm("name")
```

## 方式2：

```go
type Person struct {
    Name string `form:"name"`
}

func formFunc(c *gin.Context) {
    var person Person
    if c.ShouldBind(&person) == nil {
        log.Println(person.Name)
    }
}
```

# Json Body

son Body 是被使用最多的方式，基本上各种语言库对 json 格式的解析非常完善了，而且还在不断的推陈出新。

gin 对 json 的解析只有一种方式。

```go
type Person struct {
    Name string `json:"name"`
}

func jsonFunc(c *gin.Context) {
    var person Person
    if c.ShouldBind(&person) == nil {
        log.Println(person.Name)
    }
}
```

gin 默认是使用的 go 内置的 encoding/json 库，内置的 json 在 go 1.12 后性能得到了很大的提高。不过 Go 对接 PHP 的接口，如果用内置的 json 库简直就是一种折磨，gin 可以使用 jsoniter 来代替，只需要在编译的时候加上标志即可：”go build -tags=jsoniter .”，强烈建议对接 PHP 接口的同学，尝试 jsoniter 这个库，让你不再受 PHP 接口参数类型不确定之苦。

当然 gin 还支持其他类型参数的解析，如 Header，XML，YAML，Msgpack，Protobuf 等，这里就不再具体介绍了。

## Bind 系列函数的源码剖析


使用 gin 解析 request 的参数，按照我的实践来看，使用 Bind 系列函数还是比较好一点，因为这样请求的参数会比较好归档、分类，也有助于后续的接口升级，而不是将接口的请求参数分散不同的 handler 里面。

# 初始化 binding 相关对象


gin 在程序启动就会默认初始化好 binding 相关的变量

```go
// binding:L74
var (
 JSON          = jsonBinding{}
 XML           = xmlBinding{}
 Form          = formBinding{}
 Query         = queryBinding{}
 FormPost      = formPostBinding{}
 FormMultipart = formMultipartBinding{}
 ProtoBuf      = protobufBinding{}
 MsgPack       = msgpackBinding{}
 YAML          = yamlBinding{}
 Uri           = uriBinding{}
 Header        = headerBinding{}
)
```

# ShoudBind 与 MustBind 的区别

bind 相关的系列函数大体上分为两类 ShoudBind 和 MustBind。实现上基本一样，为了有区别的 MustBind 在解析失败的时候，返回 HTTP 400 状态。

MustBindWith:

```go
func (c *Context) MustBindWith(obj interface{}, b binding.Binding) error {
    if err := c.ShouldBindWith(obj, b); err != nil {
        c.AbortWithError(http.StatusBadRequest, err).SetType(ErrorTypeBind) // nolint: errcheck
        return err
    }
    return nil
}
```

ShoudBindWith:

```go
func (c *Context) ShouldBindWith(obj interface{}, b binding.Binding) error {
   return b.Bind(c.Request, obj)
}
```

# 匹配对应的参数 decoder

不管是 MustBind 还是 ShouldBind，总体上解析又可以分为两类：一种是让 gin 自己判断使用哪种 decoder，另外一种就是指定某种 decoder。自己判断使用哪种 decoder 比 指定 decoder 多了一步判断，其他的都是一样的。

```go
func (c *Context) ShouldBind(obj interface{}) error {
    b := binding.Default(c.Request.Method, c.ContentType())
    return c.ShouldBindWith(obj, b)
}

func Default(method, contentType string) Binding {
    if method == http.MethodGet {
        return Form
    }

    switch contentType {
    case MIMEJSON:
        return JSON
    case MIMEXML, MIMEXML2:
        return XML
    case MIMEPROTOBUF:
        return ProtoBuf
    case MIMEMSGPACK, MIMEMSGPACK2:
        return MsgPack
    case MIMEYAML:
        return YAML
    case MIMEMultipartPOSTForm:
        return FormMultipart
    default: // case MIMEPOSTForm:
        return Form
    }
}
```

ShouldBind/MustBind 会根据传入的 ContentType 来判断该使用哪种 decoder。不过对于 Header 和 Uri 方式的参数，只能用指定方式的decoder 了。

# 总结
本篇文章主要介绍了 gin 是如何快速处理客户端传递过的参数的。