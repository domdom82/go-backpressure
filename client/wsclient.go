package client

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	config *Config
}

// sync
func (c *WsClient) Run() {
	fmt.Println("Starting client for websocket server", c.config.ServerAddr, "on port", c.config.ServerPort)

	rand.Seed(time.Now().Unix())

	timePerRequest := time.Duration(1.0 / float64(c.config.RequestsPerSecond) * float64(time.Second))
	expectedTotalTime := time.Duration(float64(c.config.RequestsTotal) / float64(c.config.RequestsPerSecond) * float64(time.Second))

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws", c.config.ServerAddr, c.config.ServerPort), nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Tell the server how much we are going to send
	conn.WriteJSON(ConnectionParams{PayloadSize: c.config.PayloadSize, RequestsTotal: c.config.RequestsTotal})

	tStart := time.Now()
	r := 1
	for ; r <= c.config.RequestsTotal; r++ {
		reqStart := time.Now()
		payload := makePayload(c.config.PayloadSize)
		err := conn.WriteMessage(websocket.BinaryMessage, payload)
		fmt.Printf("req: %d (wrote %d bytes)", r, c.config.PayloadSize)
		if err != nil {
			fmt.Printf(" (%s)", err)
		}
		fmt.Println()
		reqStop := time.Now()
		requestTime := reqStop.Sub(reqStart)
		if requestTime < timePerRequest {
			time.Sleep(timePerRequest - requestTime)
		}
	}
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	tEnd := time.Now()
	actualTotalTime := tEnd.Sub(tStart)

	fmt.Printf("\nFired %d requests in %s\n", r-1, actualTotalTime)
	fmt.Printf("Expected: %d requests in %s\n", c.config.RequestsTotal, expectedTotalTime)
}
