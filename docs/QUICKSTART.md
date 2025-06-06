# Quick Start Guide

This guide will help you get started with RSS Aggregator CLI (rssagg) in minutes.

## Installation

### Prerequisites

- Go 1.19+
- PostgreSQL 14+

### Steps

1. **Install the CLI**

   ```bash
   go install github.com/abahnj/rssagg@latest
   ```

   Alternatively, build from source:

   ```bash
   git clone https://github.com/abahnj/rssagg.git
   cd rssagg
   go build -o rssagg
   ```

2. **Set Up Database**

   You'll need to create a PostgreSQL database first. Choose the method that works for your system:

   ```bash
   # OPTION A: Using createdb (PostgreSQL utility command)
   createdb rssagg                          
   
   # OPTION B: Using psql terminal
   psql -c "CREATE DATABASE rssagg;"        
   
   # OPTION C: With explicit user credentials
   psql -U postgres -c "CREATE DATABASE rssagg;"
   ```

   > **Note:** `createdb` and `psql` are PostgreSQL utilities that come with a PostgreSQL installation.
   > They are not part of this application. If these commands aren't found, make sure PostgreSQL
   > is properly installed and its bin directory is in your PATH.

   **Option A: Manual SQL Execution**
   ```bash
   # Run schema migrations manually
   psql -d rssagg -f sql/schema/001_users.sql
   psql -d rssagg -f sql/schema/002_feeds.sql
   psql -d rssagg -f sql/schema/003_feed_follows.sql
   psql -d rssagg -f sql/schema/004_feeds_last_fetched_at.sql
   psql -d rssagg -f sql/schema/005_posts.sql
   ```

   **Option B: Using Goose (Recommended)**
   ```bash
   # Install Goose if you don't have it
   go install github.com/pressly/goose/v3/cmd/goose@latest
   
   # Run all migrations with Goose
   goose -dir sql/schema postgres "postgres://username:password@localhost:5432/rssagg" up
   ```

3. **Configure the Application**

   Create a config file by copying the example:

   ```bash
   # Copy the example config file
   cp config.json.example config.json
   ```

   Then edit `config.json` with your database connection:

   ```json
   {
     "db_url": "postgres://username:password@localhost:5432/rssagg",
     "current_user_name": ""
   }
   ```

   Replace `username` and `password` with your PostgreSQL credentials.

## Usage

### First Run

```bash
# Register a new user
rssagg register yourusername

# You are now logged in as this user
```

### Adding and Following Feeds

```bash
# Add a popular tech feed
rssagg addfeed "Hacker News" https://hnrss.org/newest

# Follow the feed
rssagg follow https://hnrss.org/newest

# Add and follow another feed
rssagg addfeed "Reddit Golang" https://www.reddit.com/r/golang/.rss
rssagg follow https://www.reddit.com/r/golang/.rss

# List your followed feeds
rssagg following
```

### Managing Feeds

```bash
# List all available feeds in the system
rssagg feeds

# Unfollow a feed you no longer want
rssagg unfollow https://www.reddit.com/r/golang/.rss
```

### Browsing Content

```bash
# Browse the latest 10 posts from your followed feeds
rssagg browse

# Browse the latest 20 posts from your followed feeds
rssagg browse 20
```

### Content Aggregation

```bash
# Continuously fetch new content every minute
rssagg agg 1m

# Fetch new content every 30 seconds
rssagg agg 30s
```

## Tips

- The CLI will automatically log you in as the last user you registered or logged in as
- When running `agg`, leave it running in a terminal to continuously update your feeds
- Add multiple feeds for a more diverse content stream
- For continuous operation, consider running as a background service