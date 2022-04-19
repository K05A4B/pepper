package pepper

import (
	"net/http"
	"os"
)

type Pepper struct {

	// 调试模式
	DebugMode bool

	KeyFile string // 密钥文件路径 填写后会自动开启 HTTPS
	CrtFile string // 证书文件路径 填写后会自动开启 HTTPS

	// 设置错误页
	HttpErrorPages *ErrorPages

	// 中间件
	middleware []*MiddlewareFunc

	// 使用数据结构 树 来保存处理函数一个节点对应一个路径
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
	tree *TreeNode
}

func (p *Pepper) NewHandler(method string, uri string, handler HandlerFunc) {
	if p.tree == nil {
		p.tree = new(TreeNode)
	}

	if uri == "/" {
		p.tree.NewHandler(method, handler)
		return
	}

	p.tree.PushNode(method, uri, handler, true)
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
	p.All(uri+"!:", func(res Response, req *Request) {
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

// 使用一个组
func (p *Pepper) UseGroup(uri string, group *Group) {
	if p.tree == nil {
		p.tree = new(TreeNode)
	}

	if uri == "/" {
		return
	}

	p.tree.PushNodeByGroup(uri, group)
}

// 运行服务器
func (p *Pepper) Run(addr string) error {
	srv := http.NewServeMux()
	srv.HandleFunc("/", p.pepperHandler)

	if p.CrtFile != "" && p.KeyFile != "" {
		return http.ListenAndServeTLS(addr, p.CrtFile, p.KeyFile, srv)
	}

	return http.ListenAndServe(addr, srv)
}

func NewPepper() *Pepper {
	p := &Pepper{}
	if p.HttpErrorPages == nil {
		p.HttpErrorPages = new(ErrorPages)
		p.HttpErrorPages.Other = make(map[int]string)
	}
	return p
}
