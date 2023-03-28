package grun

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct{
	*RouterGroup
	router *Router
	groups []*RouterGroup // 存储所有的组
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 获得上下文
	cxt := newContext(w, req)
	// 加入中间件
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) { // 如果是当前组的中间件
			cxt.handlers = append(cxt.handlers, group.middlewares...)
		}	
	}
	e.router.handle(cxt)
}

// 初始化引擎
func New() *Engine {
	e := &Engine{
		router: newRouter(),
	}
	e.RouterGroup = &RouterGroup{engine: e}
	e.groups = []*RouterGroup{e.RouterGroup}
	return e
}

// 进行分组，所有的组共用同一个engine
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	e := g.engine // 获得当前组的engine
	newGroup := &RouterGroup{ // 初始化
		prefix: g.prefix + prefix,
		parent: g,
		engine: e,
	}
	e.groups = append(e.groups, newGroup) // 加入该组
	return newGroup
}

// 使用中间件,可以使用多个
func (g *RouterGroup) Use(handler ...HandlerFunc) {
	g.middlewares = append(g.middlewares, handler...)
}

// 传入addr 启动web服务
func (e *Engine) Run(addr string) error {
	log.Println("web serve is starting....")
	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", addr), e)
} 

// 添加路由 method 方法 pattern 路径 handler 处理函数
func (g *RouterGroup) addRoute(method , path string, handler HandlerFunc)  {
	// 访问时需加上prefix
	g.engine.router.addRoute(method, g.prefix + path, handler)
}

// 具体化方法
func (g *RouterGroup) POST(path string, handler HandlerFunc) {
	g.addRoute("POST", path, handler)
}

func (g *RouterGroup) GET(path string, handler HandlerFunc) {
	g.addRoute("GET", path, handler)
}

func (g *RouterGroup) PUT(path string, handler HandlerFunc) {
	g.addRoute("PUT", path, handler)
}

func (g *RouterGroup) DELETE(path string, handler HandlerFunc) {
	g.addRoute("DELETE", path, handler)
}