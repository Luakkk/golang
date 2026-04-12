package v1

import (
	"net/http"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userRoutes struct {
	t usecase.UserInterface
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface, l logger.Interface, rateLimiter *utils.RateLimiterStore) {
	r := &userRoutes{t, l}

	h := handler.Group("/users")
	{
		// Public routes
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)

		// Protected routes — require valid JWT
		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		protected.Use(rateLimiter.Middleware())
		{
			protected.GET("/me", r.GetMe)                                           // Task 1
			protected.GET("/protected/hello", r.ProtectedFunc)

			// Admin-only routes — Task 2
			admin := protected.Group("/")
			admin.Use(utils.RoleMiddleware("admin"))
			{
				admin.PATCH("/promote/:id", r.PromoteUser)
			}
		}
	}
}

// POST /users/ — register a new user
func (r *userRoutes) RegisterUser(c *gin.Context) {
	var createUserDTO entity.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(createUserDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	role := createUserDTO.Role
	if role == "" {
		role = "user"
	}

	user := entity.User{
		Username: createUserDTO.Username,
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Role:     role,
	}

	createdUser, sessionID, err := r.t.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully. Please check your email for verification code.",
		"session_id": sessionID,
		"user":       createdUser,
	})
}

// POST /users/login — login and get JWT token
func (r *userRoutes) LoginUser(c *gin.Context) {
	var input entity.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.t.LoginUser(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GET /users/me — Task 1: return info about current user from JWT
// Requires: JWTAuthMiddleware
// No extra data in request — reads userID from JWT context only
func (r *userRoutes) GetMe(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	user, err := r.t.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// PATCH /users/promote/:id — Task 2: promote a user to admin
// Requires: JWTAuthMiddleware + RoleMiddleware("admin")
func (r *userRoutes) PromoteUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := r.t.PromoteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User promoted to admin successfully",
		"user":    user,
	})
}

// GET /users/protected/hello — example protected endpoint
func (r *userRoutes) ProtectedFunc(c *gin.Context) {
	c.JSON(200, gin.H{"message": "OK"})
}
