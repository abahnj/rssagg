# Installation Guide

This document provides comprehensive installation instructions for the RSS Aggregator CLI (rssagg).

## Prerequisites

### Go

The application requires Go 1.19 or later.

#### Installing Go

- **Linux/macOS**:
  ```bash
  # Using Homebrew (macOS)
  brew install go

  # Using apt (Ubuntu/Debian)
  sudo apt update
  sudo apt install golang-go
  ```

- **Windows**:
  Download and run the installer from [go.dev/dl](https://go.dev/dl/)

Verify installation:
```bash
go version
```

### PostgreSQL

The application requires PostgreSQL 14 or later.

#### Installing PostgreSQL

- **Linux**:
  ```bash
  # Ubuntu/Debian
  sudo apt update
  sudo apt install postgresql postgresql-contrib

  # Start the service
  sudo systemctl start postgresql
  sudo systemctl enable postgresql
  ```

- **macOS**:
  ```bash
  # Using Homebrew
  brew install postgresql

  # Start the service
  brew services start postgresql
  ```

- **Windows**:
  Download and run the installer from [postgresql.org/download/windows](https://www.postgresql.org/download/windows/)

#### Creating a Database User

```bash
# Access PostgreSQL as the postgres user
sudo -u postgres psql

# Create a new user (replace username and password)
CREATE USER yourusername WITH PASSWORD 'yourpassword';

# Grant necessary privileges
ALTER USER yourusername CREATEDB;

# Exit psql
\q
```

## Installing RSS Aggregator CLI

### Method 1: Using Go Install

This is the easiest way to install the application:

```bash
go install github.com/abahnj/rssagg@latest
```

This will download, compile, and install the binary to your `$GOPATH/bin` directory. Ensure this directory is in your `PATH` environment variable.

### Method 2: Building from Source

For the latest development version or to make local modifications:

```bash
# Clone the repository
git clone https://github.com/abahnj/rssagg.git
cd rssagg

# Build the application
go build -o rssagg

# Optionally, move the binary to a directory in your PATH
sudo mv rssagg /usr/local/bin/  # Linux/macOS
# or
move rssagg %GOPATH%\bin\  # Windows
```

## Database Setup

### Creating the Database

You can create the database using one of these methods:

```bash
# Method 1: Using createdb (Unix/macOS)
createdb rssagg

# Method 2: Using psql (all platforms)
# Without credentials (if using peer authentication)
psql -c "CREATE DATABASE rssagg;"

# With credentials
psql -U postgres -c "CREATE DATABASE rssagg;"

# Method 3: Using pgAdmin or another GUI tool
# 1. Open pgAdmin
# 2. Connect to your PostgreSQL server
# 3. Right-click on "Databases" and select "Create" > "Database"
# 4. Enter "rssagg" as the database name and save
```

### Setting Up the Schema

You have two options for setting up the database schema:

#### Option 1: Manual SQL Execution

You can manually run the SQL migration files to set up the database schema:

```bash
# For Linux/macOS
psql -d rssagg -f sql/schema/001_users.sql
psql -d rssagg -f sql/schema/002_feeds.sql
psql -d rssagg -f sql/schema/003_feed_follows.sql
psql -d rssagg -f sql/schema/004_feeds_last_fetched_at.sql
psql -d rssagg -f sql/schema/005_posts.sql

# For Windows
psql -U postgres -d rssagg -f sql/schema/001_users.sql
psql -U postgres -d rssagg -f sql/schema/002_feeds.sql
psql -U postgres -d rssagg -f sql/schema/003_feed_follows.sql
psql -U postgres -d rssagg -f sql/schema/004_feeds_last_fetched_at.sql
psql -U postgres -d rssagg -f sql/schema/005_posts.sql
```

Run these commands in order, as each migration builds upon the previous one.

#### Option 2: Using Goose for Migrations (Recommended)

[Goose](https://github.com/pressly/goose) is a database migration tool that helps manage schema changes more effectively.

1. Install Goose:

```bash
# Using Go install
go install github.com/pressly/goose/v3/cmd/goose@latest

# Or on macOS using Homebrew
brew install goose
```

2. Run all migrations:

```bash
# Navigate to the project directory
cd /path/to/rssagg

# Run migrations
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/rssagg" up
```

Replace `username` and `password` with your PostgreSQL credentials.

Goose has several advantages:
- Automatically tracks migration versions
- Allows rolling back migrations if needed
- Ensures migrations run in the correct order

## Configuration

Create a configuration file named `config.json` in the project root directory. You can copy from the example file:

```bash
# Copy the example config file
cp config.json.example config.json
```

Edit the configuration file to include your database connection:

```json
{
  "db_url": "postgres://yourusername:yourpassword@localhost:5432/rssagg",
  "current_user_name": ""
}
```


Replace `yourusername` and `yourpassword` with your PostgreSQL credentials.

### Connection URL Format

The database connection URL follows this format:
```
postgres://username:password@host:port/database_name
```

- `username`: Your PostgreSQL username
- `password`: Your PostgreSQL password
- `host`: The host where PostgreSQL is running (usually `localhost`)
- `port`: The port PostgreSQL is listening on (default is `5432`)
- `database_name`: The name of the database (in this case, `rssagg`)

## Verifying Installation

```bash
# Check that the application runs
rssagg

# You should see a list of available commands
```

## Running as a Service (Optional)

### Linux (systemd)

Create a systemd service file:

```bash
sudo touch /etc/systemd/system/rssagg.service
sudo chmod 644 /etc/systemd/system/rssagg.service
```

Edit the file with the following content:

```ini
[Unit]
Description=RSS Aggregator Service
After=network.target postgresql.service

[Service]
Type=simple
User=yourusername
ExecStart=/usr/local/bin/rssagg agg 5m
Restart=on-failure
Environment=HOME=/home/yourusername

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable rssagg
sudo systemctl start rssagg
```

## Troubleshooting

### Common Issues

1. **Database connection errors**:
   - Ensure PostgreSQL is running: `systemctl status postgresql`
   - Verify your credentials in the connection URL
   - Check that the database exists: `psql -l`

2. **Permission denied errors**:
   - Ensure the user has appropriate permissions in PostgreSQL
   - Check file permissions for configuration files

3. **Command not found**:
   - Ensure the binary is in your PATH
   - Try running with the full path to the binary