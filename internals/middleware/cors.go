package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS sets the appropriate request header values for Cross-Origin Resource Sharing (CORS).
func CORS(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if ctx.Request.Method == http.MethodOptions {
		// If OPTIONS is requested, return the CORS headers and stop request processing
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	// Continue with subsequent middleware/request handling
	ctx.Next()
}
