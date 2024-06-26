package model

import (
	"fmt"
	"net"
)

// User This is a struct that will hold the user configuration
type User struct {
	// Name string
	Name string
	// Address string
	Address string
	// C chan string
	C chan string
	// conn net.Conn
	Conn net.Conn
}

// NewUser This is a constructor for the User struct
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

// ListenMessage This is a method that will listen for messages
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		_, err := u.Conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Error writing to connection:", err.Error())
			return
		}
	}
}
