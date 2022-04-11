package pepper

type Group struct {
	Router *Router
}

func (p *Pepper) CreateGroup() *Group {
	return NewGroup()
}

// 创建处理函数
func (g *Group) NewHandler(method string, uri string, handler HandlerFunc) {
	if g.Router == nil {
		g.Router = new(Router)
	}

	if uri == "/" {
		g.Router.NewHandler(method, handler)
		return
	}

	g.Router.PushRouter(method, uri, handler, true)
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
	if g.Router == nil {
		g.Router = new(Router)
	}
	g.Router.PushRouterByGroup(uri, group)
}

func NewGroup() *Group {
	return &Group{}
}