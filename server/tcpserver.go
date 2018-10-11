package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type TcpServer struct {
	config *Config
}

// sync
func (srv *TcpServer) Run() {
	fmt.Println("Starting tcp server on port", srv.config.Port)
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", srv.config.Port))
	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	defer conn.Close()
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
