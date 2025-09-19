package main

import (
	"fmt"
	"log"
	"os"

	"github.com/konradgj/boot.blogaggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	st := state{
		cfg: &cfg,
	}
	cmds := commands{
		registeredCmds: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) == 1 {
		fmt.Println("expected a command: <command>")
		os.Exit(1)
	}
	if len(args) < 2 {
		fmt.Println("username required: login <name>")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}
	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
