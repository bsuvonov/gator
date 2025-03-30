# Gator

Gator is a multi-user command-line RSS feed aggregator written in Go. It allows users to add and follow RSS feeds, aggregate the posts, and browse them efficiently. The program utilizes a PostgreSQL database for storage and supports user authentication to manage personal feed subscriptions.

## Features

- **User Authentication**: Register, login, and manage user sessions.
- **RSS Feed Management**: Add and follow RSS feeds.
- **Feed Aggregation**: Periodically fetch RSS feed data.
- **Post Browsing**: View recent posts from followed feeds.
- **Database-backed Storage**: Uses PostgreSQL for data persistence.

## Installation

### Prerequisites
Ensure you have the following installed:

- Go (1.24+)
- PostgreSQL

### Clone the Repository
```sh
git clone https://github.com/bsuvonov/gator.git
cd gator
```

### Set Up Environment Variables
Create a configuration file in your home directory:
```sh
touch ~/.gatorconfig.json
```
Edit the file and add:
```json
{
    "db_url": "postgres://user:password@hostname:5432/gator",
    "current_user_name": ""
}
```
Replace `user`, `hostname`, `password`, and `gator` with your PostgreSQL credentials.

### Run database migrations

```sh
# Install goose to handle database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run `up` migrations
goose postgres "postgres://user:password@hostname:5432/gator" --dir sql/schema up
```
Replace the PosgtreSQL connection string with a correct one.

### Build and Run
To build the project:
```sh
go mod tidy
go build -o gator
```
Run the program:
```sh
./gator <command> [args]
```

## Commands

| Command         | Description |
|----------------|-------------|
| `register <name>` | Register a new user |
| `login <name>` | Login to existing user |
| `reset` | Reset users and their feeds |
| `users` | List existing users |
| `addfeed <name> <url>` | Add a new RSS feed |
| `feeds` | List all created feeds |
| `follow <url>` | Follow a specific feed |
| `following` | List feeds the user follows |
| `unfollow <url>` | Unfollow a feed |
| `agg <duration>` | Fetch new posts from subscribed feeds periodically |
| `browse [limit]` | Browse recently published posts (default: 5 posts) |

## Example Usage

1. Register a user:
```sh
./gator register alice
```

2. Login as the user:
```sh
./gator login alice
```

3. Add an RSS feed:
```sh
./gator addfeed TechCrunch https://techcrunch.com/feed/
```

4. Follow an existing feed:
```sh
./gator follow https://techcrunch.com/feed/
```

5. Start fetching feeds every 10 minutes:
```sh
./gator agg 10m
```

6. Browse the latest posts:
```sh
./gator browse 10
```

## Database Schema
This application uses PostgreSQL with the following main tables:

- `users`: Stores registered users.
- `feeds`: Stores added RSS feeds.
- `feed_follows`: Tracks which users follow which feeds.
- `posts`: Stores posts from feeds.
