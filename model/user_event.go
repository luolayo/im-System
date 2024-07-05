package model

const (
	// UserJoin is the event type for a user joining
	UserJoin = "join"
	// UserLeave is the event type for a user leaving
	UserLeave = "leave"
	// UserMessage is the event type for a user message
	UserMessage = "message"
	// UserList is the event type for a user list
	UserList = "users"
	// UserRename is the event type for a user renaming
	UserRename = "rename"
)

// UserEvent represents events related to users
type UserEvent struct {
	Type string // Type is the event type ("join", "leave", etc.)
	User *User  // User is the user associated with the event
}
