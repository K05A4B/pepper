package pepper

import "strings"

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

// 创建一种请求方式的处理器
func (h *HandlerPool) NewHandler(method string, handler HandlerFunc) {
	if handler == nil {
		return
	}

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

// 获取对应处理函数
func (h *HandlerPool) GetHandlerFunc(method string) *HandlerFunc {
	var callback *HandlerFunc

	if h.All != nil {
		return h.All
	}

	method = strings.ToUpper(method)
	switch method {
	case METHOD_GET:
		callback = h.Get
	case METHOD_POST:
		callback = h.Post
	case METHOD_HEAD:
		callback = h.Head
	case METHOD_PUT:
		callback = h.Put
	case METHOD_DELETE:
		callback = h.Delete
	case METHOD_CONNECT:
		callback = h.Connect
	case METHOD_OPTIONS:
		callback = h.Options
	case METHOD_TRACE:
		callback = h.Trace
	case METHOD_PATCH:
		callback = h.Patch
	}

	return callback
}