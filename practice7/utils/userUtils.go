package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserIDFromContext extracts and parses the userID UUID from gin context.
// The value is set by JWTAuthMiddleware as a string.
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDStr, _ := c.Get("userID")
	return uuid.Parse(userIDStr.(string))
}
