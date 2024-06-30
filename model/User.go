package model

import (
	"log"
	"net"
)

// User holds the user configuration
type User struct {
	Name    string      // User's name
	Address string      // User's address
	C       chan string // Channel for user messages
	Conn    net.Conn    // User's connection
}

// NewUser is a constructor for the User struct
func NewUser(conn net.Conn) *User {
	user := &User{
		Name:    conn.RemoteAddr().String(),
		Address: conn.RemoteAddr().String(),
		C:       make(chan string),
		Conn:    conn,
	}
	go user.ListenMessage()
	return user
}

// ListenMessage listens for messages on the user's channel and sends them to the connection
func (u *User) ListenMessage() {
	for msg := range u.C {
		_, err := u.Conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			return
		}
	}
}
