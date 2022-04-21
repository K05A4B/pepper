package pepper

type Group struct {
	ExistRoot bool
	Node      *TreeNode
}

// 创建处理函数
func (g *Group) NewHandler(method string, uri string, handler HandlerFunc) {
	if g.Node == nil {
		g.Node = new(TreeNode)
	}

	if uri == "/" {
		g.Node.NewHandler(method, handler)
		g.ExistRoot = true
		return
	}

	g.Node.PushNode(method, uri, handler, true)
}

func (g *Group) All(uri string, handler HandlerFunc) {
	g.NewHandler(METHOD_ALL, uri, handler)
}

func (g *Group) Get(uri string, handler HandlerFunc) {
	g.NewHandler(METHOD_GET, uri, handler)
}

func (g *Group) Post(uri string, handler HandlerFunc) {
	g.NewHandler(METHOD_POST, uri, handler)
}

// 使用组
func (g *Group) UseGroup(uri string, group *Group) {
	if g.Node == nil {
		g.Node = new(TreeNode)
	}
	g.Node.PushNodeByGroup(uri, group)
}

// 创建组
func NewGroup() *Group {
	return &Group{}
}
