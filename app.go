package main

import (
	"bufio"
	"fmt"
	"im-System/client"
	"im-System/server"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Enter '1' to start the server or '2' to start the client:")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	switch choice {
	case "1":
		startServer()
	case "2":
		startClient()
	default:
		fmt.Println("Invalid choice")
	}
}

func startServer() {
	fmt.Print("Enter IP address (default 127.0.0.1): ")
	reader := bufio.NewReader(os.Stdin)
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)
	if ip == "" {
		ip = "127.0.0.1"
	}

	fmt.Print("Enter port (default 30001): ")
	portInput, _ := reader.ReadString('\n')
	portInput = strings.TrimSpace(portInput)
	port := 30001
	if portInput != "" {
		port, _ = strconv.Atoi(portInput)
	}

	newServer := server.NewServer(ip, port)
	newServer.Start()
}

func startClient() {
	fmt.Print("Enter server address (default 127.0.0.1:30001): ")
	reader := bufio.NewReader(os.Stdin)
	serverAddr, _ := reader.ReadString('\n')
	serverAddr = strings.TrimSpace(serverAddr)
	if serverAddr == "" {
		serverAddr = "127.0.0.1:30001"
	}

	newClient := client.NewClient(serverAddr)

	// Prompt for newClient name
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	newClient.Name = strings.TrimSpace(name)

	// Automatically rename the user after connection
	newClient.SendMessage("/rename " + newClient.Name)

	// Start receiving messages
	go newClient.ReceiveMessages()

	// Send messages to server
	for {
		fmt.Print("> ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}
		newClient.SendMessage(msg)
		if msg == "exit" {
			break
		}
	}
}
