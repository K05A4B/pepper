package pepper

import (
	"net/http"
	"os"
	"strings"
)

type Pepper struct {
	KeyFile string // 密钥文件路径 填写后会自动开启 HTTPS
	CrtFile string // 证书文件路径 填写后会自动开启 HTTPS

	// 设置错误页
	HttpErrorPages *ErrorPages

	// 中间件
	middleware []*MiddlewareFunc

	// 通过这个属性来储存所有路由
	// 在内存中储存的结果差不多是这样：
	// / -
	//   |_ /router1
	//   |      |  |
	//   |      |  |
	//   |      |  /child1
	//   |      |
	//   |      |______ /child2
	//   |_ /router2
	//          |  |
	//          |  |
	//          |  /child1
	//          |
	//          |———— /child2
	//                      |
	//                      |
	//                      /subchild
	router *Router
}

func (p *Pepper) NewHandler(method string, uri string, handler HandlerFunc) {
	if p.router == nil {
		p.router = new(Router)
	}

	if uri == "/" {
		p.router.NewHandler(method, handler)
		return
	}

	p.router.PushRouter(method, uri, handler, true)
}

func (p *Pepper) All(uri string, handle HandlerFunc) {
	p.NewHandler(METHOD_ALL, uri, handle)
}

func (p *Pepper) Get(uri string, handle HandlerFunc) {
	p.NewHandler(METHOD_GET, uri, handle)
}

func (p *Pepper) Post(uri string, handle HandlerFunc) {
	p.NewHandler(METHOD_POST, uri, handle)
}

func (p *Pepper) Static(uri string, dir string) {
	p.All(uri + "!:", func(res Response, req *Request) {
		file := dir + req.TrimPath
		info, err := os.Stat(file)		
		if err != nil {
			res.SendErrorPage(404)
			return
		}

		if info.IsDir() {
			res.SendErrorPage(403)
			return
		}

		res.WriteFile(file, 2150)
	})
}

// 使用中间件
func (p *Pepper) Use(middleware ...MiddlewareFunc) {
	for i := 0; i < len(middleware); i++ {
		p.middleware = append(p.middleware, &middleware[i])
	}
}

func (p *Pepper) UseGroup(uri string, group *Group) {
	p.router.PushRouterByGroup(uri, group)
}

// 运行服务器
func (p *Pepper) Run(addr string) {
	srv := http.NewServeMux()
	srv.HandleFunc("/", p.pepperHandler)

	if p.CrtFile != "" && p.KeyFile != "" {
		http.ListenAndServeTLS(addr, p.CrtFile, p.KeyFile, srv)
	}

	http.ListenAndServe(addr, srv)
}

// 处理请求
func (p *Pepper) pepperHandler(res http.ResponseWriter, req *http.Request) {
	response := Response{
		Resp: res,
		ErrorPages: p.HttpErrorPages,
	}

	request := &Request{
		Req:    req,
		Method: req.Method,
		Path:   req.URL.Path,
		Proto:  req.Proto,
		Host:   req.Host,
	}

	middlewareLen := len(p.middleware)

	for i := 0; i < middlewareLen; i++ {
		// 调用中间件
		(*p.middleware[i])(p, response, request)
	}

	p.callHandler(res, req, response, request)
}

func (p *Pepper) callHandler(res http.ResponseWriter, req *http.Request, response Response, request *Request) {

	// 判断路由是否存在
	if p.router == nil {
		p.router = new(Router)
	}

	method := req.Method
	URI := req.URL.Path
	if URI == "/" {
		p.router.Call(method, response, request, p)
		return
	}

	URI = strings.TrimSuffix(URI, "/")
	path := strings.Split(URI, "/")
	pathLen := len(path)
	currentRouter := p.router

	for i := 0; i < pathLen; i++ {
		pathItem := path[i]

		if pathItem == "" {
			continue
		}

		if currentRouter.Trusteeship {
			request.TrimPath = strings.TrimPrefix(request.Path, currentRouter.PrefixPath)
			currentRouter.Call(method, response, request, p)
			return
		}

		nextRouter, ok := currentRouter.Next[pathItem]
		if !ok {
			p.HttpErrorPages.SendPage(404, response)
			return
		}

		if pathLen == (i + 1) {
			nextRouter.Call(method, response, request, p)
		}

		currentRouter = nextRouter
	}
}

func NewPepper() *Pepper {
	p := &Pepper{}
	if p.HttpErrorPages == nil {
		p.HttpErrorPages = new(ErrorPages)
		p.HttpErrorPages.Other = make(map[int]string)
	}
	return p
}
