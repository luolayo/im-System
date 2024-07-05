package model

type Message struct {
	User    *User  // User is the sender of the message
	Content string // Content is the message content
}
