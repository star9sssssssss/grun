# go web 框架
## 1.0 version
#### 1.http库的基本使用
```GO
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
  // 处理"/hello"请求，第二个参数为处理函数
	http.HandleFunc("/hello", helloHandler)
  // 监听本地的9999端口号，第二个参数填nil时，使用默认的路由管理器
	log.Fatal(http.ListenAndServe(":9999", nil))
}

//参数1: 用于响应请求 参数2: 代表请求
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
```
```GO
func ListenAndServe(addr string, handler Handler) {}
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
所以我们需要实现ServeHttp() 方法实现，自定义路由管理器
#### 2.模仿gin的调用模式
```GO
// 初始化
r := gin.New()
// 设置路由
r.Get(path, handler)
// 运行
r.Run()
```
```GO
// 抽象处理请求的函数
type HandlerFunc func(http.ResponseWriter, *http.Request)
// 所以建立一个模块
type Engine struct{
  // 建立请求与处理函数的映射
	router map[string]HandlerFunc 
}

// 添加路由 method 方法 pattern 路径 handler 处理函数
func (e *Engine) addRoute(method , path string, handler HandlerFunc)  {
  // 键的格式为Get-/hello
	key := method + "-" + path
	e.router[key] = handler
}

// 传入addr 启动web服务
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(fmt.Sprintf("localhost:%s", addr), e)
} 
// 根据传来的请求，选择映射的函数进行处理
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	log.Println(req.Method, req.URL.Path)
	if handler, ok := e.router[key]; ok {
		handler(w, req)
	} else {
		log.Println("appear error")
	}
}

```
## 2.0 version
### 1.添加动态路由
为了让类似于/hello/:Go/cc, /hello/*的请求能够响应, 使用动态路由，以满足该模式该路径匹配
```Go
// 动态路由的实现，使用trie树
type node struct {
	path string  // 总请求路径  /hello/hh/c
	part string  // 该节点的配对路径 /hh
	parent *node // 父节点 /hello
	children []*node // 子节点 /c
	isSpecial bool // 是否是特别匹配 ':' or '*' 为true
}
// 定义node模块，part代表每段路径，path代表总路径，且每段路径都有父节点和子节点，并且包括特殊字符的匹配
```
* 方法`func (n *node) matchChild(part string) *node` 根据part找到第一个匹配的节点,用于插入
* 方法`func (n *node) insert(path string, parts []string, height int)` 插入节点
* 方法`func (n *node) matchChildren(part string) []*node` 找到所有符合part的路径，用于查找
* 	方法`func (n *node) search(parts []string, height int) *node` 找到路径的父节点
```Go
type Router struct {
	nodes  map[string]*node       // 动态路由, 每一种方法使用一个trie树
	routes map[string]HandlerFunc // 路由映射
}
```
### 2.分组
#### 1.起因
在现实状况中，存在一些请求具有相同性，所以往往具有相同的请求前缀，所以我们可以根据功能进行分组，然后分别进行处理
```Go
func main() {
	r := grun.New()
	rg := r.Group("/gg") // 分组
	rg.GET("/hello/:/cc", func(c *grun.Context) { // 动态匹配
		c.HTML(200, "<h1>这是标题1</h1>")
	})
	r.Run("8080")
}
```
#### 2.实现
为了使分组后依然能实现路由的管理，`RouterGroup` 应该含有 `*Engine` 该属性， 使所有分组都属于同一个`Engine`,
整个框架的所有资源都是由Engine统一协调的，那么就可以通过Engine间接地访问各种接口了
```GO
type Engine struct{
	*RouterGroup
	router *Router
	groups []*RouterGroup // 存储所有的组
}
}
我们还可以进一步地抽象，将`Engine`作为最顶层的分组，也就是说`Engine`拥有`RouterGroup`所有的能力
那我们就可以将和路由有关的函数，都交给`RouterGroup`实现了
type RouterGroup struct {
	prefix string 		      // 分组路由的前缀
	parent *RouterGroup 	  // 父分组
	engine *Engine 		      // 所有分组共享一个引擎
	middlewares []HandlerFunc // 中间件
}

// 添加路由 method 方法 pattern 路径 handler 处理函数
func (g *RouterGroup) addRoute(method , path string, handler HandlerFunc)  {
	// 访问时需加上prefix
	g.engine.router.addRoute(method, g.prefix + path, handler)
}
```
## 3.0 version
### 添加中间件(运行用户自定义逻辑，完美嵌入到框架中)
```GO
// 在 RouterGroup 中添加 middlewares []HandlerFunc 属性，是为了使用户定制的中间件可以在任何组内使用，进行区分

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


// 处理中间件和请求函数，按顺序
// 设置Next() 是为了能进行回调中间件，使用Next()后会处理后续的函数，如果后续函数处理完，再处理该函数的后续部分
func (c *Context) Next() {
	c.index ++
	len := len(c.handlers)
	for ; c.index < len; c.index ++ {
		c.handlers[c.index](c)
	}
}
```

```Go
// 定义中间件
func Logger() HandlerFunc {
	return func(c *Context) {
		log.SetPrefix("[before]")
		log.Println("middleware logger is starting...")
		c.Next()
		log.SetPrefix("[after]")
		log.Println("middleware logger is stoping...")
	}
}
// 使用
r.Use(grun.Logger())

```









