package pepper_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/kz91/pepper"
	"github.com/kz91/pepper/middleware/log"
	"github.com/kz91/pepper/upload"
)

func TestPepperGroup(t *testing.T) {
	g2 := pepper.NewGroup()
	g2.All("/test2", func(res pepper.Response, req *pepper.Request) {
		res.WriteString("Wellcome to g2 test1 page!")
	})

	g1 := pepper.NewGroup()
	g1.UseGroup("/", g2)
	g1.All("/test1", func(res pepper.Response, req *pepper.Request) {
		fmt.Println("runing test1")
		res.WriteString("Wellcome to g1 test1 page!")
	})

	app := pepper.NewPepper()
	app.Get("/group1", func(res pepper.Response, req *pepper.Request) {
		res.WriteString("Wellcome to group1 page")
	})
	app.UseGroup("/group1/", g1)
	app.All("/group2/group3/a!:", func(res pepper.Response, req *pepper.Request) {
		fmt.Println("path: ", req.TrimPath)
		res.WriteString(req.TrimPath)
	})

	app.All("/group2", func(res pepper.Response, req *pepper.Request) {
		res.WriteString("Group2 Get")
	})

	fmt.Println("Running")
	app.Run("127.0.0.1:8080")
}

func TestPepperMiddleware(t *testing.T) {
	app := pepper.NewPepper()

	// 调用中间件
	app.Use(log.NewMiddleware("./test/log"))

	app.All("/", func(res pepper.Response, req *pepper.Request) {
		res.Json(map[string]string{
			"msg":  "首页",
			"type": "root",
		})
	})

	app.All("/test1", func(res pepper.Response, req *pepper.Request) {
		res.Json(map[string]string{
			"msg":  "test1",
			"type": "page",
		})
	})

	app.All("/test2", func(res pepper.Response, req *pepper.Request) {
		res.Json(map[string]string{
			"msg":  "test2",
			"type": "page",
		})
	})

	fmt.Println("Running")
	app.Run(":8080")
}

func TestPepperStatic(t *testing.T) {
	app := pepper.NewPepper()
	app.HttpErrorPages.NotFound = "./404.html"
	app.Post("/", func(res pepper.Response, req *pepper.Request) {
		// fmt.Println("test")
		fmt.Println(req.Query("test"), "test")
		fmt.Println(req.GetFormStringValue("test1"))
		res.WriteFile("./pepper_test.go", 5120)
	})

	app.Get("/1", func(res pepper.Response, req *pepper.Request) {
		res.Redirect("/static/pepper.go")
	})
	app.Static("/static", "./")
	fmt.Println("Running")
	app.Run("127.0.0.1:8080")
}

func TestPepperUpload(t *testing.T) {
	app := pepper.NewPepper()
	app.HttpErrorPages.NotFound = "./404.html"
	app.Post("/", func(res pepper.Response, req *pepper.Request) {
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

		r := &upload.Rule{
			MaxSize: 1024 * 1024 * 30,
			// MaxNumber: 2,
		}

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
			file, err := f.Next()
			if err == io.EOF {
				break
			}

			if err != nil {
				result.Success = false
				result.Error = err.Error()
				break
			}

			_, err = file.Receive("./test", 1024)
			var recvError string
			if err != nil {
				recvError = err.Error()
			}

			fmt.Println("receive", file.Name)

			result.Result = append(result.Result, FileRes{
				Error:   recvError,
				Success: err == nil,
				Name:    file.Name,
				Msg:     file.Path,
			})
		}

		res.Json(result)
	})
	fmt.Println("Running")
	app.Run("127.0.0.1:8080")
}


func TestPepperTpl(t *testing.T) {
	app := pepper.NewPepper()
	app.HttpErrorPages.NotFound = "./404.html"
	app.All("/", func(res pepper.Response, req *pepper.Request) {
		res.Template("./test/test.html", map[string]string{
			"name":"xiaoming",
		}, pepper.FuncMap{
			"get_name": func() string {
				return "xiaoming"
			},
		})
	})
	fmt.Println("Running")
	app.Run("127.0.0.1:8080")
}