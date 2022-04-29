# Pepper对象

通过 pepper.NewPepper() 创建一个pepper对象
可以通过这个对象 运行服务、创建处理、使用中间件、错误页的设置 等功能

## 创建处理

可以使用 pepper 对象的 NewHandler、All、Get、Post 函数来创建处理

```go
package main

import "github.com/K05A4B/pepper"

func main() {
    app := pepper.NewPepper()

    // 处理根目录的任意请求
    // 注意：任意请求的优先级最高
    // app.All("/", handlerGet)

    // 处理根目录的GET请求
    app.Get("/", handlerGet)

    // 处理根目录的POST请求
    app.Post("/", handlerPost)

    // 处理根目录的 HEAD 请求
    // 创建一种请求方式，有以下常量, 第一个参数是以 METHOD_ 开头的常量 接下来与Get、Post、All 一样
    // METHOD_ALL
	// METHOD_GET
	// METHOD_POST
	// METHOD_HEAD
	// METHOD_PUT 
	// METHOD_DELETE
	// METHOD_CONNECT
	// METHOD_OPTIONS
	// METHOD_TRACE
	// METHOD_PATCH
    app.NewHandler(METHOD_HEAD, "/", handlerHead)

    // 这个时候服务并没有启动
    // 在 pepper 对象中有一个叫做 Listen 的函数 在里面传入地址即可
    app.Listen(":1592") // 表示监听 1592 端口
}

func handlerGet(res pepper.Response, req *pepper.Request) {
    res.WriteString("method: GET")
}

func handlerPost(res pepper.Response, req *pepper.Request) {
    res.WriteString("method: Post")
}

func handlerHead(res pepper.Response, req *pepper.Request) {
    res.WriteString("method: Head")
}
```

## 中间件

使用中间件的方法很简单 pepper 对象中有一个叫做 Use 的函数在里面传入中间件即可

例子
```go
package main

import (
    "github.com/K05A4B/pepper"
    mwlog "github.com/K05A4B/pepper/middleware/log" 
)

func main() {
    app := pepper.NewPepper()

    // 这里以 pepper 框架自带的日志中间件为例
    // 中间件参数里填写日志保存的位置
    app.Use(mwlog.NewMiddleware("./log"))

    // 后面的 !: “/” 目录后面所有路径都由 handler 函数接管
    app.All("/!:", handler)

    app.Listen(":1592")
}

func hander(res pepper.Response, req *pepper.Request) {
    res.WriteString(res.TrimPath)
}
```

接下来是中间件的创建

例子

```go
package main

import (
    "github.com/K05A4B/pepper"
    "fmt" 
)

func main() {
    app := pepper.NewPepper()

    // 这里以 pepper 框架自带的日志中间件为例
    // 中间件参数里填写日志保存的位置
    app.Use(middleware)

    // 后面的 !: “/” 目录后面所有路径都由 handler 函数接管
    app.All("/!:", handler)

    app.Listen(":1592")
}

func hander(res pepper.Response, req *pepper.Request) {
    res.WriteString(res.TrimPath)
}

// 返回的布尔值代表是否继续执行
// true 为继续执行 false 结束执行
func middleware(p *Pepper, res Response, req *Request) bool {
    // 每次执行都会打印客户端地址
    fmt.Println("Address:", res.RemoteAddr)

    // 停止处理
    // return false

    //  继续处理
    return true
}
```

# 静态目录的设置

设置静态目录也很简单 调用 Pepper 对象中的 Static 函数即可

例子

```go
package main

import "github.com/K05A4B/pepper"

func main() {
    app := pepper.NewPepper()

    // 这个的意思是将 ./web/static 内的所有内容都映射在 /static 下
    app.Static("/static", "./web/static")

    app.Listen(":1592")
}
```

# 使用组

可以通过组来实现模块化接下来会演示怎么样使用、创建组

```go
package main

import "github.com/K05A4B/pepper"

func main() {
    // 创建一个组
    group1 := pepper.NewGroup()
    group1.All("/test", handler)

    // 访问 /group1 时执行的是下面的函数
    group1.All("/", func(res pepper.Response, req *pepper.Request) {
        res.WriteString("不是预期结果吧?")
    })

    app := pepper.NewPepper()

    // 使用group1这个组
    // 访问 /group1/test 就可以访问到group1的test处理函数
    app.UseGroup("/group1", group1)

    // 访问 /group1 时可能无法执行下面的函数，因为访问的是 group1这个组 的根目录
    app.All("/group1", func(res pepper.Response, req *pepper.Request) {
        res.WriteString("group1")
    })

    app.Listen(":1592")
}

func hander(res pepper.Response, req *pepper.Request) {
    res.WriteString("Test")
}
```

# 错误页的设置与使用

在 pepper 对象中有 一个叫做 HttpErrorPages 的对象

例子
```go
package main

import (
    "github.com/K05A4B/pepper"
)

func main() {
    app := pepper.NewPepper()

    // 设置 404 页面
    app.HttpErrorPages.NotFound = "./ErrorPages/404.html"

    // 上面的效果与下面的一样
    // app.HttpErrorPages.Other[404] = "./ErrorPages/404.html"

    app.All("/404", handler)

    app.Listen(":1592")
}

func hander(res pepper.Response, req *pepper.Request) {
    // 手动返回404页面
    res.SendErrorPage(404)
}
```

# 给服务加证书

通过给服务加证书来启动 HTTPS 

```go
package main

import (
    "github.com/K05A4B/pepper"
)

func main() {
    app := pepper.NewPepper()

    // 设置 密钥文件路径
    app.KeyFile = "xxx.key"

    // 设置证书文件路径
    app.CrtFile = "xxx.crt"

    app.All("/", handler)

    app.Listen(":1592")
}

func hander(res pepper.Response, req *pepper.Request) {
    res.WriteString("Hello World")
}
```

只要证书文件和密钥文件没有问题那么 直接访问[https://127.0.0.1:1592](https://127.0.0.1:1592)是没有问题的

# 调试模式

通过 pepper 对象的 DebugMode 来启动调试模式

```go
package main

import (
    "github.com/K05A4B/pepper"
)

type Test struct {
    Test string
}

func main() {
    app := pepper.NewPepper()

    // 开启调试模式(调试模式默认关闭，且中间件可以修改调试模式的状态)
    app.DebugMode = true

    app.All("/", handler)

    app.Listen(":1592")
}

func hander(res pepper.Response, req *pepper.Request) {
    var t *Test

    // 这样写肯定会报错
    // 在调试模式下会将错误信息传输到客户端
    // 在非调试模式下会将 “服务器内部错误” 的页面发送到客户端 且控制台不会输出详细错误信息
    t.Test = "test"
}
```