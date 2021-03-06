package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/domdom82/go-backpressure/client"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	config         *Config
	upgrader       websocket.Upgrader
	clientTemplate *template.Template
}

// sync
func (srv *WsServer) Run() {
	fmt.Println("Starting web socket server on port", srv.config.Port)

	srv.upgrader = websocket.Upgrader{
		ReadBufferSize:  srv.config.Bufsize,
		WriteBufferSize: srv.config.Bufsize,
	}
	srv.setupClientTemplate()

	http.HandleFunc("/ws", srv.handleDefault)
	http.HandleFunc("/client", srv.serveClient)

	http.ListenAndServe(fmt.Sprintf(":%s", srv.config.Port), nil)
}

func (srv *WsServer) handleDefault(w http.ResponseWriter, r *http.Request) {
	conn, err := srv.upgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}

	go srv.HandleConn(conn)

}

func (srv *WsServer) HandleConn(conn *websocket.Conn) {
	fmt.Printf("\nNew connection: %v\n", conn.RemoteAddr())

	defer conn.Close()

	// We expect the client to tell us how much he'll send us
	cp := client.ConnectionParams{}
	err := conn.ReadJSON(&cp)

	if err != nil {
		fmt.Printf("error: Expected connection params but got: %v\n", err)
		return
	}

	bytesTotal := cp.RequestsTotal * cp.PayloadSize
	bytesRead := 0

	for {
		msgType, buf, err := conn.ReadMessage()
		if msgType == websocket.BinaryMessage {
			bytesRead += len(buf)
			completedPercent := float64(bytesRead) / float64(bytesTotal) * 100.0
			fmt.Printf("%s : read %d bytes (%d / %d bytes total %.2f%%)\n", conn.RemoteAddr(), len(buf), bytesRead, bytesTotal, completedPercent)
		} else if msgType == websocket.TextMessage {
			fmt.Printf("read text message: %v \n", string(buf))
		} else if msgType == websocket.CloseMessage {
			fmt.Println("Received close message. Closing connection.")
			break
		}

		if err != nil {
			fmt.Printf("%s : %s\n", conn.RemoteAddr(), err)
			break
		}
		time.Sleep(srv.config.Delay * time.Millisecond)
	}
}

func (srv *WsServer) setupClientTemplate() {
	html, _ := ioutil.ReadFile("server/client.html")
	srv.clientTemplate = template.Must(template.New("").Parse(string(html)))
}

func (srv *WsServer) serveClient(w http.ResponseWriter, r *http.Request) {
	srv.clientTemplate.Execute(w, "")
}
