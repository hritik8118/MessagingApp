package main

import (
	"MessasingApp/Backend/auth"
	"MessasingApp/Backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)

	authGroup := r.Group("/chat")
	authGroup.Use(auth.JWTAuthMiddleware())
	authGroup.GET("/ws", handlers.WebSocketHandler)

	r.Run(":8080")
}
