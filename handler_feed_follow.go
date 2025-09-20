package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/konradgj/boot.blogaggregator/internal/database"
)

func handlerAddFeedFollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected one argument: follow <url>")
	}
	url := cmd.args[0]

	feed, err := state.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}

	feedFollow, err := addFeedFollow(state, feed.ID, user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s now following %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerGetFeedFollowsForUser(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("expected no arguments: following")
	}

	feedFollows, err := state.db.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching feedfollows: %w", err)
	}

	fmt.Printf("%s follows:\n", feedFollows[0].UserName)
	for _, feedFollow := range feedFollows {
		fmt.Printf("- %s\n", feedFollow.FeedName)
	}

	return nil
}

func addFeedFollow(state *state, feed_id uuid.UUID, user_id uuid.UUID) (database.CreateFeedFollowRow, error) {
	var feedFollow database.CreateFeedFollowRow
	feedFollow, err := state.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed_id,
			UserID:    user_id,
		},
	)
	if err != nil {
		return feedFollow, fmt.Errorf("error following feed: %w", err)
	}

	return feedFollow, nil
}
