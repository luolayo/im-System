package main

import Server "im-System/server"

func main() {
	server := Server.NewServer("0.0.0.0", 30001)
	server.Start()
}
