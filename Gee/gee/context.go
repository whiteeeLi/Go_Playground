package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 这里的 H 会被作为输入扔到生成json的函数中作为参数
// 在main.go 中可以直接使用gee.H{xxx}直接生成输入
// 如果没有 H 就得使用 map[string]interface{} 做类型声明
type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
}

// 需要一个创建Context的函数，值得注意的是这里是私有的函数，仅在context包中使用
// 对应 gee.go中的 c := newContext(w, req)
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	//返回一个Context对象的地址
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Query 提供对url中参数的提取功能，返回用户查询参数对应的值
func (c *Context) Query(key string) string {
	//下面这玩意想要写出来需要查http包的相关函数，
	//相当于是封装了一下，不用用户自己打这个代码
	return c.Req.URL.Query().Get(key)
}

// PostForm 提供查询Post提交表单的内容的功能
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// SetHeader 用于设置请求头
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// Status 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// String 往响应中写入string
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Data 用于直接向响应体中写入二进制数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// JSON 向响应体中写入json数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// HTML 向响应体写入html
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
