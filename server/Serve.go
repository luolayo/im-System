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

// listenMessage listens for messages and broadcasts them to all online users
func (s *Server) listenMessage() {
	for msg := range s.Message {
		s.MapLock.Lock()
		for _, user := range s.Online {
			user.C <- msg
		}
		s.MapLock.Unlock()
	}
}

// broadcast sends a message to the Message channel
func (s *Server) broadcast(user *model.User, msg string) {
	s.Message <- fmt.Sprintf("%s: %s", user.Name, msg)
}

// handler handles new connections
func (s *Server) handler(conn net.Conn) {
	user := model.NewUser(conn)
	s.userOnline(user)
	// If the user sends a message, send the message to everyone
	go s.sendMsg(user)
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
	go s.listenMessage()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handler(conn)
	}
}

// userOnline marks a user as online
func (s *Server) userOnline(user *model.User) {
	s.MapLock.Lock()
	defer s.MapLock.Unlock()
	s.Online[user.Name] = user
	log.Printf("User %s has connected", user.Name)
	s.broadcast(user, "has connected")
}

// userOffline marks a user as offline
func (s *Server) userOffline(user *model.User) {
	s.MapLock.Lock()
	defer s.MapLock.Unlock()
	delete(s.Online, user.Name)
	log.Printf("User %s has disconnected", user.Name)
	// Check if the connection is already closed before attempting to close it
	if user.Conn != nil {
		err := user.Conn.Close()
		if err != nil {
			log.Printf("Error closing connection for user %s", user.Name)
		}
	}
	s.broadcast(user, "has disconnected")
}

// sendMsg handles message sending from a user
func (s *Server) sendMsg(user *model.User) {
	buf := make([]byte, 4096)
	for {
		if user.Conn == nil {
			return
		}
		cnt, err := user.Conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from %s: %v", user.Name, err)
			s.userOffline(user)
			return
		}
		msg := string(buf[:cnt-1])
		s.handleUserMessage(user, msg)
	}
}

// listUsers returns a string of all online users
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

// handleUserMessage processes the message sent by the user
func (s *Server) handleUserMessage(user *model.User, msg string) {
	command := model.ParseCommand(msg)
	switch command {
	case model.CmdExit:
		log.Printf("User %s sent exit command", user.Name)
		s.userOffline(user)
	case model.CmdList:
		log.Printf("User %s sent list command", user.Name)
		user.C <- s.listUsers()
	case model.CmdRename:
		newName := strings.TrimSpace(strings.TrimPrefix(msg, "/rename "))
		if newName == "" {
			user.C <- "Invalid new name."
		} else {
			oldName := user.Name
			s.renameUser(user, newName)
			log.Printf("User %s renamed to %s", oldName, newName)
			s.broadcast(user, fmt.Sprintf("has renamed to %s", newName))
		}
	default:
		s.broadcast(user, msg)
	}
}

// renameUser changes the username of a user
func (s *Server) renameUser(user *model.User, newName string) {
	s.MapLock.Lock()
	defer s.MapLock.Unlock()
	delete(s.Online, user.Name)
	user.Name = newName
	s.Online[newName] = user
}
