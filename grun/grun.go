package grun

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct{
	router map[string]HandlerFunc
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	log.Println(req.Method, req.URL.Path)
	if handler, ok := e.router[key]; ok {
		handler(w, req)
	} else {
		log.Println("appear error")
	}
}

// 初始化引擎
func New() *Engine {
	return &Engine{
		router: map[string]HandlerFunc{},
	}
}

// 传入addr 启动web服务
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(fmt.Sprintf("localhost:%s", addr), e)
} 

// 添加路由 method 方法 pattern 路径 handler 处理函数
func (e *Engine) addRoute(method , path string, handler HandlerFunc)  {
	key := method + "-" + path
	e.router[key] = handler
}

// 具体化方法
func (e *Engine) POST(path string, handler HandlerFunc) {
	e.addRoute("POST", path, handler)
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) PUT(path string, handler HandlerFunc) {
	e.addRoute("PUT", path, handler)
}

func (e *Engine) DELETE(path string, handler HandlerFunc) {
	e.addRoute("DELETE", path, handler)
}