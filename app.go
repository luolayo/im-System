package main

import (
	"fmt"
	"im-System/client"
	Server "im-System/server"
	"im-System/util"
)

func main() {
	fmt.Printf("Enter '1' to start the Server or '2' to start the client:")
	choice := util.InputString()
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
	fmt.Printf("Enter IP address (default 127.0.0.1): ")
	ip := util.InputString()
	if ip == "" {
		ip = "127.0.0.1"
	}
	fmt.Printf("Enter port (default 30001): ")
	port := util.InputString()
	if port == "" {
		port = "30001"
	}
	newServer := Server.NewServer(ip, port)
	newServer.Start()
}

func startClient() {
	fmt.Printf("Enter Server address (default 127.0.0.1:30001): ")
	serverAddr := util.InputString()
	if serverAddr == "" {
		serverAddr = "127.0.0.1:30001"
	}
	newClient := client.NewClient(serverAddr)
	newClient.Start()
}
