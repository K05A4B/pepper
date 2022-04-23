package pepper

import "strings"

// 字典树
type TrieNode struct {
	Trusteeship bool

	// PrefixPath  string

	// 子节点
	Next map[string]*TrieNode
	// 此节点的上一个节点
	Prev *TrieNode

	Middleware []*MiddlewareFunc

	pool HandlerPool
}

// 添加处理函数到节点
func (t *TrieNode) Insert(method string, uri string, handler *HandlerFunc) {
	var trusteeship bool

	if uri == "" {
		return
	}

	if uri == "/" {
		t.pool.NewHandler(method, *handler)
		return
	}

	rUri := []rune(uri)
	if rUri[len(rUri)-1] == []rune("/")[0] {
		trusteeship = true
	}

	var currentNode *TrieNode = t

	uri = strings.Trim(uri, "/")
	node := strings.Split(uri, "/")
	nodeNum := len(node)

	for i := 0; i < nodeNum; i++ {
		item := node[i]

		if item == "" {
			continue
		}

		prevRouter := currentNode
		currentNode = t.setChildNode(currentNode, item)
		currentNode.Prev = prevRouter

		currentNode.Trusteeship = trusteeship
		// currentNode.PrefixPath = uri

		if nodeNum == (i + 1) {
			currentNode.pool.NewHandler(method, *handler)
		}
	}
}

func (t *TrieNode) InsertByGroup(uri string, group *Group) {

	if uri == "/" {
		return
	}

	var currentNode *TrieNode = t

	uri = strings.Trim(uri, "/")
	node := strings.Split(uri, "/")
	nodeNum := len(node)

	for i := 0; i < nodeNum; i++ {
		item := node[i]

		if item == "" {
			continue
		}

		prevRouter := currentNode
		currentNode = t.setChildNode(currentNode, item)
		currentNode.Prev = prevRouter

		if nodeNum == (i + 1) {
			currentNode.pool = group.Trie.pool
			currentNode.Next = group.Trie.Next
			currentNode.Middleware = group.Trie.Middleware
			currentNode.Trusteeship = false
		}
	}
}

// 通过URI查找字典树 返回对应节点
// 如果没有找到则返回nil
func (t *TrieNode) SearchNode(uri string, callback func(*TrieNode)) *TrieNode {
	var currentNode *TrieNode = t

	uri = strings.Trim(uri, "/")
	node := strings.Split(uri, "/")
	nodeNum := len(node)

	if callback == nil {
		callback = func(tn *TrieNode) {}
	}

	for i := 0; i < nodeNum; i++ {
		item := node[i]

		if item == "" {
			continue
		}

		nextNode, ok := currentNode.Next[item]
		if !ok {
			return nil
		}

		callback(nextNode)

		if currentNode.Trusteeship {
			return currentNode
		}

		if nodeNum == (i + 1) {
			return nextNode
		}

		currentNode = nextNode
	}

	return currentNode
}

// 通过URI查找树并返回 处理池 (*HandlePool)
func (t *TrieNode) Search(uri string) *HandlerPool {
	node := t.SearchNode(uri, nil)
	if node == nil {
		return nil
	}
	return &node.pool
}

// 获取对应处理函数
func (t *TrieNode) GetHandlerFunc(method string, uri string) *HandlerFunc {
	pool := t.Search(uri)
	return pool.GetHandlerFunc(method)
}

// 调用当前节点的对应处理的函数
func (t *TrieNode) CallHandler(res Response, req *Request) {
	f := t.pool.GetHandlerFunc(req.Method)
	if f == nil {
		res.SendErrorPage(404)
		return
	}
	(*f)(res, req)
}

// 设置子节点返回子节点指针
func (t *TrieNode) setChildNode(currentNode *TrieNode, s string) *TrieNode {
	if currentNode == nil {
		return nil
	}

	if currentNode.Next == nil {
		currentNode.Next = make(map[string]*TrieNode)
	}

	if _, ok := currentNode.Next[s]; !ok {
		currentNode.Next[s] = new(TrieNode)
	}

	return currentNode.Next[s]
}
