package middleware

import (
	xss "github.com/araujo88/gin-gonic-xss-middleware"
	"github.com/gin-gonic/gin"
)

// XSS func
func XSS() gin.HandlerFunc {
	var xssMiddleware xss.XssMw
	return xssMiddleware.RemoveXss()
}
