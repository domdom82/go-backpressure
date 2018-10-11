package server

import (
	"fmt"
	"net/http"

	"bufio"
	"time"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	config   *Config
	upgrader websocket.Upgrader
}

// sync
func (srv *WsServer) Run() {
	fmt.Println("Starting websocket server on port", srv.config.Port)

	srv.upgrader = websocket.Upgrader{}

	http.HandleFunc("/", srv.handleDefault)

	http.ListenAndServe(fmt.Sprintf(":%s", srv.config.Port), nil)
}

func (srv *WsServer) handleDefault(w http.ResponseWriter, r *http.Request) {
	conn, err := srv.upgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}
	_, reader, err := conn.NextReader()
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	buf := make([]byte, srv.config.Bufsize)

	for {
		nbytes, err := bufio.NewReaderSize(reader, srv.config.Bufsize).Read(buf)
		fmt.Printf("read %d bytes", nbytes)
		if err != nil {
			fmt.Printf(" (%s)", err)
		}
		fmt.Println()
		time.Sleep(srv.config.Delay * time.Millisecond)
	}

}
