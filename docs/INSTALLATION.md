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

```bash
# For Linux/macOS
createdb rssagg

# For Windows using psql
psql -U postgres -c "CREATE DATABASE rssagg"
```

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

Note: The `.gatorconfig.json` format is also supported (copy from `.gatorconfig.json.example`), but `config.json` in the project root is the recommended approach.

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