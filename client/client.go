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

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.config.ServerAddr, c.config.ServerPort))
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriterSize(conn, c.config.PayloadSize)
	for {
		payload := c.makePayload(c.config.PayloadSize)
		w.Write(payload)
		fmt.Fprintf(conn, "\n")
	}
}
