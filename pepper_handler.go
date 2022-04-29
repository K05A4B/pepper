package pepper

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// 处理请求
func (p *Pepper) handlerPepper(res http.ResponseWriter, req *http.Request) {
	response := Response{
		Resp:       res,
		ErrorPages: p.ErrorPages,
	}

	request := &Request{
		Req:           req,
		Body:          req.Body,
		Method:        req.Method,
		Path:          req.URL.Path,
		Proto:         req.Proto,
		Host:          req.Host,
		RemoteAddr:    req.RemoteAddr,
		ContentLength: req.ContentLength,
	}

	response.req = request
	request.res = response

	defer p.panicRecover(response)

	for _, mwf := range p.middleware {
		if !(*mwf)(p, response, request) {
			return
		}
	}

	continueExecute := true

	node := p.trie.SearchNode(request.Path, func(tn *TrieNode) {
		middleware := tn.Middleware
		mwLen := len(middleware)
		if mwLen == 0 {
			return
		}

		for _, mwf := range middleware {
			if !(*mwf)(p, response, request) {
				continueExecute = false
				return
			}
		}
	})

	if !continueExecute {
		return
	}

	if node == nil {
		response.SendErrorPage(404)
		return
	}

	node.CallHandler(response, request)
}

// 恢复处理函数中的 panic 并且显示错误
func (p *Pepper) panicRecover(res Response) {
	if err := recover(); err != nil {

		const size = 64 << 10

		if !p.DebugMode {
			log.Printf("pepper: panic serving: %v\n", err)
			res.SendErrorPage(500)
			return
		}

		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		errContent := fmt.Sprintf("pepper: panic serving: %v\n%s\n", err, buf)

		log.Print(errContent)
		res.WriteString(errContent)
	}
}
