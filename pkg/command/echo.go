package command

import "fmt"

type EchoCommand struct {
	message string
}

func Echo(message string) *EchoCommand {
	return &EchoCommand{message: message}
}

func (cmd *EchoCommand) Encode() string {
	return fmt.Sprintf("ECHO \"%s\"", cmd.message)
}
