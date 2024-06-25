package gee

import "net/http"

// HandlerFunc defines the request handler used by gee
// 这里是最先发生变化的地方，类型改了
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
// 因为将路由相关的东西封装到了router.go中，所以这里直接引用router
type Engine struct {
	router *router
}

// New is the constructor of gee.Engine
func New() *Engine {
	//本来直接创map，但封到router中
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
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
