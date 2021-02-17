package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (dash *dashboard) showLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (dash *dashboard) login(c *gin.Context) {
	email := c.PostForm("userEmail")
	password := c.PostForm("userPassword")

	log.Printf("%s , %s", email, password)
}

/*
// Handler for the login request
func login(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	if response := auth.Login(username, password); response.Token != "" {
		// If authentication succeeds set the cookies and
		// respond with an HTTP success
		// status and include the token in the response
		c.SetCookie("username", username, 3600, "", "", false, true)
		c.SetCookie("token", response.Token, 3600, "", "", false, true)

		c.JSON(http.StatusOK, response)
	} else {
		// Respond with an HTTP error if authentication fails
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

// Handler for the logout request
func logout(c *gin.Context) {
	// Obtain the username and token from the cookies
	username, err1 := c.Cookie("username")
	token, err2 := c.Cookie("token")

	if err1 == nil && err2 == nil && auth.Logout(username, token) {
		// Clear the cookies and
		// respond with an HTTP success status
		c.SetCookie("username", "", -1, "", "", false, true)
		c.SetCookie("token", "", -1, "", "", false, true)

		c.JSON(http.StatusOK, nil)
	} else {
		// Respond with an HTTP error
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

// Handler to serve the protected content
func serveProtectedContent(c *gin.Context) {
	// Obtain the username and token from the cookies
	username, err1 := c.Cookie("username")
	token, err2 := c.Cookie("token")

	if err1 == nil && err2 == nil && auth.Authenticate(username, token) {
		// Respond with an HTTP success status and include the
		// content in the response

		c.JSON(http.StatusOK, gin.H{"content": "This should be visible to authenticated users only."})
	} else {
		// Respond with an HTTP error
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

*/
