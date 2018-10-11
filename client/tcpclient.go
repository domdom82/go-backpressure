package client

import (
	"fmt"
	"net"

	"math/rand"
	"time"

	"bufio"
)

type TcpClient struct {
	config *Config
}

// sync
func (c *TcpClient) Run() {
	fmt.Println("Starting client for tcp server", c.config.ServerAddr, "on port", c.config.ServerPort)

	rand.Seed(time.Now().Unix())

	timePerRequest := time.Duration(1.0 / float64(c.config.RequestsPerSecond) * float64(time.Second))
	expectedTotalTime := time.Duration(float64(c.config.RequestsTotal) / float64(c.config.RequestsPerSecond) * float64(time.Second))

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.config.ServerAddr, c.config.ServerPort))
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriterSize(conn, c.config.PayloadSize)
	tStart := time.Now()
	r := 1
	for ; r < c.config.RequestsTotal; r++ {
		reqStart := time.Now()
		payload := makePayload(c.config.PayloadSize)
		nbytes, err := w.Write(payload)
		fmt.Printf("req: %d (wrote %d bytes)", r, nbytes)
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
	tEnd := time.Now()
	actualTotalTime := tEnd.Sub(tStart)

	fmt.Printf("\nFired %d requests in %s\n", r, actualTotalTime)
	fmt.Printf("Expected: %d requests in %s\n", c.config.RequestsTotal, expectedTotalTime)
}
