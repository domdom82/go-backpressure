package server

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TcpServer struct {
	config *Config
}

func (srv *TcpServer) HandleConn(conn net.Conn) {
	fmt.Printf("\nNew connection: %v\n", conn.RemoteAddr())
	defer conn.Close()
	buf := make([]byte, srv.config.Bufsize)
	for {
		nbytes, err := conn.Read(buf)
		fmt.Printf("read %d bytes", nbytes)
		if err != nil {
			fmt.Printf(" (%s)", err)
			if err == io.EOF {
				break
			}
		}
		fmt.Println()
		time.Sleep(srv.config.Delay * time.Millisecond)
	}
}

// sync
func (srv *TcpServer) Run() {
	fmt.Println("Starting tcp server on port", srv.config.Port)
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%s", srv.config.Port))
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go srv.HandleConn(conn)
	}
}
