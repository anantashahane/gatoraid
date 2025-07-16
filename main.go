package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/anantashahane/gatoraid/internal"
	"github.com/anantashahane/gatoraid/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	data, err := internal.Read()
	if err != nil {
		fmt.Println("Error recovring data file.", err)
		os.Exit(1)
	}
	programState := state{configuration: &data}
	db, err := sql.Open("postgres", programState.configuration.DbURL)
	if err != nil {
		fmt.Println("Error opening up database.")
		os.Exit(1)
	}
	programState.dbConnection = database.New(db)

	cmds := commands{commandMap: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", resetData)
	cmds.register("users", presentAllUsers)
	cmds.register("agg", aggHandler)
	cmds.register("addfeed", addFeed)
	cmds.register("feeds", getAllFeed)

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
