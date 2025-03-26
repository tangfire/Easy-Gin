# 自定义验证v10

## 介绍

Validator 是基于 tag（标记）实现结构体和单个字段的值验证库，它包含以下功能：

- 使用验证 tag（标记）或自定义验证器进行跨字段和跨结构体验证。
- 关于 slice、数组和 map，允许验证多维字段的任何或所有级别。
- 能够深入 map 键和值进行验证。
- 通过在验证之前确定接口的基础类型来处理类型接口。
- 处理自定义字段类型（如 sql 驱动程序 Valuer）。
- 别名验证标记，它允许将多个验证映射到单个标记，以便更轻松地定义结构体上的验证。
- 提取自定义的字段名称，例如，可以指定在验证时提取 JSON 名称，并在生成的 FieldError 中使用该名称。
- 可自定义 i18n 错误消息。
- Web 框架 gin 的默认验证器。


## 安装

使用 go get：

```go
go get github.com/go-playground/validator/v10
```



然后将 Validator 包导入到代码中：

```go
import "github.com/go-playground/validator/v10"

```

## 变量验证

Var 方法使用 tag（标记）验证方式验证单个变量。

```go
func (*validator.Validate).Var(field interface{}, tag string) error
```

它接收一个 interface{} 空接口类型的 field 和一个 string 类型的 tag，返回传递的非法值得无效验证错误，否则将 nil 或 ValidationErrors 作为错误。如果错误不是 nil，则需要断言错误去访问错误数组，例如：


```go
validationErrors := err.(validator.ValidationErrors)
```

如果是验证数组、slice 和 map，可能会包含多个错误。

示例代码：


```go
func main() {
  validate := validator.New()
  // 验证变量
  email := "admin#admin.com"
  email := ""
  err := validate.Var(email, "required,email")
  if err != nil {
    validationErrors := err.(validator.ValidationErrors)
    fmt.Println(validationErrors)
    // output: Key: '' Error:Field validation for '' failed on the 'email' tag
    // output: Key: '' Error:Field validation for '' failed on the 'required' tag
    return
  }
}
```



这段代码使用了 Go 的 `validator` 库对变量 `email` 进行验证，但存在一些需要修正和注意的问题。以下是关键解析：

---

### **一、代码问题分析**
1. **变量重复声明错误**  
   代码中 `email` 被声明了两次（`email := "admin#admin.com"` 和 `email := ""`），这在 Go 中会导致编译错误。正确的做法应该是：
   ```go
   email := "admin#admin.com"  // 第一次声明
   email = ""                  // 后续赋值（去掉冒号）
   ```

2. **验证逻辑矛盾**
    - `required` 规则要求字段非空，而 `email` 规则要求符合邮箱格式。
    - 当 `email` 为空时，`required` 验证会失败，此时 `email` 格式验证不会执行（因为值为空）。但根据注释，代码预期会同时触发两个错误，这需要修正逻辑。

---

### **二、验证规则解析**
#### **验证器初始化**
```go
validate := validator.New()  // 创建 validator 实例
```
- 通过 `validator.New()` 初始化一个验证器实例，用于后续验证操作。

#### **变量验证**
```go
err := validate.Var(email, "required,email")
```
- **`Var` 方法作用**  
  验证单个变量（非结构体字段）是否符合规则。
- **规则说明**
    - `required`：字段值不能为空（空字符串、零值等）。
    - `email`：字段必须符合 RFC 5322 邮箱格式（如 `user@example.com`）。

#### **错误处理**
```go
if err != nil {
    validationErrors := err.(validator.ValidationErrors)
    fmt.Println(validationErrors)
}
```
- **错误类型断言**  
  将错误转换为 `validator.ValidationErrors` 类型，以便获取多个验证失败的详细信息。
- **输出示例**  
  若 `email` 为空，输出：
  ```text
  Key: '' Error:Field validation for '' failed on the 'required' tag
  ```
  若 `email` 为 `admin#admin.com`（格式错误），输出：
  ```text
  Key: '' Error:Field validation for '' failed on the 'email' tag
  ```

---

### **三、修正后的代码**
```go
package main

import (
    "fmt"
    "github.com/go-playground/validator/v10"
)

func main() {
    validate := validator.New()
    email := "admin#admin.com"  // 测试邮箱格式错误
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
```

---

### **四、运行结果示例**
1. **当 `email` 为空**
   ```text
   Error: '' failed on the 'required' tag
   ```
2. **当 `email` 格式错误**
   ```text
   Error: '' failed on the 'email' tag
   ```
3. **当 `email` 合法（如 `test@example.com`）**
   ```text
   Validation passed
   ```

---

### **五、扩展知识**
- **多规则执行顺序**  
  `validator` 默认会执行所有规则，即使前一个规则已失败。例如，若 `email` 为空，`required` 和 `email` 规则都会触发错误。
- **自定义错误消息**  
  可通过 `RegisterTranslation` 方法实现错误消息中文化（参考网页1中的翻译器方案）。
- **验证器性能优化**  
  `validator` 基于反射实现，高频调用时建议复用验证器实例。

如需进一步了解其他验证规则（如 `min`、`max`、`oneof`），可参考网页1和网页4的标签说明。


## 结构体验证

结构体验证结构体公开的字段，并自动验证嵌套结构体，除非另有说明。

```go
func (*validator.Validate).Struct(s interface{}) error

```



它接收一个 interface{} 空接口类型的 s，返回传递的非法值得无效验证错误，否则将 nil 或 ValidationErrors 作为错误。如果错误不是 nil，则需要断言错误去访问错误数组，例如：

```go
validationErrors := err.(validator.ValidationErrors)

```

