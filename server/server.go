package server

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port    string        `yaml:"port"`
	Bufsize int           `yaml:"bufsize"`
	Delay   time.Duration `yaml:"delay"`
}

type Server struct {
	config *Config
}

func (cfg *Config) NewServer() *Server {
	server := &Server{config: cfg}

	return server
}

func NewServerConfigFromFile(filename string) (*Config, error) {
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

// sync
func (srv *Server) Run() {
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", srv.config.Port))
	conn, _ := listener.Accept()
	buf := make([]byte, srv.config.Bufsize)
	for {

		nbytes, err := bufio.NewReaderSize(conn, srv.config.Bufsize).Read(buf)
		fmt.Printf("read %d bytes", nbytes)
		if err != nil {
			fmt.Printf(" (%s)", err)
		}
		fmt.Println()
		time.Sleep(srv.config.Delay * time.Millisecond)
	}
}
