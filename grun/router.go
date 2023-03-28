package grun

import (
	"net/http"
	"strings"
)

type Router struct {
	nodes  map[string]*node       // 动态路由, 每一种方法使用一个trie树
	routes map[string]HandlerFunc // 路由映射
}

// 初始化路由
func newRouter() *Router {
	return &Router{
		nodes:  make(map[string]*node),
		routes: make(map[string]HandlerFunc),
	}
}

// 解析路径
func parsePath(path string) []string {
	ss := strings.Split(path, "/")

	parts := make([]string, 0)
	for _, item := range ss {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // 如果含有*，则无需继续考虑
				break
			}
		}
	}
	return parts
}

// 添加路由映射关系
func (r *Router) addRoute(method, path string, handler HandlerFunc) {
	parts := parsePath(path)

	key := method + "-" + path
	_, ok := r.nodes[method]
	if !ok {
		r.nodes[method] = &node{}
	}
	r.nodes[method].insert(path, parts, 0)
	r.routes[key] = handler
}

// 获得该路径的父节点
func (r *Router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePath(path)	
	// 返回真实的路由的映射
	params := make(map[string]string)

	// 获得该方法的trie树
	node, ok := r.nodes[method]
	if !ok {
		return nil, nil
	}

	n := node.search(searchParts, 0)
	if n != nil {
		parts := parsePath(n.path)
		// hello/:lan/cc
		// hello/java/cc params[lan] = java
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}


// 处理请求
func (r *Router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node == nil {
		c.Data(-1, []byte("处理请求出现错误"))
		return
	}
	c.Params = params
	// 获取键
	key := c.Method + "-" + node.path
	// 处理请求
	if handler, ok := r.routes[key]; ok {
		// 加入待处理的函数中
		c.handlers = append(c.handlers, handler)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
