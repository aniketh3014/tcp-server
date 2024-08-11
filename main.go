package main

import (
	"fmt"
	"net"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	quitchan   chan struct{}
	magchan    chan []byte
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitchan:   make(chan struct{}),
		magchan:    make(chan []byte, 10),
	}
}

func (s *Server) start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.listener = ln

	go s.acceptLoop()

	<-s.quitchan
	close(s.magchan)
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("can not accept connection", err)
			continue
		}
		fmt.Println("accepted new connection", conn.RemoteAddr())
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {

	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading messege", err)
			continue
		}

		s.magchan <- buf[:n]

		conn.Write([]byte("Thank you for your message!"))
	}
}

func main() {
	server := NewServer(":8080")
	go func() {
		for msg := range server.magchan {
			fmt.Println("message recived from connection:", string(msg))
		}
	}()

	server.start()
}
