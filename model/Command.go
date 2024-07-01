package model

import "strings"

// Command type to define different commands
type Command int

const (
	CmdMsg Command = iota
	CmdExit
	CmdList
	CmdRename
)

// ParseCommand parses a message and returns the corresponding command
func ParseCommand(msg string) Command {
	switch {
	case strings.HasPrefix(msg, "/rename "):
		return CmdRename
	case msg == "exit":
		return CmdExit
	case msg == "list":
		return CmdList
	default:
		return CmdMsg
	}
}
