package utils

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ─────────────────────────────────────────────
// Password helpers
// ─────────────────────────────────────────────

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// ─────────────────────────────────────────────
// JWT helpers
// ─────────────────────────────────────────────

func GenerateJWT(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ─────────────────────────────────────────────
// JWTAuthMiddleware — validates Bearer token and
// stores user_id / role in the gin context
// ─────────────────────────────────────────────

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["user_id"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

// ─────────────────────────────────────────────
// Task 2: RoleMiddleware — allows only users
// with the specified role to pass through
// ─────────────────────────────────────────────

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
			return
		}
		if role.(string) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access denied: requires role " + requiredRole,
			})
			return
		}
		c.Next()
	}
}

// ─────────────────────────────────────────────
// Task 3: Rate Limiter middleware
// Authenticated: identify by UserID from JWT
// Anonymous:     identify by ClientIP
// Uses sync.Mutex to prevent race conditions
// ─────────────────────────────────────────────

type clientData struct {
	count    int
	resetAt  time.Time
}

type RateLimiterStore struct {
	mu       sync.Mutex
	clients  map[string]*clientData
	max      int
	window   time.Duration
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiterStore {
	return &RateLimiterStore{
		clients: make(map[string]*clientData),
		max:     maxRequests,
		window:  window,
	}
}

func (rl *RateLimiterStore) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine identifier: prefer UserID from JWT, fall back to IP
		var key string
		if userID, exists := c.Get("userID"); exists {
			key = "user:" + userID.(string)
		} else {
			key = "ip:" + c.ClientIP()
		}

		rl.mu.Lock()
		data, ok := rl.clients[key]
		now := time.Now()

		if !ok || now.After(data.resetAt) {
			rl.clients[key] = &clientData{count: 1, resetAt: now.Add(rl.window)}
			rl.mu.Unlock()
			c.Next()
			return
		}

		data.count++
		if data.count > rl.max {
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			return
		}
		rl.mu.Unlock()
		c.Next()
	}
}
