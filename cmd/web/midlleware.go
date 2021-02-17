package main

import (
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
	"homeSerBot/pkg/mysqlmodels"
	"net/http"
)

func loggingMiddleware(c *gin.Context) {

}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (dash *dashboard) AuthRequired(c *gin.Context) {
	exists := dash.session.Exists(c.Request, "authenticatedUserID")
	if !exists {
		c.AbortWithError(http.StatusUnauthorized, mysqlmodels.ErrInvalidCredentials)
		return
	}
	uid := dash.session.GetInt(c.Request, "authenticatedUserID")
	user := dash.users.VerifyId(uint(uid))
	if user == nil {
		c.AbortWithError(http.StatusUnauthorized, mysqlmodels.ErrInvalidCredentials)
		return
	}

	c.Next()
}
