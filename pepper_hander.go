package pepper

import (
	"net/http"
	"strings"
)

// 处理请求
func (p *Pepper) pepperHandler(res http.ResponseWriter, req *http.Request) {
	response := Response{
		Resp:       res,
		ErrorPages: p.HttpErrorPages,
	}

	request := &Request{
		Req:        req,
		Method:     req.Method,
		Path:       req.URL.Path,
		Proto:      req.Proto,
		Host:       req.Host,
		RemoteAddr: req.RemoteAddr,
	}

	middlewareLen := len(p.middleware)

	for i := 0; i < middlewareLen; i++ {
		// 调用中间件
		if !(*p.middleware[i])(p, response, request) {
			return
		}
	}

	p.callHandler(res, req, response, request)
}

func (p *Pepper) callHandler(res http.ResponseWriter, req *http.Request, response Response, request *Request) {

	// 判断路由是否存在
	if p.tree == nil {
		p.tree = new(TreeNode)
	}

	method := req.Method
	URI := req.URL.Path
	if URI == "/" {
		p.tree.Call(method, response, request, p)
		return
	}

	URI = strings.TrimSuffix(URI, "/")
	path := strings.Split(URI, "/")
	pathLen := len(path)
	currentNode := p.tree

	for i := 0; i < pathLen; i++ {
		pathItem := path[i]

		if pathItem == "" {
			continue
		}

		if currentNode.Trusteeship {
			request.TrimPath = strings.TrimPrefix(request.Path, currentNode.PrefixPath)
			currentNode.Call(method, response, request, p)
			return
		}

		nextNode, ok := currentNode.Next[pathItem]
		if !ok {
			p.HttpErrorPages.SendPage(404, response)
			return
		}

		if pathLen == (i + 1) {
			nextNode.Call(method, response, request, p)
		}

		currentNode = nextNode
	}
}