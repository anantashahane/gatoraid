package main

import (
	"fmt"

	"github.com/anantashahane/gatoraid/internal"
	"github.com/anantashahane/gatoraid/internal/config"
)

type state struct {
	configuration *config.Config
}

type command struct {
	command    string
	arguements []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) (err error) {
	if len(cmd.arguements) != 1 {
		return fmt.Errorf("Login expects exactly 1 argument, received %v", cmd.arguements)
	}
	err = internal.SetUser(*s.configuration, cmd.arguements[0])
	if err != nil {
		return err
	}
	data, err := internal.Read()
	if err != nil {
		return err
	}
	s.configuration = &data
	fmt.Println("Logged in as " + s.configuration.CurrentUserName)
	return nil
}

func (c *commands) run(s *state, cmd command) (err error) {
	runner, exist := c.commandMap[cmd.command]
	if !exist {
		return fmt.Errorf("No such command as %s", cmd.command)
	}
	err = runner(s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}
