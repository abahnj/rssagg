# RSS Aggregator CLI (rssagg)

A command-line RSS feed aggregator written in Go that allows you to follow, manage, and browse content from your favorite RSS feeds.

## Features

- 📰 **Feed Management**: Add, list, follow, and unfollow RSS feeds
- 👥 **User Management**: Register users and manage authentication
- 🔄 **Automated Aggregation**: Continuously fetch and update content from followed feeds
- 📱 **Content Browsing**: View aggregated posts from your followed feeds
- 🔍 **Smart Duplicates Handling**: Automatically detects and prevents duplicate posts

## Prerequisites

- [Go](https://golang.org/doc/install) 1.19 or later
- [PostgreSQL](https://www.postgresql.org/download/) 14 or later

## Installation

### Installing From Source

Clone the repository and build from source:

```bash
# Clone the repository
git clone https://github.com/abahnj/rssagg.git
cd rssagg

# Build the binary
go build -o rssagg

# Optionally, install to your GOPATH
go install
```

### Using Go Install

```bash
go install github.com/abahnj/rssagg@latest
```

## Configuration

The application requires a database connection. Create a `config.json` file in the project root directory (you can copy from the provided `config.json.example`):

Example configuration:

```json
{
  "db_url": "postgres://username:password@localhost:5432/rssagg",
  "current_user_name": ""
}
```


## Database Setup

First, create a new PostgreSQL database:

```bash
# Option 1: Using createdb (Unix/macOS)
createdb rssagg

# Option 2: Using psql (works on all platforms)
psql -c "CREATE DATABASE rssagg;"

# Option 3: Using pgAdmin or another GUI tool
# Right-click on "Databases" and select "Create" > "Database"
# Enter "rssagg" as the database name
```

### Option 1: Manual SQL Execution

You can manually set up the database schema using the SQL files in the `sql/schema/` directory.

### Option 2: Using Goose for Migrations (Recommended)

[Goose](https://github.com/pressly/goose) is a database migration tool that helps manage schema changes. To use Goose:

1. Install Goose:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. Run the migrations:
   ```bash
   cd /path/to/rssagg
   goose -dir sql/schema postgres "postgres://username:password@localhost:5432/rssagg" up
   ```

   Replace `username` and `password` with your PostgreSQL credentials.

## Usage

```bash
# Start the application
rssagg
```

### Available Commands

```
Available commands:
  login <username>        - Log in as a user
  register <username>     - Register a new user
  users                   - List all users
  reset                   - Delete all users
  feeds                   - List all feeds
  addfeed <name> <url>    - Add a new feed
  follow <url>            - Follow an existing feed
  unfollow <url>          - Unfollow a feed
  following               - List feeds you're following
  browse [limit]          - View posts from feeds you follow (default limit: 10)
  agg <duration>          - Aggregate and show feed content every <duration> (e.g. 30s, 1m)
```

## Examples

### Basic Workflow

```bash
# Register a new user
rssagg register johndoe

# Add a new feed
rssagg addfeed "Hacker News" https://hnrss.org/newest

# Follow the feed
rssagg follow https://hnrss.org/newest

# List followed feeds
rssagg following

# Browse posts (limit to 20)
rssagg browse 20

# Continuously aggregate content every 30 seconds
rssagg agg 30s
```

## Development

### Project Structure

```
├── commands.go              # Command definitions
├── internal/
│   ├── cli/                 # CLI framework
│   ├── config/              # Configuration management
│   ├── database/            # Database models and queries
│   ├── feeds/               # Feed management
│   ├── middleware/          # Request middleware
│   ├── posts/               # Post management
│   ├── types/               # Shared type definitions
│   └── users/               # User management
├── main.go                  # Application entry point
├── sql/
│   ├── queries/             # SQLC query definitions
│   └── schema/              # Database migrations
└── sqlc.yaml                # SQLC configuration
```

### Database Schema

The application uses the following tables:
- `users`: Stores user information
- `feeds`: Stores feed information and metadata
- `feed_follows`: Manages relationships between users and feeds
- `posts`: Stores posts from feeds

## Contributing

Contributions are welcome! Feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.