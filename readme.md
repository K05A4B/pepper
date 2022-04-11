# 关于Pepper

Pepper 是一个用go语言编写的轻量HTTP库

# 实现原理
中考结束后补齐

# 安装 Pepper

```shell
go get github.com/kz91/Pepper
```

# 创建一个简单的服务

这样就完成了最简单的服务
访问 [http://127.0.0.1](http://127.0.0.1) 时就会在页面上显示 “Hello World”

```go
package main

import "github.com/kz91/pepper"

func main() {
    // 创建 Pepper 对象
    app := pepper.NewPepper()
    app.All("/", root)
    app.Run(":80")
}

func root(res pepper.Response, req *pepper.Request) {
    res.WriteString("Hello World")
})

```
\* 文档不全中考结束后在补全

[快查](./list.md)

下一章: [Response 对象](./response.md)