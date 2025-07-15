package main

import (
	"fmt"
	"os"

	"github.com/anantashahane/gatoraid/internal"
)

func main() {
	data, err := internal.Read()
	if err != nil {
		fmt.Println("Error recovring data file.", err)
		os.Exit(1)
	}
	programState := state{configuration: &data}

	cmds := commands{commandMap: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Expected atleast one arguement.")
		for k, _ := range cmds.commandMap {
			fmt.Println("\t", k)
		}
		os.Exit(1)
	}

	cmd := command{command: args[1], arguements: args[2:]}
	if err = cmds.run(&programState, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
