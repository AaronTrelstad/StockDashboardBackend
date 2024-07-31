package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error while upgrading connection: ", err)
		return
	}

	defer conn.Close()

	var prev_value = 50.00;

	for {
		stockPrice := math.Max(rand.NormFloat64() * 5 + prev_value, 0);
		prev_value = stockPrice;
		message := fmt.Sprintf(`{"timestamp":%d, "price":%.2f}`, time.Now().Unix() * 1000, stockPrice)

		err = conn.WriteMessage(websocket.TextMessage, []byte(message))

		if err != nil {
			fmt.Println("Error while writing message: ", err)
			return 
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
    fmt.Println("WebSocket server started at ws://localhost:8082/ws")
    if err := http.ListenAndServe(":8082", nil); err != nil {
        fmt.Println("Error while starting server:", err)
    }
}
