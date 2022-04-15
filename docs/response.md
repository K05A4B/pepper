# Response对象是什么?

首先我们要知道 `Response` 对象是干什么的？
`Response` 对象用于动态响应客户端请求，并将动态生成的响应结果返回到客户端浏览器中。使用Response对象可以直接发送信息给浏览器，重定向浏览器到另一个URL或设置cookie的值等



# Response 对象一般出现在哪？

 `HandlerFunc ` 或者 `MiddlewareFunc` 类型的函数里面都有 `Response` 对象



# 该怎么样使用Response对象？


### 向客户端传输数据: `WriteString()` 函数

传输字符串
例如说我们一开始看到的例子 [创建一个简单的服务](./started.md#创建一个简单的服务) 里面用到了 Response的`WriteString()` 函数

```go
func (r **Response*) WriteString(str *string*) (int, error)
```

+ 参数
    - `str` 要传输的字符串数据
    
+ 返回值
    - `int`: 传输的数据大小
    - `error`: 错误信息(且一般此函数不会出错)

----



### 向客户端传输数据: Write() 函数

传输字二进制数据
`Write()` 函数 使用方法和 `WriteString()` 函数使用方法一致，但是传入的是`字节串`

```go
func (r *Response) Write(b []byte) (int, error)
```

+ 参数
    - b 要传输的二进制数据

+ 返回值
    - `int`: 传输的数据大小
    - `error`: 错误信息(如果没有则返回`nil`, 且一般此函数不会出错)

---



### 向客户端传输数据: Json() 函数

传输自动编码的 Json 数据

传入结构体或者集合数据类型会自动转成 json 数据并且发送给客户端

```go
func (r *Response) Json(v interface{}) (err error)
```

+ 参数
  - v : 结构体或者集合(map)
+ 返回值
  - `err error`: 错误信息

---

### 向客户端传输数据: WriteFile() 函数

发送文件到客户端

```go
func (r *Response) WriteFile(file string, bufferSize int) (err error)
```

+ 参数
  - `file string`: 文件路径
  - `bufferSize int`: 缓冲区大小
+ 返回值
  - `err error`: 错误信息

---

### 向客户端传输数据: WriteReader() 函数

io.Reader 接口发送数据到客户端

```go
func (r *Response) WriteReader(reader io.Reader, bufferSize int) error
```

+ 参数
  - `file string`: 文件路径
  - `bufferSize int`: 缓冲区大小
+ 返回值
  - `err error`: 错误信息

---

### 向客户端传输数据: SendErrorPage() 函数

发送错误页

```go
func (r *Response) SendErrorPage(code int) error
```

+ 参数
  - `code int`: 错误码(会通过错误码自动找设置的错误页)
+ 返回值
  - `error`: 错误信息

---



### 向客户端传输数据: Template() 函数

解析模板并且自动发送(因为是封装的 "html/template" 所以模板的语法一致)

```go
func (r *Response) Template(file string, tpl interface{}, fc FuncMap) error
```

+ 参数
  - `file string`: 文件路径
  - `tpl interface{}`: 模板的数据
  - `fc FuncMap`: 模板里的自定义函数
+ 返回值
  - `error`: 错误信息

---



### 设置响应头数据: SetStatusCode() 函数

设置HTTP[状态码](https://www.runoob.com/http/http-status-codes.html)

```go
func (r *Response) SetStatusCode(code int)
```

+ 参数
  - `code int`: 状态码

---



### 设置响应头数据:  SetHeader() 函数

设置响应头

```go
func (r *Response) SetHeader(key, value string)
```

- 参数
  - `key string`: 响应头的键的名称
  - `value string`: 值的数据

----



### 设置响应头数据: SetCookie() 函数

设置Cookie

```go
func (r *Response) SetCookie(opt *http.Cookie)
```

+ 参数
  + `opt http.Cookie`: 传入 http.Cookie 对象

---



### 设置响应头数据: Redirect() 函数

重定向

```go
func (r *Response) Redirect(url string)
```

+ 参数
  - `url string`: 要重定向到的url

---

上一章: [开始使用](./started.md)

下一章: [Request 对象](./request.md)