package main

import (
	"context"
	"fmt"
)

func handlerAgg(state *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %+v\n", feed)

	return nil
}
