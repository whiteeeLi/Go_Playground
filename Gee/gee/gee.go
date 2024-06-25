package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc
// 我们需要一个类型来储存handler
// 其实这玩意就是http包里面对HandlerFunc的定义，这里不知道维之门没有直接引用http中的类型
// 毕竟我们还是要把HandlerFunc传给http的函数
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engineer 用于储存请求与handler之间的映射关系
type Engineer struct {
	route map[string]HandlerFunc
}

// New 在使用之前很自然需要先创建一个Engineer对象
func New() *Engineer {
	//结构体的初始化
	var e Engineer = Engineer{
		route: make(map[string]HandlerFunc, 0),
	}
	return &e
}

//ok通过上面的操作我们已经有了一个初始化后的map了，下一笔就需要考虑如何为map添加路由了

// GET 首先当然是经典的Get方法
// Get方法的输入是路由对应的url以及处理这个路由的方法
func (e *Engineer) GET(addr string, handler HandlerFunc) {
	//因为每个Restful风格下每个url对应多种请求类型，所以依据Get以及addr构建key
	//本来这里可以直接往map中添加数据的，但是为了代码复用单开一个函数处理这个行为
	e.addRouter("GET", addr, handler)
}

func (e *Engineer) addRouter(method string, addr string, handler HandlerFunc) {
	key := method + "-" + addr
	if handler == nil {
		fmt.Printf("url %s do not have a correct handler", addr)
	}
	e.route[key] = handler
}

// POST defines the method to add POST request
func (e *Engineer) POST(pattern string, handler HandlerFunc) {
	e.addRouter("POST", pattern, handler)
}

// Run 为了把http包的服务包装起来，需要有一个函数将ListenAndServe封装起来
// 并且把我们的engineer传入上面那个函数
// 同时为了处理异常，最好添加一个err的返回值
func (e *Engineer) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// 为了把engineer 传入到ListenAndServe函数中，需要实现方法ServeHTTP
func (e *Engineer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//内部的方法应当可以处理所有的请求
	//所以这里就应该直接查询当前请求对应的handler是哪一个，
	//然后调用各个handler直接处理请求了
	key := req.Method + "-" + req.URL.Path
	if handler, ok := e.route[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
