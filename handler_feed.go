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
	printFeed(feed, user)
	fmt.Println("=====================================")

	return nil
}

func handlerGetFeeds(state *state, _ command) error {

	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds in database")
		return nil
	}

	for _, feed := range feeds {
		user, err := state.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting user: %w", err)
		}

		printFeed(feed, user)
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
