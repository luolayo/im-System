package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	Conn net.Conn
	Name string
}

// NewClient creates a new client and connects to the server
func NewClient(serverAddr string) *Client {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
	}
	return &Client{
		Conn: conn,
	}
}

// SendMessage sends a message to the server
func (c *Client) SendMessage(msg string) {
	_, err := c.Conn.Write([]byte(msg + "\n"))
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// ReceiveMessages receives messages from the server
func (c *Client) ReceiveMessages() {
	reader := bufio.NewReader(c.Conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}
		fmt.Print(msg)
	}
}

func main() {
	serverAddr := "127.0.0.1:30001" // Server address
	client := NewClient(serverAddr)

	// Prompt for client name
	fmt.Print("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	client.Name = strings.TrimSpace(name)

	// Start receiving messages
	go client.ReceiveMessages()

	// Send messages to server
	for {
		fmt.Print("> ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}
		client.SendMessage(msg)
		if msg == "exit" {
			break
		}
	}
}
