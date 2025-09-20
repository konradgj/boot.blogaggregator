package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/konradgj/boot.blogaggregator/internal/config"
	"github.com/konradgj/boot.blogaggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	st := state{
		cfg: &cfg,
		db:  dbQueries,
	}
	cmds := commands{
		registeredCmds: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args
	if len(args) == 1 {
		fmt.Println("expected a command: <command>")
		os.Exit(1)
	}
	if len(args) < 2 {
		fmt.Println("expected argument: command <arg>")
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
