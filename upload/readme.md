# pepper/upload

接收文件需要创建接收器

简单例子
```go
package main

import (
    "fmt"

    "github.com/kz91/pepper"
    "github.com/kz91/pepper/upload"
)

func main() { 
    app := pepper.NewPepper()
	
    app.DebugMode = true

	app.Post("/", func(res pepper.Response, req *pepper.Request) {
		// 创建接收器
		// 接收 表单名称 为 files1 的文件
		f, err := upload.NewReceive(req, "files1", nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			// 接收 一个 文件
			// 将放到 upload 里
			// 并且设置接收缓冲区大小为 2048 (2KB)
			file, err := f.Receive("./upload", 2048)
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Println(err)
				return
			}

            // 文件保存的路径
			fmt.Println(file.Path)
			// 文件可信的 MIME
			fmt.Println(file.Mime)
			// 文件不可信的 MIME (客户端提供的)
			fmt.Println(file.NotTrustworthyMime)
			// 文件大小
			fmt.Println(file.Size)
			// 原始文件名
			fmt.Println(file.Name)
		}
	})
	app.Run(":1592")
```

小例子

```go
package main

import (
    "fmt"

    "github.com/kz91/pepper"
    "github.com/kz91/pepper/upload"
)

func main() { 
    type FileRes struct {
        Success bool   `json:"success"`
        Error   string `json:"err"`
        Msg     string `json:"msg"`
        Name    string `json:"name"`
    }

    type Result struct {
        Code    int       `json:"code"`
        Success bool      `json:"success"`
        Result  []FileRes `json:"result"`
        Error   string    `json:"err"`
		}

    app := pepper.NewPepper()
	
    app.DebugMode = true

	app.Post("/", func(res pepper.Response, req *pepper.Request) {
		// 创建上传规则
		r := &upload.Rule{
			// 可以以上传的最大大小 (单位: 字节)
			MaxSize: 1024 * 1024 * 30,
			// 可以以上传的最小大小 (单位: 字节)
			MinSize: 1024 * 1024 * 3,
			// 可以以上传几个文件
			MaxNumber: 2,
		}
		// 设置允许的文件 mime 类型
		r.Mime.Append("text/html")
		r.Mime.Append("audio/mpeg")

		f, err := upload.NewReceive(req, "files1", r)
		if err != nil {
			fmt.Println(err)
		}

		result := Result{
			Code:    200,
			Success: true,
			Result:  []FileRes{},
		}

		for {
			var ErrorText string
			file, err := f.Receive("./upload", 2048)
			if err == io.EOF {
				break
			}

			result.Code = 0
			result.Success = true

			if err != nil {
				result.Code = 500
				result.Success = false
				ErrorText = err.Error()
			}

            fmt.Println("receive", file.Name)
            fmt.Println("Content-Type", file.Mime)
            result.Result = append(result.Result, FileRes{
                Error:   ErrorText,
                Success: err == nil,
                Name:    file.Name,
                Msg:     file.Path,
            })
		}

		res.Json(result)
	})
	app.Run(":1592")

}
```