package grun

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		log.SetPrefix("[before]")
		log.Println("middleware logger is starting...")
		c.Next()
		log.SetPrefix("[after]")
		log.Println("middleware logger is stoping...")
	}
}

// 捕获错误
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func(c *Context) {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}(c)
		c.Next()
	}
}


// 获取panci错误的堆栈信息
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}