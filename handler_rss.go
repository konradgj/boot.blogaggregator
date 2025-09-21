package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/konradgj/boot.blogaggregator/internal/database"
)

func handlerAgg(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected one argument: agg <time_between_reqs>")
	}

	interval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not parse interval: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", interval)

	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		err = scrapeFeed(state)
		if err != nil {
			log.Printf("error scraping feed: %v", err)
		}
	}
}

func scrapeFeed(state *state) error {
	feedToFetch, err := state.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not get next feed to fetch: %w", err)
	}

	fmt.Printf("Fetching feed %s\n", feedToFetch.Name)
	rssFeed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}

	err = state.db.MarkFeedFetched(
		context.Background(),
		database.MarkFeedFetchedParams{
			ID:        feedToFetch.ID,
			UpdatedAt: time.Now().UTC(),
		},
	)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %w", err)
	}

	fmt.Println("Found the following items:")
	for _, rssItem := range rssFeed.Channel.Item {
		fmt.Printf("- %s\n", rssItem.Title)
	}

	return nil
}
