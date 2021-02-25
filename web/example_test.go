package web_test

import (
	"github.com/alfred-zhong/goutil/web"
	"github.com/gin-gonic/gin"
)

func ExampleLoggerMiddleware() {
	e := gin.New()

	e.Use(web.RequestIDMiddleware)
	e.Use(web.LoggerMiddleware(true))

	e.GET("hello", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	if err := e.Run(":9527"); err != nil {
		panic(err)
	}
}
