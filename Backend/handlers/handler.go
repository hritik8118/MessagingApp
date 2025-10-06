package handlers

import (
	"fmt"
	"net/http"
	"time"

	"MessasingApp/backend/auth"
	"MessasingApp/backend/db"
	"MessasingApp/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// allow from any origin for local testing, refine for production
		return true
	},
}

func RegisterHandler(c *gin.Context) {
	fmt.Println("[Backend] /register endpoint hit")
	var req models.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}
	ok := db.CreateUser(req.Username, req.Password)
	if !ok {
		c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
		return
	}
	c.Status(http.StatusCreated)
}

func LoginHandler(c *gin.Context) {
	fmt.Println("[Backend] /login endpoint hit")
	var req models.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if !db.ValidateUser(req.Username, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := auth.GenerateJWT(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}
	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}

// WebSocketHandler upgrades first, then validates token and reads/writes messages
func WebSocketHandler(c *gin.Context) {
	fmt.Println("[Backend] /chat/ws endpoint hit (ws handler) - attempting upgrade")
	token := c.Query("token")
	// Upgrade the connection first so we can send WS close messages with reasons
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("[Backend] upgrade failed:", err)
		return
	}
	defer conn.Close()
	fmt.Println("[Backend] websocket upgrade OK; validating token")

	username, err := auth.ValidateJWT(token)
	if err != nil {
		fmt.Println("[Backend] invalid token on ws:", err)
		// Send a close message with reason
		closeMsg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "invalid token")
		_ = conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
		return
	}
	fmt.Println("[Backend] WebSocket authenticated for user:", username)

	// Send a welcome message
	_ = conn.WriteJSON(map[string]string{"system": "welcome", "user": username})

	// simple echo-like loop that stores messages in-memory
	for {
		var incoming models.Message
		if err := conn.ReadJSON(&incoming); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("[Backend] client closed connection:", err)
			} else {
				fmt.Println("[Backend] read error:", err)
			}
			break
		}
		// populate message fields server-side
		msg := db.NewMessage(username, incoming.Receiver, incoming.Content)
		db.AddMessage(msg)
		// echo back to sender
		if err := conn.WriteJSON(msg); err != nil {
			fmt.Println("[Backend] write error:", err)
			break
		}
	}
	fmt.Println("[Backend] ws handler exit for user:", username)
}
