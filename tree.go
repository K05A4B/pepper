package pepper

import "strings"

// 树
type TreeNode struct {
	HandlerPool
	Trusteeship bool
	PrefixPath  string
	Next        map[string]*TreeNode // 路由表
	Prev        *TreeNode            // 上一个路由对象
}

// 推送新路由
func (t *TreeNode) PushNode(method string, uri string, handler HandlerFunc, allowTrusteeship bool) {
	uri = strings.TrimRight(uri, "/")
	path := strings.Split(uri, "/")
	pathLen := len(path)

	var currentNode *TreeNode = t

	if uri == "/" {
		t.NewHandler(method, handler)
		return
	}

	for i := 0; i < pathLen; i++ {
		pathItem := path[i]

		if pathItem == "" {
			continue
		}

		if currentNode.Next == nil {
			currentNode.Next = make(map[string]*TreeNode)
		}

		if _, ok := currentNode.Next[pathItem]; !ok {
			currentNode.Next[pathItem] = new(TreeNode)
		}

		prevRouter := currentNode
		currentNode = currentNode.Next[pathItem]
		currentNode.Prev = prevRouter

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
					delete(currentNode.Prev.Next, pathItem+suffix)

					// 获取去除后缀的路由名字
					name := strings.TrimSuffix(pathItem, suffix)

					// 创建新路由命名为去除后缀的名字
					currentNode.Prev.Next[name] = new(TreeNode)
					currentNode = currentNode.Prev.Next[name]

					// 设置path前缀
					currentNode.PrefixPath = strings.TrimSuffix(uri, suffix)

					// 启用托管
					currentNode.Trusteeship = true
				}
			}

			currentNode.NewHandler(method, handler)
		}
	}
}

// 推送没有处理函数的路由表
func (t *TreeNode) PushNodeByGroup(uri string, g *Group) {
	uri = strings.TrimRight(uri, "/")
	path := strings.Split(uri, "/")
	pathLen := len(path)

	var currentNode *TreeNode = t

	if uri == "/" {
		for k, v := range g.Node.Next {
			_, ok := t.Next[k]
			if !ok {
				t.Next[k] = v
			}
		}
		return
	}

	for i := 0; i < pathLen; i++ {
		pathItem := path[i]

		if pathItem == "" {
			continue
		}

		if currentNode == nil {
			currentNode = new(TreeNode)
		}

		if currentNode.Next == nil {
			currentNode.Next = make(map[string]*TreeNode)
		}

		if _, ok := currentNode.Next[pathItem]; !ok {
			currentNode.Next[pathItem] = new(TreeNode)
		}

		prevNode := currentNode
		currentNode = currentNode.Next[pathItem]
		currentNode.Prev = prevNode
	}

	currentNode.Next = g.Node.Next
}
