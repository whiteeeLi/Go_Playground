package gee

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by gee
// 这里是最先发生变化的地方，类型改了
type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	newPrefix := r.prefix + prefix
	newGroup := &RouterGroup{
		prefix: newPrefix,
		parent: r,
		engine: r.engine,
	}
	return newGroup
}

func (r *RouterGroup) addRouter(method string, pattern string, handler HandlerFunc) {
	newPattern := r.prefix + pattern
	log.Printf("Route %4s - %s", method, newPattern)
	r.engine.router.addRoute(method, newPattern, handler)
}

func (r *RouterGroup) GET(pattern string, handler HandlerFunc) {
	r.addRouter("GET", pattern, handler)
}

func (r *RouterGroup) POST(pattern string, handler HandlerFunc) {
	r.addRouter("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//这里可以看出来context是对 w 与 req 的封装
	c := newContext(w, req)
	engine.router.handle(c)
}
