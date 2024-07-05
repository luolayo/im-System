package server

import (
	"fmt"
	Logger "im-System/logger"
	"im-System/model"
	"net"
	"strings"
	"sync"
)

// Server represents the chat server
type Server struct {
	ip         string                   // IP is the server's IP address
	port       string                   // Port is the server's port number
	mu         sync.Mutex               // mu is used to protect the clients map
	clients    map[net.Conn]*model.User // clients maps connections to users
	listener   net.Listener             // listener is the TCP listener
	quit       chan struct{}            // quit is used to signal server shutdown
	userEvents chan model.UserEvent     // userEvents is used to send user-related events
	messages   chan model.Message       // messages is used to broadcast messages
	logger     Logger.Logger            // logger is the server logger
}

// NewServer creates a new Server instance
func NewServer(ip, port string) *Server {
	return &Server{
		ip:         ip,
		port:       port,
		clients:    make(map[net.Conn]*model.User),
		quit:       make(chan struct{}),
		userEvents: make(chan model.UserEvent),
		messages:   make(chan model.Message),
		logger:     *Logger.NewLogger(Logger.InfoLevel),
	}
}

// Start starts the server
func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.ip, s.port))
	if err != nil {
		s.logger.Error("Error starting server: %s", err)
	}
	s.listener = listener
	s.logger.Info("Starting server on %s:%s", s.ip, s.port)
	// Start the goroutine to handle incoming connections
	go s.acceptConnections()

	// Start the goroutine to handle user events
	go s.handleUserEvents()

	// Start the goroutine to broadcast messages
	go s.broadcastMessages()

	// Wait for a signal to quit
	<-s.quit
}

// Stop stops the server
func (s *Server) Stop() {
	close(s.quit)
	err := s.listener.Close()
	if err != nil {
		s.logger.Error("Error stopping server: %s", err)
	}
}

// acceptConnections accepts incoming client connections
func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("Error accepting connection:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection handles individual client connections
func (s *Server) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		if conn == nil {
			return
		}
		err := conn.Close()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			s.logger.Error("Error closing connection: %s", err)
		}
	}(conn)

	// For now, just add the user to the list of clients
	user := model.NewUser(conn, "")
	s.addUser(user)

	// Read messages from the connection and broadcast them
	buf := make([]byte, 1024)
	for {
		if conn == nil {
			s.removeUser(conn)
			return
		}
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				s.removeUser(conn)
				return
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				s.removeUser(conn)
				return
			}
			s.logger.Error("Error reading from connection: %s", err)
			s.removeUser(conn)
			return
		}
		s.handleUserMessage(user, buf[:n-1])
	}
}

func (s *Server) broadcastMessages() {
	for {
		select {
		case msg := <-s.messages:
			s.mu.Lock()
			for conn, user := range s.clients {
				if user != msg.User {
					_, err := conn.Write([]byte(fmt.Sprintf("%s: %s", msg.User.Name(), msg.Content)))
					if err != nil {
						s.logger.Error("Error writing to connection: %s", err)
					}
				}
			}
			s.mu.Unlock()
		case <-s.quit:
			return
		}
	}
}

// renameUser renames a user and notifies all users
func (s *Server) renameUser(user *model.User, newName string) {
	oldName := user.Name()
	user.SetName(newName)

	// Notify the user who changed the name
	confirmationMessage := fmt.Sprintf("You have successfully changed your name to %s\n", newName)
	_, err := user.Conn.Write([]byte(confirmationMessage))
	if err != nil {
		s.logger.Error("Error writing to connection: %s", err)
	}

	// Notify all other users about the name change
	s.mu.Lock()
	defer s.mu.Unlock()
	notification := fmt.Sprintf("%s has changed their name to %s\n", oldName, newName)
	for conn := range s.clients {
		if conn != user.Conn {
			_, err := conn.Write([]byte(notification))
			if err != nil {
				s.logger.Error("Error writing to connection: %s", err)
			}
		}
	}
	s.userEvents <- model.UserEvent{Type: model.UserRename, User: user}
}

// addUser adds a user to the server's clients
func (s *Server) addUser(user *model.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[user.Conn] = user
	s.userEvents <- model.UserEvent{Type: model.UserJoin, User: user}
	_, err := user.Conn.Write([]byte("Welcome to the chat!\n"))
	if err != nil {
		s.logger.Error("Error writing to connection: %s", err)
		return
	}
}

// removeUser removes a user from the server's clients
func (s *Server) removeUser(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user, ok := s.clients[conn]; ok {
		delete(s.clients, conn)
		s.userEvents <- model.UserEvent{Type: model.UserLeave, User: user}
	}
}

// handleUserEvents handles user-related events
func (s *Server) handleUserEvents() {
	for {
		select {
		case event := <-s.userEvents:
			switch event.Type {
			case model.UserJoin:
				s.logger.Info("%s joined the chat", event.User.Name())
			case model.UserLeave:
				s.logger.Info("%s left the chat", event.User.Name())
			case model.UserMessage:
				s.logger.Info("%s send a message", event.User.Name())
			case model.UserList:
				s.logger.Info("%s List of users requested", event.User.Name())
			case model.UserRename:
				s.logger.Info("%s renamed", event.User.Name())
			}
		case <-s.quit:
			return
		}
	}
}

func (s *Server) handleUserMessage(user *model.User, msg []byte) {
	message := string(msg)
	if message == "/exit" {
		s.removeUser(user.Conn)
		user.Close() // Close the user's connection
		return
	}
	if message == "/users" {
		s.listUsers(user)
		return
	}
	if strings.HasPrefix(message, "/rename") {
		newName := strings.TrimPrefix(message, "/rename ")
		s.renameUser(user, newName)
		return
	}
	// Add line breaks to user messages
	message += "\n"
	s.userEvents <- model.UserEvent{Type: model.UserMessage, User: user}
	s.messages <- model.Message{User: user, Content: message}
}

// listUsers sends the list of online users to the requesting user
func (s *Server) listUsers(requestingUser *model.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userEvents <- model.UserEvent{Type: model.UserList, User: requestingUser}
	var users []string
	for _, user := range s.clients {
		users = append(users, user.Name())
	}

	userList := strings.Join(users, ", ")
	message := fmt.Sprintf("Online users: %s\n", userList)
	_, err := requestingUser.Conn.Write([]byte(message))
	if err != nil {
		s.logger.Error("Error writing to connection: %s", err)
	}
}
