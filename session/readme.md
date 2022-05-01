# pepper/session

管理 session

通过内存储存数据
```go
package main

import (
    "http"

    "github.com/K05A4B/pepper"
    "github.com/K05A4B/pepper/session"
)

var mery *Memory

func main(){
    // 创建一个 Session 管理器(内存)
	mery = session.NewMemory(&session.Options{
        // 生命周期(单位: 秒)
		LifeCycle: 60*60*12,
	})

	app := pepper.NewPepper()

	app.DebugMode = true
	
	app.Get("/set", set)
	app.Get("/get", get)

	app.Listen(":8080")
}

func set(res pepper.Response, req *pepper.Request) {
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
}

func get(res pepper.Response, req *pepper.Request) {
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
}
```

通过文件加内存的方式储存数据

```go
package main

import (
    "http"

    "github.com/K05A4B/pepper"
    "github.com/K05A4B/pepper/session"
)

var store *Store

func main(){
    var err error

    // 创建一个 Session 管理器(文件+内存)
	store, err = session.NewStore("./test_dir/session.dat", &session.Options{
		LifeCycle: 60*60*12,
	})
	if err != nil {
		fmt.Println("StoreError:", err)
		return
	}

	app := pepper.NewPepper()

	app.DebugMode = true
	
	app.Get("/set", set)
	app.Get("/get", get)

	app.Listen(":8080")
}

func set(res pepper.Response, req *pepper.Request) {
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
}

func get(res pepper.Response, req *pepper.Request) {
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
}
```