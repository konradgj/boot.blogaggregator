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
		commands: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerResetUsers)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", handlerAddFeedFollow)
	cmds.register("following", handlerGetFeedFollowsForUser)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		os.Exit(1)
	}
	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&st, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
