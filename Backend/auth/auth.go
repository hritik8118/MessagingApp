package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("replace-this-secret-with-env-or-config")

// GenerateJWT creates a token for the given username (short expiry for testing)
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates and returns the subject (username)
func ValidateJWT(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("empty token")
	}
	// Allow tokens with "Bearer " prefix if provided
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	token, err := jwt.Parse(tokenString, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, _ := claims["sub"].(string)
		return sub, nil
	}
	return "", errors.New("invalid token")
}

// JWTAuthMiddleware is useful for REST endpoints; it also accepts token via ?token= for convenience
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Authorization header first
		authHeader := c.GetHeader("Authorization")
		token := authHeader
		if token == "" {
			// fallback to query param
			token = c.Query("token")
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		user, err := ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		// store user in context for handlers
		c.Set("username", user)
		c.Next()
	}
}
