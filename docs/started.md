# 关于Pepper

Pepper 是一个用go语言编写的轻量HTTP框架

# 实现原理
使用[字典树](https://zhuanlan.zhihu.com/p/28891541)的某一个节点来保存对应的处理函数池
一个节点对应一个路径


# 安装 Pepper

```shell
go get github.com/K05A4B/Pepper
```

# 创建一个简单的服务

这样就完成了最简单的服务
访问 [http://127.0.0.1](http://127.0.0.1) 时就会在页面上显示 “Hello World”

```go
package main

import "github.com/K05A4B/pepper"

func main() {
    // 创建 Pepper 对象
    app := pepper.NewPepper()
    app.All("/", root)
    app.Listen(":80")
}

func root(res pepper.Response, req *pepper.Request) {
    res.WriteString("Hello World")
})

```
\* 文档不全中考结束后在补全

[快查](./list.md)

下一章: [Pepper 对象](./pepper.md)