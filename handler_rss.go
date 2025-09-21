package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/konradgj/boot.blogaggregator/internal/database"
	"github.com/lib/pq"
)

func handlerBrowse(state *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("expected one or no arguments: browse <amount>(optional)")
	}
	limit := 2
	if len(cmd.args) == 1 {
		l, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = l
	}

	posts, err := state.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
		},
	)
	if err != nil {
		return fmt.Errorf("error fetching posts: %w", err)
	}

	for _, post := range posts {
		printPost(post)
	}

	return nil
}

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

	now := time.Now().UTC()

	for _, rssItem := range rssFeed.Channel.Item {
		pubDate, err := parseXmlDate(rssItem.PubDate)
		if err != nil {
			log.Printf("could not parse pubDate for %s: %v", rssItem.Title, err)
			pubDate = time.Time{}
		}

		_, err = state.db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   now,
				UpdatedAt:   now,
				Title:       rssItem.Title,
				Url:         rssItem.Link,
				Description: wrapNullString(rssItem.Description),
				PublishedAt: wrapNullTime(pubDate),
				FeedID:      feedToFetch.ID,
			},
		)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				continue
			}
			log.Printf("error creating post for %s: %v", rssItem.Title, err)
			continue
		}
	}

	return nil
}

func wrapNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func wrapNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func parseXmlDate(date string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, date); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse date: %s", date)
}

func printPost(post database.Post) {
	fmt.Printf("Title:        %s\n", post.Title)
	fmt.Printf("URL:          %s\n", post.Url)

	if post.Description.Valid {
		fmt.Printf("Description:  %s\n", post.Description.String)
	} else {
		fmt.Printf("Description:  <NULL>\n")
	}

	if post.PublishedAt.Valid {
		fmt.Printf("PublishedAt:  %s\n", post.PublishedAt.Time.Format(time.RFC3339))
	} else {
		fmt.Printf("PublishedAt:  <NULL>\n")
	}
	fmt.Println("--------------------------------------")
}
