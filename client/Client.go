package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Client struct
type Client struct {
	conn   net.Conn
	name   string
	server string
}

// NewClient creates a new client
func NewClient(serverAddr string) *Client {
	return &Client{
		server: serverAddr,
	}
}

// Connect to the server
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.server)
	if err != nil {
		return fmt.Errorf("error connecting to server: %v", err)
	}
	c.conn = conn

	// Start a goroutine to handle incoming messages
	go c.handleServerMessages()

	return nil
}

// handleServerMessages listens for messages from the server and prints them
func (c *Client) handleServerMessages() {
	for {
		message := make([]byte, 1024)
		n, err := c.conn.Read(message)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Print("\r" + string(message[:n]) + "> ") // Print the server message and reprint the prompt
		if strings.Contains(string(message[:n]), "disconnected") {
			os.Exit(0)
		}
	}
}

// SendMessage sends a message to the server
func (c *Client) SendMessage(message string) {
	_, err := c.conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

// StartClient starts the client and displays the menu
func (c *Client) StartClient() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Exit")
		fmt.Println("2. Rename User")
		fmt.Println("3. Enter Public Chat Mode")
		fmt.Println("4. Query all online users")
		fmt.Print("Enter your choice: ")

		if scanner.Scan() {
			choice := scanner.Text()
			switch choice {
			case "1":
				fmt.Println("Exiting...")
				c.SendMessage("/exit")
				return
			case "2":
				fmt.Print("Enter new username: ")
				if scanner.Scan() {
					newName := scanner.Text()
					c.SendMessage("/rename " + newName)
				}
			case "3":
				c.enterPublicChatMode()
			case "4":
				c.SendMessage("/users")
			default:
				fmt.Println("Invalid choice. Please try again.")
			}
		}
	}
}

// enterPublicChatMode handles the public chat mode
func (c *Client) enterPublicChatMode() {
	fmt.Println("Entering public chat mode. Type '/exit' to return to the menu.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			text := scanner.Text()
			if text == "/exit" {
				return
			}
			c.SendMessage(text)
		}
	}
}
