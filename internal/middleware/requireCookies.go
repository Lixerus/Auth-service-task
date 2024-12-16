package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AccessCookie  string = "AccessCookieString"
	RefreshCookie string = "RefreshCookieString"
)

func RequireCookies(c *gin.Context) {
	cookie, err := c.Request.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "No auth cookie"})
		return
	}
	if cookie.Value == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "Bad auth cookie"})
		return
	}
	c.Set(AccessCookie, cookie.Value)

	cookie, err = c.Request.Cookie("Authrefresh")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "No refresh cookie"})
		return
	}
	if cookie.Value == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "Bad auth cookie"})
		return
	}
	c.Set(RefreshCookie, cookie.Value)
	c.Next()
}
