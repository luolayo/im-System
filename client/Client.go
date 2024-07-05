package client

import (
	"bufio"
	"fmt"
	Logger "im-System/logger"
	"im-System/util"
	"net"
	"strings"
)

type Client struct {
	Conn        net.Conn
	Name        string
	MessageChan chan string // Channel to notify message received
	logger      Logger.Logger
}

// NewClient creates a new client and connects to the Server
func NewClient(serverAddr string) *Client {
	logger := *Logger.NewLogger(Logger.ErrorLevel)
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		logger.Error("Error connecting to server: %v", err)
	}
	return &Client{
		Conn:        conn,
		MessageChan: make(chan string), // Initialize the message channel
		logger:      logger,
	}
}

// SendMessage sends a message to the Server
func (c *Client) SendMessage(msg string) {
	_, err := c.Conn.Write([]byte(msg + "\n"))
	if err != nil {
		c.logger.Error("Error sending message: %v", err)
	}
}

// ReceiveMessages receives messages from the Server
func (c *Client) ReceiveMessages(done chan struct{}) {
	reader := bufio.NewReader(c.Conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			c.logger.Error("Error receiving message: %v", err)
			close(done)
			return
		}
		fmt.Print("\r" + msg + "> ") // Print the server message and reprint the prompt
		// Check if the message indicates that the user has been kicked out
		if strings.TrimSpace(msg) == "You have been inactive for too long. You will be disconnected." {
			fmt.Println("\nYou have been inactive for too long. You will be disconnected.")
			close(done)
			return
		}
		c.MessageChan <- "" // Notify the main routine to reprint the prompt
	}
}

func (c *Client) Start() {
	// Prompt for newClient name
	fmt.Print("Enter your name: ")
	name := util.InputString()
	c.Name = strings.TrimSpace(name)

	// Automatically rename the user after connection
	c.SendMessage("/rename " + c.Name)

	done := make(chan struct{})

	// Start receiving messages
	go c.ReceiveMessages(done)

	// Send messages to Server
	for {
		fmt.Print("> ")
		msg := util.InputString()
		if msg == "" {
			continue
		}
		c.SendMessage(msg)
		if msg == "/exit" {
			c.SendMessage("/exit")
			break
		}
		// Wait for a message to be received or done channel to be closed before printing the next prompt
		select {
		case <-c.MessageChan:
			// Message received, continue to next iteration
		case <-done:
			// Done signal received, break the loop
			return
		}
	}
}
