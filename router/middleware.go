package router

import (
	"github.com/gin-gonic/gin"
)

// RequestLogMiddleware :
func RequestLogMiddleware(method *string, path *string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method = &ctx.Request.Method
		path = &ctx.Request.RequestURI
	}
}
