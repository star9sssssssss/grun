package grun

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// 上下文
type Context struct {
	W http.ResponseWriter 		// 响应数据
	Req *http.Request 			// 请求
	Path string  				// 请求路径
	Method string 				// 请求方法
	StatusCode int 				// 响应状态码
	Params map[string]string	// 路由的真实映射
	handlers []HandlerFunc  	// 所有的处理函数，包括请求和中间件
	index int					// Handlers调用的下标
}

// 初始化上下文
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		W: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
		index: -1,
	}
}

// 处理中间件和请求函数，按顺序
// 设置Next() 是为了能进行回调中间件，使用Next()后会处理后续的函数，如果后续函数处理完，再处理该函数的后续部分
func (c *Context) Next() {
	c.index ++
	len := len(c.handlers)
	for ; c.index < len; c.index ++ {
		c.handlers[c.index](c)
	}
}

// 设置响应头
func (c *Context) SetHeader(key, value string) {
	c.W.Header().Set(key, value)
}

// 设置响应状态
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.W.WriteHeader(code)
}

// 封装一些特点的处理方法
func (c *Context) JSON(code int, data interface{}) {
	c.SetHeader("Content-Type", "application/json;charset=UTF-8")
	c.Status(code)
	// 返回一个在c.W写入json数据的encoder
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(data); err != nil {
		http.Error(c.W, err.Error(), 500)
	}
}

// 解析表单(常用于POST请求)
func (c *Context) PostForm(key string) string {
	return c.Req.PostFormValue(key)
}

// 查询路径参数(常用于GET请求)
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 写入返回的正常数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.W.Write(data)
}
// 返回数据以字符串格式
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain;charset=UTF-8")
	c.Status(code)
	c.W.Write([]byte(fmt.Sprintf(format, values...)))
}

// 返回html
func (c *Context) HTML(code int, html string) {
	// 设置格式
	c.SetHeader("Content-Type", "text/html;charset=UTF-8")
	c.Status(code)
	c.W.Write([]byte(html))
}

// 发生错误的返回
func (c *Context) Fail(code int, data string) {
	c.SetHeader("Content-Type", "text/plain;charset=UTF-8")
	c.Status(code)
	c.W.Write([]byte(data))
} 

// 获得:,*匹配的真实路径
// hello/:lan/cc
// hello/java/cc
// hello/go/cc
func (c *Context) GetUrlPart(key string) string {
	// map[lan] = java map[lan] = go
	return c.Params[key]
}