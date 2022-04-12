package pepper

type Group struct {
	Node *TreeNode
}

func (p *Pepper) CreateGroup() *Group {
	return NewGroup()
}

// 创建处理函数
func (g *Group) NewHandler(method string, uri string, handler HandlerFunc) {
	if g.Node == nil {
		g.Node = new(TreeNode)
	}

	if uri == "/" {
		g.Node.NewHandler(method, handler)
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

// 创建组
func (g *Group) CreateGroup() *Group {
	return NewGroup()
}

// 使用组
func (g *Group) UseGroup(uri string, group *Group) {
	if g.Node == nil {
		g.Node = new(TreeNode)
	}
	g.Node.PushNodeByGroup(uri, group)
}

func NewGroup() *Group {
	return &Group{}
}