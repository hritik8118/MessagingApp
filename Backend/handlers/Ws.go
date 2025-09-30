package handlers

import (
	"MessasingApp/Backend/auth"
	"MessasingApp/Backend/db"
	"MessasingApp/Backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(c *gin.Context) {
	token := c.Query("token")
	fmt.Println("WebSocket token received:", token)
	username, err := auth.ValidateJWT(token)
	if err != nil {
		fmt.Println("JWT validation failed in WebSocketHandler:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer conn.Close()

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		msg.ID = time.Now().Format("20060102150405")
		msg.Sender = username
		msg.Timestamp = time.Now().Unix()
		db.AddMessage(msg)
		conn.WriteJSON(msg)
	}
}
