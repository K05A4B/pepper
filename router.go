package pepper

import (
	"strings"
)

/* 路由 */
// 继承至 Handler对象，着所有请求方式的处理函数都在这里面
type Router struct {
	HandlerPool
	Trusteeship bool
	PrefixPath  string
	Next        map[string]*Router // 路由表
	Prev        *Router            // 上一个路由对象
}

// 推送新路由
func (r *Router) PushRouter(method string, uri string, handler HandlerFunc, allowTrusteeship bool) {
	uri = strings.TrimRight(uri, "/")
	path := strings.Split(uri, "/")
	pathLen := len(path)

	var currentRouter *Router = r

	if uri == "/" {
		r.NewHandler(method, handler)
		return
	}

	for i := 0; i < pathLen; i++ {
		pathItem := path[i]

		if pathItem == "" {
			continue
		}

		if currentRouter.Next == nil {
			currentRouter.Next = make(map[string]*Router)
		}

		if _, ok := currentRouter.Next[pathItem]; !ok {
			currentRouter.Next[pathItem] = new(Router)
		}

		prevRouter := currentRouter
		currentRouter = currentRouter.Next[pathItem]
		currentRouter.Prev = prevRouter

		if pathLen == (i + 1) {
			// 后缀
			suffix := "!:"

			runePathItem := []rune(pathItem)

			itemLen := len(runePathItem)
			if itemLen > 2 {
				// 这个路由名字的后缀
				itemSuffix := string(runePathItem[itemLen-2:])

				if itemSuffix == suffix && allowTrusteeship {
					// 删除自动创建的路由
					delete(currentRouter.Prev.Next, pathItem+suffix)

					// 获取去除后缀的路由名字
					name := strings.TrimSuffix(pathItem, suffix)

					// 创建新路由命名为去除后缀的名字
					currentRouter.Prev.Next[name] = new(Router)
					currentRouter = currentRouter.Prev.Next[name]

					// 设置path前缀
					currentRouter.PrefixPath = strings.TrimSuffix(uri, suffix)

					// 启用托管
					currentRouter.Trusteeship = true
				}
			}

			currentRouter.NewHandler(method, handler)
		}
	}
}

// 推送没有处理函数的路由表
func (r *Router) PushRouterByGroup(uri string, g *Group) {
	uri = strings.TrimRight(uri, "/")
	path := strings.Split(uri, "/")
	pathLen := len(path)

	var currentRouter *Router = r

	if uri == "/" {
		for k, v := range g.Router.Next {
			_, ok := r.Next[k]
			if !ok {
				r.Next[k] = v
			}
		}
		return
	}

	for i := 0; i < pathLen; i++ {
		pathItem := path[i]

		if pathItem == "" {
			continue
		}

		if currentRouter == nil {
			currentRouter = new(Router)
		}

		if currentRouter.Next == nil {
			currentRouter.Next = make(map[string]*Router)
		}

		if _, ok := currentRouter.Next[pathItem]; !ok {
			currentRouter.Next[pathItem] = new(Router)
		}

		prevRouter := currentRouter
		currentRouter = currentRouter.Next[pathItem]
		currentRouter.Prev = prevRouter
	}

	currentRouter.Next = g.Router.Next
}
