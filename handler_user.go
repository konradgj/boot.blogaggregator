package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/konradgj/boot.blogaggregator/internal/database"
	"github.com/lib/pq"
)

func handlerLogin(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected argument: login <name>")
	}
	name := cmd.args[0]

	user, err := state.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user %s does not exist", name)
	}

	err = state.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Logged in as %s\n", name)

	return nil
}

func handlerRegister(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected argument: register <name>")
	}
	name := cmd.args[0]

	user, err := state.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		})
	if err != nil {
		// check for unique violation
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("user already exists with name %s", name)
		}
		return fmt.Errorf("error registering user: %w", err)
	}

	state.cfg.CurrentUserName = name
	fmt.Printf("Succefully registered user %s\n", name)

	err = state.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	return nil
}
