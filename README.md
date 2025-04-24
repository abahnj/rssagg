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

First, you'll need to create a new PostgreSQL database. There are several ways to do this depending on your setup:

### Creating the Database

```bash
# Option 1: Using createdb (a PostgreSQL utility command)
# This command is available if you have PostgreSQL command-line tools installed
# Typically used on Unix/macOS systems
createdb rssagg

# Option 2: Using psql (the PostgreSQL interactive terminal)
# This works on all platforms where PostgreSQL is installed
psql -c "CREATE DATABASE rssagg;"

# With explicit user (useful on Windows or when using non-default users)
psql -U postgres -c "CREATE DATABASE rssagg;"
```

**Note about PostgreSQL tools:**
- `createdb` and `psql` are PostgreSQL client utilities that come with PostgreSQL installation
- These are not part of our application but external tools you need to use
- On Windows, you might need to add the PostgreSQL bin directory to your PATH to use these commands
- For more information, see the [PostgreSQL documentation](https://www.postgresql.org/docs/current/app-createdb.html)

**Alternative: Using pgAdmin (Graphical Interface)**
1. Open pgAdmin (install from [pgadmin.org](https://www.pgadmin.org/download/) if needed)
2. Connect to your PostgreSQL server
3. Right-click on "Databases" in the browser tree
4. Select "Create" > "Database"
5. Enter "rssagg" as the database name
6. Click "Save"

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