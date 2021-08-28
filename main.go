package main

import (
	"log"
	"os"
)

type Command struct {
}

func NewCommand() *Command {
	return &Command{}
}

func (c *Command) Run() error {
	return nil
}

func main() {
	command := NewCommand()
	if err := command.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
