package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)
type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServerWS() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWSOrderBook(ws *websocket.Conn){
	fmt.Println("new incoming connection from client to orderbook feed", ws.RemoteAddr().String())

	for {
		payload := fmt.Sprintf("orderbook feed: %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(1 * time.Second)
	}
}

func (s *Server) handleWS(ws *websocket.Conn){
	fmt.Println("new incoming connection from client", ws.RemoteAddr().String())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF{
				break
		}
		msg := buf[:n]
		fmt.Println("received message from client:", string(msg))
		ws.Write([]byte("thank you for sending a message"))

		s.broadcast(msg)
	}
}
}

func (s *Server) broadcast(msg []byte) {
	for ws := range s.conns{
		go func(ws *websocket.Conn){
			if _, err := ws.Write(msg); err != nil {
				fmt.Println("error broadcasting message to client", err.Error())
				ws.Close()
				delete(s.conns, ws)
			} 
		}(ws)
	}
}

func main() {
	server := NewServerWS()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/ws/orderbook", websocket.Handler(server.handleWSOrderBook))
	http.ListenAndServe(":3000", nil)
}