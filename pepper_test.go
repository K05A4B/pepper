package pepper_test

import (
	"fmt"
	"testing"

	"github.com/kz91/pepper"
)

func TestPepperCore(t *testing.T) {
	app := pepper.NewPepper()
	app.DebugMode = true
	app.All("/test1/child", func(res pepper.Response, req *pepper.Request) {
		res.WriteString("Hello World")
	})

	app.All("/test1", func(res pepper.Response, req *pepper.Request) {
		res.WriteString("Test1")
	})

	app.All("/test2/", func(res pepper.Response, req *pepper.Request) {
		res.WriteString(req.Path)
	})

	app.Static("/static/", "./")
	app.Run(":8080")
}

func TestPepperGroup(t *testing.T) {
	group1 := pepper.NewGroup()

	group1.Use(func(p *pepper.Pepper, res pepper.Response, req *pepper.Request) bool {
		fmt.Println("group1", req.Query("auth"))
		if req.Query("auth") == "true" {
			return true
		} else {
			res.SendErrorPage(403)
			return false
		}
	})

	group1.All("/", func(res pepper.Response, req *pepper.Request) {
		res.WriteString("Welcome to group1!")
	})

	group1.All("/child", func(res pepper.Response, req *pepper.Request) {
		res.WriteString("Welcome to group1.child!")
	})

	app := pepper.NewPepper()
	app.DebugMode = true

	app.ErrorPages.Forbidden = func(res pepper.Response, req *pepper.Request) {
		res.WriteString("<h1>403 Forbidden</h1><p>path: "+req.Path+"</p><p>address: "+req.RemoteAddr+"</p>")
	}

	app.UseGroup("/group1", group1)
	app.Run(":8080")
}
