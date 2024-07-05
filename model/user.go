package model

import (
	Logger "im-System/logger"
	"net"
)

// User represents a connected user
type User struct {
	Conn    net.Conn // Conn is the network connection associated with the user
	name    string   // Name is the user's name
	address string   // Address is the user's network address
}

// NewUser creates a new User instance
func NewUser(conn net.Conn, name string) *User {
	return &User{
		Conn:    conn,
		name:    name,
		address: conn.RemoteAddr().String(),
	}
}

func (u *User) Name() string {
	if u.name == "" {
		return u.address
	}
	return u.name
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) Address() string {
	return u.address
}

// Close closes the user's connection
func (u *User) Close() {
	if u.Conn != nil {
		err := u.Conn.Close()
		if err != nil {
			logger := Logger.NewLogger(Logger.ErrorLevel)
			logger.Error("Error closing user connection: %v", err)
		}
	}
}
