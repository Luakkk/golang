package v1

import "github.com/gin-gonic/gin"

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{"error": msg})
}
