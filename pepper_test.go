package pepper_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/K05A4B/pepper"
	"github.com/K05A4B/pepper/middleware/log"
	"github.com/K05A4B/pepper/session"
	_ "github.com/mattn/go-sqlite3"
)

// 测试框架核心功能
func TestPepperCore(t *testing.T) {

	// 创建 Pepper 对象
	app := pepper.NewPepper()

	// 调试模式
	app.DebugMode = true

	// 使用中间件
	app.Use(log.NewMiddleware("./test_dir/log"))

	// 创建节点
	app.All("/test1/child", func(res pepper.Response, req *pepper.Request) {
		res.WriteString("Hello World")
	})

	app.All("/test1", func(res pepper.Response, req *pepper.Request) {
		res.WriteString("Test1")
	})

	app.All("/test2/", func(res pepper.Response, req *pepper.Request) {
		res.WriteString(req.Path)
	})

	// 设置静态目录
	app.Static("/static/", "./")

	// 运行服务器
	app.Listen(":8080")
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
		res.WriteString("<h1>403 Forbidden</h1><p>path: " + req.Path + "</p><p>address: " + req.RemoteAddr + "</p>")
	}

	app.UseGroup("/group1", group1)
	app.Listen(":8080")
}

func TestSessionByMemory(t *testing.T) {

	// 创建一个 Session 管理器
	mery := session.NewMemory(&session.Options{
		LifeCycle: 60*60*12,
	})

	app := pepper.NewPepper()

	app.DebugMode = true
	
	app.Get("/set", func(res pepper.Response, req *pepper.Request) {
		name := req.Query("name")
		
		id := req.GetCookieValue("SessionID")
		
		if !mery.Exist(id) {
			id, _ = mery.Create()
			res.SetCookie(&http.Cookie{
				Name: "SessionID",
				Value: id,
				MaxAge: 60*60*12,
			})
		}

		sess := mery.Get(id)
		sess.Set("name", name)
	})

	app.Get("/get", func(res pepper.Response, req *pepper.Request) {
		id := req.GetCookieValue("SessionID")
		
		if !mery.Exist(id) {
			res.SendErrorPage(403)
			return
		}

		sess := mery.Get(id)
		name := sess.Get("name")
		if name == nil {
			res.WriteString("nil")
			return
		}

		res.WriteString(name.(string))
	})

	app.Listen(":8080")
}

func TestSessionByStore(t *testing.T) {

	// 创建一个 Session 管理器
	store, err := session.NewStore("./test_dir/session.dat", &session.Options{
		LifeCycle: 12,
	})
	if err != nil {
		fmt.Println("StoreError:", err)
		return
	}

	defer store.Close()

	app := pepper.NewPepper()

	app.DebugMode = true
	
	app.Get("/set", func(res pepper.Response, req *pepper.Request) {
		name := req.Query("name")
		
		id := req.GetCookieValue("SessionID")
		
		if !store.Exist(id) {
			id, _ = store.Create()
			res.SetCookie(&http.Cookie{
				Name: "SessionID",
				Value: id,
				MaxAge: 60*60*12,
			})
		}
				
		sess := store.Get(id)
		sess.Set("name", name)
		if err := store.Save();err != nil {
			res.WriteString(err.Error())
		}
	})

	app.Get("/get", func(res pepper.Response, req *pepper.Request) {
		id := req.GetCookieValue("SessionID")
		
		if !store.Exist(id) {
			res.SendErrorPage(403)
			return
		}

		sess := store.Get(id)
		name := sess.Get("name")
		if name == nil {
			res.WriteString("nil")
			return
		}

		res.WriteString(name.(string))
	})

	app.Listen(":8080")
}

func TestSessionByDataBase(t *testing.T) {

	// 创建一个 Session 管理器
	store, err := session.NewDataBase("sqlite3", "./test_dir/session.db", &session.Options{
		LifeCycle: 60*60*12,
		CleanInterval: 15,
	})

	if err != nil {
		fmt.Println("StoreError:", err)
		return
	}

	app := pepper.NewPepper()

	app.DebugMode = true
	
	app.Get("/set", func(res pepper.Response, req *pepper.Request) {
		name := req.Query("name")
		
		id := req.GetCookieValue("SessionID")
		
		if !store.Exist(id) {
			id, _ = store.Create()
			res.SetCookie(&http.Cookie{
				Name: "SessionID",
				Value: id,
				MaxAge: 60*60*12,
			})
		}
				
		sess := store.Get(id)
		sess.Set("name", name)

		if err := store.Push();err != nil {
			res.WriteString(fmt.Sprint(err, name))
		}
	})

	app.Get("/get", func(res pepper.Response, req *pepper.Request) {
		id := req.GetCookieValue("SessionID")
		
		if !store.Exist(id) {
			res.SendErrorPage(403)
			return
		}

		if err := store.Pull(); err != nil {
			res.WriteString(err.Error())
		}

		sess := store.Get(id)
		name := sess.Get("name")
		if name == nil {
			res.WriteString("nil")
			return
		}

		res.WriteString(name.(string))
	})

	app.Get("/remove", func(res pepper.Response, req *pepper.Request) {
		id := req.GetCookieValue("SessionID")

		if err := store.Pull(); err != nil {
			res.WriteString(fmt.Sprint("拉取数据库数据失败：", err))
			return 
		}
		
		if !store.Exist(id) {
			res.SendErrorPage(403)
			return
		}

		sess := store.Get(id)

		sess.Remove("name")

		if err := store.Push(); err != nil {
			res.WriteString(fmt.Sprint("推送数据到数据库数据失败：", err))
			return 
		}
	})

	app.Get("/empty", func(res pepper.Response, req *pepper.Request) {
		id := req.GetCookieValue("SessionID")

		if !store.Exist(id) {
			res.SendErrorPage(403)
			return
		}

		sess := store.Get(id)
		sess.Empty()

		if err := store.Push(); err != nil {
			res.WriteString(fmt.Sprint("推送数据到数据库数据失败：", err))
			return 
		}
	})
	app.Listen(":8080")
}
