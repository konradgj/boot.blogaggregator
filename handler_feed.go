package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/konradgj/boot.blogaggregator/internal/database"
)

func handlerAddFeed(state *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("expected two arguments: register <name> <url>")
	}

	user, err := state.db.GetUser(context.Background(), state.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feed, err := state.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
			Url:       cmd.args[1],
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println("=====================================")

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}
