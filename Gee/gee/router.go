package gee

import (
	"net/http"
	"strings"
)

type router struct {
	//为什么这里会有多个前缀树？好的因为需要为不同类型的请求分别设置路由
	//如GET POST 对应的前缀树是不一样的
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute 基于trie模块提供的insert函数实现路由添加
// 输入为 url 大概不需要返回值，因为对应不同类型请求有不同的前缀树所以要传入method
// 最后需要为了建立映射关系自然需要传入handler
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	//首先对pattern进行解析
	parts := parsePattern(pattern)
	//就操作上其实还是通过map储存请求与handler的匹配关系
	//不过这个key值是从前缀树中获取的
	key := method + "-" + pattern
	//先看看对应的method是否存在前缀树
	if _, ok := r.roots[method]; !ok {
		//构造一个空的头节点
		r.roots[method] = &node{}
	}
	//尝试往前缀树插入新节点
	r.roots[method].insert(pattern, parts, 0)
	//储存key与handler关系
	r.handlers[key] = handler
}

// getRoute 基于前缀树获取
// 获取匹配节点，以及动态路由获得的参数
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	//正常来讲如果获取的n不为空就证明找到了，但是我们还得尝试获取参数
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			} else if part[0] == '*' && len(part) == 1 {
				// /static/css/geektutu.css匹配到/static/*filepath
				//解析结果为{filepath: "css/geektutu.css"}。
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// handle 执行handler 这个函数提供为ServeHTTP 使用
func (r *router) handle(c *Context) {
	n, param := r.getRoute(c.Method, c.Path)
	//获取到匹配的路径后执行handler
	if n != nil {
		//储存Param到Context里面
		c.Params = param
		//获取handler并执行
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
