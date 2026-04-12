package v1

import (
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"
	"practice-7/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, t usecase.UserInterface, l logger.Interface) {
	// Global rate limiter: 100 requests per minute per user/IP
	rateLimiter := utils.NewRateLimiter(100, time.Minute)

	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, t, l, rateLimiter)
	}
}
