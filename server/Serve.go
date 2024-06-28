package server

import (
	"fmt"
	"im-System/model"
	"net"
	"strconv"
	"sync"
)

// Server This is a struct that will hold the server configuration
type Server struct {
	// Ip string
	Ip string
	// Port int
	Port int
	// Online map[string]*model.User
	Online map[string]*model.User
	// MapLock sync.RWMutex
	MapLock sync.RWMutex
	// Message chan string
	Message chan string
}

// NewServer This is a constructor for the Serve struct
func NewServer(ip string, port int) *Server {
	fmt.Println("Server is starting")
	return &Server{
		Ip:      ip,
		Port:    port,
		Online:  make(map[string]*model.User),
		Message: make(chan string),
	}
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.MapLock.Lock()
		for _, user := range s.Online {
			user.C <- msg
		}
		s.MapLock.Unlock()
	}
}

func (s *Server) Broadcast(user model.User, msg string) {
	s.Message <- user.Name + ":" + msg
}

// Handler This is a method that will handle the connection
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("A new connection has been established")
	user := model.NewUser(conn)
	s.MapLock.Lock()
	s.Online[user.Name] = user
	s.MapLock.Unlock()
	s.Broadcast(*user, "has connected")
	// If the user sends a message, send the message to everyone
	go func() {
		buf := make([]byte, 4096)
		for {
			cnt, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				return
			}
			if cnt == 0 {
				s.Broadcast(*user, "has disconnected")
			}
			msg := string(buf[:cnt-1])
			s.Broadcast(*user, msg)
		}
	}()
	select {}
}

// Start This is a method that will start the server
func (s *Server) Start() {
	listen, err := net.Listen("tcp", s.Ip+":"+strconv.Itoa(s.Port))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			fmt.Println("Error closing listener:", err.Error())
			return
		}
	}(listen)
	go s.ListenMessage()
	for {
		conn, err := listen.Accept()
		if err != nil {
			return
		}
		go s.Handler(conn)
	}
}
