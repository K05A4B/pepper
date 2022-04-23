package pepper

import (
	"net/http"
	"os"
	"strings"
)

type Pepper struct {
	// 调试模式
	DebugMode bool

	KeyFile string // 密钥文件路径 填写后会自动开启 HTTPS
	CrtFile string // 证书文件路径 填写后会自动开启 HTTPS

	// 错误页
	ErrorPages *ErrorPages

	middleware []*MiddlewareFunc

	// 字典树
	// 通过字典树来保存处理池
	trie *TrieNode
}

// 创建对应请求处理器
func (p *Pepper) NewHandler(method string, uri string, handler HandlerFunc) {
	p.trie.Insert(method, uri, &handler)
}

// 快捷创建请求处理器
func (p *Pepper) All(u string, h HandlerFunc)  { p.NewHandler(METHOD_ALL, u, h) }
func (p *Pepper) Get(u string, h HandlerFunc)  { p.NewHandler(METHOD_GET, u, h) }
func (p *Pepper) Post(u string, h HandlerFunc) { p.NewHandler(METHOD_POST, u, h) }

// 设置静态文件
func (p *Pepper) Static(uri string, dir string) {
	p.All(uri, func(res Response, req *Request) {
		file := dir + strings.TrimPrefix(req.Path, uri)
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

func (p *Pepper) Use(mw ...MiddlewareFunc) {
	for _, f := range mw {
		p.middleware = append(p.middleware, &f)
	}
}

func (p *Pepper) UseGroup(uri string, group *Group) {
	p.trie.InsertByGroup(uri, group)
}

// 运行服务器
func (p *Pepper) Run(addr string) error {
	srv := http.NewServeMux()
	srv.HandleFunc("/", p.handlerPepper)

	if p.CrtFile != "" && p.KeyFile != "" {
		return http.ListenAndServeTLS(addr, p.CrtFile, p.KeyFile, srv)
	}

	return http.ListenAndServe(addr, srv)
}

// 创建 Pepper 对象
func NewPepper() *Pepper {
	return &Pepper{
		trie:       new(TrieNode),
		ErrorPages: NewErrorPages(),
	}
}
