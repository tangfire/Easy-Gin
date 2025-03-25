# Sessions

gorilla/sessions为自定义session后端提供cookie和文件系统session以及基础结构。

主要功能是：

- 简单的API：将其用作设置签名（以及可选的加密）cookie的简便方法。
- 内置的后端可将session存储在cookie或文件系统中。
- Flash消息：一直持续读取的session值。
- 切换session持久性（又称“记住我”）和设置其他属性的便捷方法。
- 旋转身份验证和加密密钥的机制。
- 每个请求有多个session，即使使用不同的后端也是如此。
- 自定义session后端的接口和基础结构：可以使用通用API检索并批量保存来自不同商店的session。


代码：

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/sessions"
)

// 初始化一个cookie存储对象
// something-very-secret应该是一个你自己的密匙，只要不被别人知道就行
var store = sessions.NewCookieStore([]byte("something-very-secret"))

func main() {
    http.HandleFunc("/save", SaveSession)
    http.HandleFunc("/get", GetSession)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("HTTP server failed,err:", err)
        return
    }
}

func SaveSession(w http.ResponseWriter, r *http.Request) {
    // Get a session. We're ignoring the error resulted from decoding an
    // existing session: Get() always returns a session, even if empty.

    //　获取一个session对象，session-name是session的名字
    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 在session中存储值
    session.Values["foo"] = "bar"
    session.Values[42] = 43
    // 保存更改
    session.Save(r, w)
}
func GetSession(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    foo := session.Values["foo"]
    fmt.Println(foo)
}
```

---

这段代码使用 `gorilla/sessions` 包实现了基于 Cookie 的**服务端 Session 管理**功能，核心逻辑分为 Session 存储和读取两部分。以下是逐层解析与技术要点：

---

### 一、代码核心逻辑解析
#### 1. **Session 存储初始化**
```go
var store = sessions.NewCookieStore([]byte("something-very-secret"))
```
- **作用**：创建一个基于 Cookie 的 Session 存储对象（`CookieStore`）。
- **参数说明**：
    - `[]byte("something-very-secret")`：密钥用于加密和签名 Session 数据，需替换为随机生成的强密钥以提高安全性。
- **底层机制**：Session 数据会被序列化、加密后存储在客户端的 Cookie 中（而非服务端），适合无状态场景。

#### 2. **Session 存储（`/save` 路由）**
```go
func SaveSession(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")  // 获取或创建 Session
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    session.Values["foo"] = "bar"  // 写入数据
    session.Values[42] = 43
    session.Save(r, w)  // 必须调用保存操作
}
```
- **`store.Get`**：根据请求中的 Cookie 获取名为 `session-name` 的 Session。若不存在，则创建新 Session。
- **`session.Values`**：以键值对形式存储任意类型数据，数据会被加密后写入 Cookie。
- **`session.Save`**：将修改后的 Session 数据序列化并加密，通过响应头 `Set-Cookie` 发送给客户端。

#### 3. **Session 读取（`/get` 路由）**
```go
func GetSession(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")  // 获取 Session
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    foo := session.Values["foo"]  // 读取数据
    fmt.Println(foo)
}
```
- **数据解密**：客户端再次请求时，Cookie 中的加密数据会被自动解密并还原为 `session.Values`。
- **依赖 Cookie**：Session 的持久化完全依赖客户端 Cookie 传递 Session ID 和加密数据。

---

### 二、关键技术细节与安全性
#### 1. **Cookie 存储限制**
- **数据大小限制**：单个 Cookie 通常不超过 4KB，因此 `CookieStore` 仅适合存储少量数据（如用户 ID、令牌）。
- **加密与签名**：
    - **签名**：防止数据篡改。
    - **加密**：防止敏感信息泄露（需设置 `securecookie.GenerateRandomKey` 生成强密钥）。

#### 2. **Session 生命周期管理**
- **过期时间**：默认情况下，Session 在浏览器关闭后失效。可通过以下代码设置过期时间：
  ```go
  session.Options.MaxAge = 3600  // 单位：秒
  session.Save(r, w)
  ```
- **主动删除**：设置 `MaxAge = -1` 并调用 `session.Save` 可立即删除 Session。

#### 3. **多 Session 支持**
- **命名 Session**：通过不同名称（如 `session-name`）支持多个独立 Session，适用于多角色或多模块场景。

---

### 三、扩展应用场景
#### 1. **分布式 Session 存储**
- **问题**：`CookieStore` 依赖客户端存储，不适合分布式系统。可改用 `RedisStore` 或 `MySQLStore` 实现服务端集中存储。
- **示例改造**：
  ```go
  // 使用 Redis 存储
  var store = sessions.NewRedisStore(redisClient, "session-prefix")
  ```

#### 2. **安全增强**
- **SameSite 属性**：防止 CSRF 攻击（需设置 `session.Options.SameSite = http.SameSiteStrictMode`）。
- **HttpOnly 与 Secure**：
  ```go
  session.Options.HttpOnly = true  // 禁止 JS 访问
  session.Options.Secure = true    // 仅 HTTPS 传输
  ```

---

### 四、请求测试示例
#### 1. 存储 Session
```bash
curl -v http://localhost:8080/save
```
**响应头**包含加密的 Session Cookie：
```
Set-Cookie: session-name=MTcx...; Path=/; Max-Age=3600; HttpOnly
```

#### 2. 读取 Session
```bash
curl -H "Cookie: session-name=MTcx..." http://localhost:8080/get
```
**控制台输出**：`bar`

---

### 总结
这段代码展示了 `gorilla/sessions` 的**基础用法**，核心特点包括：
1. **客户端存储加密 Session**：通过 Cookie 实现无状态会话管理。
2. **数据安全**：依赖加密和签名保障数据完整性。
3. **灵活扩展**：支持替换存储后端（如 Redis、文件系统）。

实际开发中需注意**密钥管理**和**存储介质选择**，避免单点故障和安全漏洞。对于高并发场景，建议采用服务端集中存储方案。

删除session的值：

```go
    // 删除
    // 将session的最大存储时间设置为小于零的数即为删除
    session.Options.MaxAge = -1
    session.Save(r, w)
```

