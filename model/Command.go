package model

type Command int

const (
	CmdMsg Command = iota
	CmdExit
	CmdList
)

func ParseCommand(msg string) Command {
	switch msg {
	case "exit":
		return CmdExit
	case "list":
		return CmdList
	default:
		return CmdMsg
	}
}
