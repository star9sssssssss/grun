package main

import (
	"go-web/grun"
	"log"
)

func main() {
	r := grun.New()
	r.Use(grun.Logger())
	r.Use(grun.Recovery())
	rg := r.Group("/gg")
	rg.GET("/hello/:/cc", func(c *grun.Context) {
		c.HTML(200, "<h1>这是标题1</h1>")
	})
	r.GET("/hell/*", func(c *grun.Context) {
		log.Println("执行中")
		c.JSON(200, grun.H{
			"name": "韩信",
			"age": 18,
		})
	})
	r.Run("8080")
}