实际上，Struct 方法是调用的 StructCtx 方法，因为本文不是源码讲解，所以此处不展开赘述，如有兴趣，可以查看源码。

示例代码：


```go
func main() {
  validate = validator.New()
  type User struct {
    ID     int64  `json:"id" validate:"gt=0"`
    Name   string `json:"name" validate:"required"`
    Gender string `json:"gender" validate:"required,oneof=man woman"`
    Age    uint8  `json:"age" validate:"required,gte=0,lte=130"`
    Email  string `json:"email" validate:"required,email"`
  }
  user := &User{
    ID:     1,
    Name:   "frank",
    Gender: "boy",
    Age:    180,
    Email:  "gopher@88.com",
  }
  err = validate.Struct(user)
  if err != nil {
    validationErrors := err.(validator.ValidationErrors)
    // output: Key: 'User.Age' Error:Field validation for 'Age' failed on the 'lte' tag
    // fmt.Println(validationErrors)
    fmt.Println(validationErrors.Translate(trans))
    return
  }
}
```

细心的读者可能已经发现，错误输出信息并不友好，错误输出信息中的字段不仅没有使用备用名（首字母小写的字段名），也没有翻译为中文。通过改动代码，使错误输出信息变得友好。

注册一个函数，获取结构体字段的备用名称：

```go
validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
    if name == "-" {
      return "j"
    }
    return name
  })
```


错误信息翻译为中文：

```go
zh := zh.New()
uni = ut.New(zh)
trans, _ := uni.GetTranslator("zh")
_ = zh_translations.RegisterDefaultTranslations(validate, trans)
```


# 标签

通过以上章节的内容，读者应该已经了解到 Validator 是一个基于 tag（标签），实现结构体和单个字段的值验证库。

本章节列举一些比较常用的标签：

在 Go 语言的 `validator` 库中，常用的标签可分为**基础验证**、**格式验证**、**逻辑验证**和**高级验证**四大类。以下是具体分类及说明：

---

### 一、基础验证标签
1. **`required`**
    - 作用：字段必须为非零值（字符串不为空、数字非零、指针非空等）
    - 示例：`validate:"required"`
    - 适用类型：所有类型

2. **`min` / `max`**
    - 作用：约束数值范围或字符串/集合长度
    - 示例：`validate:"min=18,max=100"`（年龄在 18-100 之间）
    - 适用类型：数字、字符串、切片、数组、Map

3. **`len`**
    - 作用：精确匹配字符串长度或集合大小
    - 示例：`validate:"len=11"`（手机号必须为 11 位）

---

### 二、格式验证标签
4. **`email`**
    - 作用：验证邮箱格式
    - 示例：`validate:"email"`

5. **`url`**
    - 作用：验证 URL 格式
    - 示例：`validate:"url"`

6. **`ip` / `ipv4` / `ipv6`**
    - 作用：验证 IP 地址格式
    - 示例：`validate:"ipv4"`

7. **`uuid`**
    - 作用：验证 UUID 格式
    - 示例：`validate:"uuid"`

8. **`datetime`**
    - 作用：验证日期时间格式（需指定格式）
    - 示例：`validate:"datetime=2006-01-02"`（匹配 `YYYY-MM-DD` 格式）

---

### 三、逻辑验证标签
9. **`oneof`**
    - 作用：字段值必须为枚举值之一
    - 示例：`validate:"oneof=male female"`（性别只能是 `male` 或 `female`）

10. **`eqfield` / `nefield`**
    - 作用：跨字段验证（如密码确认）
    - 示例：`validate:"eqfield=Password"`（确认密码需与密码一致）

11. **`gte` / `lte`**
    - 作用：数值大于等于（`gte`）或小于等于（`lte`）指定值
    - 示例：`validate:"gte=0,lte=130"`（年龄在 0-130 之间）

12. **`required_if`**
    - 作用：条件验证（其他字段满足条件时必填）
    - 示例：`validate:"required_if=PaymentMethod alipay"`（当支付方式为支付宝时必填）

---

### 四、高级验证标签
13. **`dive`**
    - 作用：递归验证嵌套结构体或集合内的元素
    - 示例：`validate:"dive"`（验证切片中的每个结构体）

14. **`excludesall` / `excludesrune`**
    - 作用：排除指定字符或 Unicode 字符
    - 示例：`validate:"excludesall=#%&"`（字符串不能包含 `#`、`%`、`&`）

15. **`startswith` / `endswith`**
    - 作用：验证字符串前缀或后缀
    - 示例：`validate:"startswith=ID-"`（字符串必须以 `ID-` 开头）

---

### 完整示例
```go
type UserRequest struct {
    Username  string `validate:"required,min=3,max=20"`        // 必填，长度3-20
    Password  string `validate:"required,min=8"`               // 必填，至少8位
    Email     string `validate:"required,email"`                // 必填且为邮箱格式
    Age       int    `validate:"gte=18,lte=100"`                // 年龄18-100
    Gender    string `validate:"oneof=male female"`             // 性别枚举值
    StartDate string `validate:"datetime=2006-01-02"`          // 日期格式校验
    Addresses []Address `validate:"dive"`                      // 递归验证嵌套结构体
}

type Address struct {
    City    string `validate:"required"`
    ZipCode string `validate:"required,numeric,len=6"`
}
```

---

### 扩展场景
- **自定义标签**：可通过 `RegisterValidation` 自定义规则（如手机号格式）。
- **错误消息翻译**：利用 `ValidationErrors.Translate()` 将错误信息本地化。

如需更完整的标签列表，可参考官方文档：[go-playground/validator](https://github.com/go-playground/validator)。