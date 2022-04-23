package pepper

type Group struct {
	Trie *TrieNode
}

// 创建对应请求处理器
func (g *Group) NewHandler(method string, uri string, handler HandlerFunc) {
	g.Trie.Insert(method, uri, &handler)
}

// 快捷创建请求处理器
func (g *Group) All(u string, h HandlerFunc)  { g.NewHandler(METHOD_ALL, u, h) }
func (g *Group) Get(u string, h HandlerFunc)  { g.NewHandler(METHOD_GET, u, h) }
func (g *Group) Post(u string, h HandlerFunc) { g.NewHandler(METHOD_POST, u, h) }

func (g *Group) UseGroup(uri string, group *Group) {
	g.Trie.InsertByGroup(uri, group)
}

func (g *Group) Use(mw ...MiddlewareFunc) {
	for _, f := range mw {
		g.Trie.Middleware = append(g.Trie.Middleware, &f)
	}
}

func NewGroup() *Group {
	return &Group{
		Trie: new(TrieNode),
	}
}