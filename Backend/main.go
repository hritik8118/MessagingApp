package main

import (
	"MessasingApp/Backend/auth"
	"MessasingApp/Backend/handlers"

	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("[Backend] Starting server on :8080...")
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/register", func(c *gin.Context) {
		fmt.Println("[Backend] /register endpoint hit")
		handlers.RegisterHandler(c)
	})
	r.POST("/login", func(c *gin.Context) {
		fmt.Println("[Backend] /login endpoint hit")
		handlers.LoginHandler(c)
	})

	authGroup := r.Group("/chat")
	authGroup.Use(auth.JWTAuthMiddleware())
	authGroup.GET("/ws", handlers.WebSocketHandler)

	r.Run(":8080")
}
