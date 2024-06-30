package server

import (
	"fmt"
	"im-System/model"
	"log"
	"net"
	"strings"
	"sync"
)

// Server holds the server configuration
type Server struct {
	Ip      string                 // Server IP address
	Port    int                    // Server port number
	Online  map[string]*model.User // Map of online users
	MapLock sync.RWMutex           // Read/Write mutex for the Online map
	Message chan string            // Channel for broadcasting messages
}

// NewServer is a constructor for the Server struct
func NewServer(ip string, port int) *Server {
	log.Println("Server is starting")
	return &Server{
		Ip:      ip,
		Port:    port,
		Online:  make(map[string]*model.User),
		Message: make(chan string),
	}
}

// ListenMessage listens for messages and broadcasts them to all online users
func (s *Server) ListenMessage() {
	for msg := range s.Message {
		s.MapLock.Lock()
		for _, user := range s.Online {
			user.C <- msg
		}
		s.MapLock.Unlock()
	}
}

// Broadcast sends a message to the Message channel
func (s *Server) Broadcast(user *model.User, msg string) {
	s.Message <- fmt.Sprintf("%s: %s", user.Name, msg)
}

// Handler This is a method that will handle the connection
func (s *Server) Handler(conn net.Conn) {
	user := model.NewUser(conn)
	s.userOnline(user)
	// If the user sends a message, send the message to everyone
	go s.sendMsg(conn, user)
	select {}
}

// Start starts the server
func (s *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Fatalf("Error listening: %v", err)
		return
	}
	defer func() {
		err := listen.Close()
		if err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}()

	// Start listening for incoming messages
	go s.ListenMessage()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.Handler(conn)
	}
}

// userOnline marks a user as online
func (s *Server) userOnline(user *model.User) {
	s.MapLock.Lock()
	defer s.MapLock.Unlock()
	s.Online[user.Name] = user
	log.Printf("User %s has connected", user.Name)
	s.Broadcast(user, "has connected")
}

// userOffline marks a user as offline
func (s *Server) userOffline(user *model.User) {
	s.MapLock.Lock()
	defer s.MapLock.Unlock()
	delete(s.Online, user.Name)
	log.Printf("User %s has disconnected", user.Name)
	err := user.Conn.Close()
	if err != nil {
		return
	}
	s.Broadcast(user, "has disconnected")
}

// sendMsg Encapsulation message sending method
func (s *Server) sendMsg(conn net.Conn, user *model.User) {
	buf := make([]byte, 4096)
	for {
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		msg := string(buf[:cnt-1])
		if msg == "exit" {
			s.userOffline(user)
			return
		}
		if msg == "list" {
			user.C <- s.listUsers()
			continue
		}
		s.Broadcast(user, msg)
	}
}

func (s *Server) listUsers() string {
	// Print out all online users
	s.MapLock.Lock()
	defer s.MapLock.Unlock()
	var users []string
	for _, user := range s.Online {
		users = append(users, user.Name)
	}
	return "Online users: " + strings.Join(users, ", ")
}
