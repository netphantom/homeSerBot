package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (dash *dashboard) routes() http.Handler {
	r := gin.New()

	r.Use(gin.Recovery(), gin.Logger())
	r.Static("/static", "./ui/static")
	r.LoadHTMLGlob("ui/pages/*")

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	r.GET("/login", dash.showLogin)
	r.POST("/login", dash.login)

	private := r.Group("/dashboard")
	private.Use(dash.AuthRequired)
	{
		private.GET("/home")
	}

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "notfound.html", nil)
	})
	return r
}
