// !! om namah bhagwate vasudevay !!
package main

import (
	"fmt"
	"log"
	"net/http"

	"MessasingApp/Backend/handlers"
)

func main() {
	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/ws", handlers.WsHandler)

	fmt.Println("🚀 Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
