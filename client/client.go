package client

import (
	"fmt"
	"io/ioutil"
	"net"

	"math/rand"
	"time"

	"bytes"

	"bufio"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerAddr        string `yaml:"server_addr"`
	ServerPort        string `yaml:"server_port"`
	PayloadSize       int    `yaml:"payload_size"`
	RequestsPerSecond int    `yaml:"requests_per_second"`
	RequestsTotal     int    `yaml:"requests_total"`
}

type Client struct {
	config *Config
}

func (cfg *Config) NewClient() *Client {
	server := &Client{config: cfg}

	return server
}

func NewClientConfigFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config *Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Client) makePayload(size int) []byte {
	payload := new(bytes.Buffer)

	for i := 0; i < size; i++ {
		payload.WriteByte(byte(rand.Intn(255)))
	}
	return payload.Bytes()
}

// sync
func (c *Client) Run() {
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
		payload := c.makePayload(c.config.PayloadSize)
		nbytes, err := w.Write(payload)
		fmt.Printf("req: %d (wrote %d bytes)", r, nbytes)
		if err != nil {
			fmt.Printf(" (%s)", err)
		}
		fmt.Println()
		fmt.Fprintf(conn, "\n")
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
