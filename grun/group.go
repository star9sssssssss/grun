package grun

type RouterGroup struct {
	prefix string 		      // 分组路由的前缀
	parent *RouterGroup 	  // 父分组
	engine *Engine 		      // 所有分组共享一个引擎
	middlewares []HandlerFunc // 中间件
}



