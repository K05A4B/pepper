package pepper

import (
	"strings"
)

const (
	METHOD_ALL     = "ALL"
	METHOD_GET     = "GET"
	METHOD_POST    = "POST"
	METHOD_HEAD    = "HEAD"
	METHOD_PUT     = "PUT"
	METHOD_DELETE  = "DELETE"
	METHOD_CONNECT = "CONNECT"
	METHOD_OPTIONS = "OPTIONS"
	METHOD_TRACE   = "TRACE"
	METHOD_PATCH   = "PATCH"
)

// 处理函数
type HandlerFunc func(res Response, req *Request)

// 处理程序池
type HandlerPool struct {
	All     *HandlerFunc
	Get     *HandlerFunc
	Post    *HandlerFunc
	Head    *HandlerFunc
	Put     *HandlerFunc
	Delete  *HandlerFunc
	Connect *HandlerFunc
	Options *HandlerFunc
	Trace   *HandlerFunc
	Patch   *HandlerFunc
}

// 创建一种请求方式
func (h *HandlerPool) NewHandler(method string, handler HandlerFunc) {
	switch strings.ToUpper(method) {
	case METHOD_GET:
		h.Get = &handler
	case METHOD_POST:
		h.Post = &handler
	case METHOD_HEAD:
		h.Head = &handler
	case METHOD_PUT:
		h.Put = &handler
	case METHOD_DELETE:
		h.Delete = &handler
	case METHOD_CONNECT:
		h.Connect = &handler
	case METHOD_OPTIONS:
		h.Options = &handler
	case METHOD_TRACE:
		h.Trace = &handler
	case METHOD_PATCH:
		h.Patch = &handler
	// 默认为任意请求方式
	default:
		h.All = &handler
	}
}

// 调用对应处理函数
func (h *HandlerPool) Call(method string, res Response, req *Request, p *Pepper) {

	var callBackFunction *HandlerFunc

	method = strings.ToUpper(method)
	switch method {
	case METHOD_GET:
		callBackFunction = h.Get
	case METHOD_POST:
		callBackFunction = h.Post
	case METHOD_HEAD:
		callBackFunction = h.Head
	case METHOD_PUT:
		callBackFunction = h.Put
	case METHOD_DELETE:
		callBackFunction = h.Delete
	case METHOD_CONNECT:
		callBackFunction = h.Connect
	case METHOD_OPTIONS:
		callBackFunction = h.Options
	case METHOD_TRACE:
		callBackFunction = h.Trace
	case METHOD_PATCH:
		callBackFunction = h.Patch
	}

	// 如果有ALL函数就调用ALL函数
	All := h.All
	if All != nil {
		callBackFunction = All
	}

	// 判断处理函数是否为 nil
	if callBackFunction != nil {
		// 调用函数
		(*callBackFunction)(res, req)
	} else {
		res.SendErrorPage(403)
	}
}
