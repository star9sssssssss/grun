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













