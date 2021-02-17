package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (dash *dashboard) routes() http.Handler {
	r := gin.New()
	r.Use(loggingMiddleware)
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Static("/static", "./ui/static")

	r.LoadHTMLGlob("ui/html/*")
	r.LoadHTMLGlob("ui/templates/*")

	r.GET("/login", dash.showLogin)
	r.POST("/login", dash.login)
	//r.GET("/protected-content", serveProtectedContent)

	return r
}
