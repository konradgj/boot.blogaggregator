# boot.blogaggregator
boot.dev blog aggregator project


## Learning Goals

 - Learn how to integrate a Go application with a PostgreSQL database
 - Practice using your SQL skills to query and migrate a database (using sqlc and goose, two lightweight tools for typesafe SQL in Go)
 - Learn how to write a long-running service that continuously fetches new posts from RSS feeds and stores them in the database

# Gator

Gator is a command-line RSS aggregator built in Go. It stores feeds, posts, and user information in PostgreSQL and provides commands for managing and browsing feeds.

---

## Requirements

- **Go** (>= 1.25 recommended)  
  [Download & install Go](https://golang.org/dl/) if you don't have it already.
- **PostgreSQL**  
  Make sure you have a running Postgres instance and can connect using a user with privileges to create databases and tables.

---

## Install the CLI

You can install the `boot.blogaggregator` CLI tool with:

```bash
go install github.com/konradgj/boot.blogaggregatpr@latest
```

## Config
A config file is required to connect to your PostgreSQL database. Create a file named `.gatorconfig.json` in your `$HOME` directory:
```json
{
  "postgres_url": "postgres://username:password@localhost:5432/dbname?sslmode=disable"
}
```

# Usage
```bash
boot.bloggagregator <command> <argument>
```

| Command       | Arguments                     | Description                                      |
|---------------|-------------------------------|--------------------------------------------------|
| `register`    | `<username>`                    | Register a new user                            |
| `login`       | `<username>`                    | Log in as an existing user                     |
| `reset`       | None                          | Reset all users (wipes alle tables)              |
| `users`       | None                          | List all users                                   |
| `addfeed`     | `<name> <url>`                | Add a new RSS feed (requires login)              |
| `following`   | None                          | List feeds you are following                     |
| `feeds`       | None                          | List all available feeds                         |
| `agg`         | `<interval>` (optional)       | Start aggregator to fetch feeds periodically     |
| `browse`      | `<limit>` (optional)          | Browse posts from feeds for current user         |
| `follow`      | `<feed_id>`                   | Follow a feed                                    |
| `unfollow`    | `<feed_id>`                   | Stop following a feed                            |
