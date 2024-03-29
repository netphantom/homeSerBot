package web

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"net/http"
)

func (dash *dashboard) routes() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery(), gin.Logger(), ginsession.New())
	r.Static("/static", "./ui/static")
	r.LoadHTMLGlob("ui/html/*")
	r.HTMLRender = templateRenderer()

	r.GET("/login", dash.showLogin)
	r.POST("/login", dash.login)

	private := r.Group("/dashboard")
	private.Use(dash.AuthRequired, dash.SetNotificationNumber)
	{
		private.GET("", dash.home)
		private.GET("/profile", dash.profile)

		private.GET("/notifications", dash.showNotifications)
		private.GET("/adminMode", dash.adminMode)

		private.GET("/changePassword", dash.showChangePassword)
		private.POST("/changePassword", dash.changePassword)

		private.GET("/showProcessesList", dash.showProcessesList)
		private.GET("/deleteProcess", dash.deleteProcess)
		private.POST("/editProcess", dash.editProcess)

		private.POST("/processAdd", dash.processAdd)
		private.GET("/processDetail", dash.showProcessDetails)
		private.GET("/showNewProcess", dash.showNewProcess)

		private.GET("/logout", dash.logout)
	}

	r.NoRoute(dash.notFound)
	return r
}
