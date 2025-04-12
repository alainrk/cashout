# HappyPoor

A PostgreSQL-backed Telegram bot application with a clean Go architecture.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Environment Setup](#environment-setup)
- [Database Management](#database-management)
  - [Working with Migrations](#working-with-migrations)
  - [Migration Commands](#migration-commands)
  - [Creating a New Migration](#creating-a-new-migration)
- [Development](#development)
  - [Building and Running](#building-and-running)
  - [Live Reloading](#live-reloading)
  - [Code Quality](#code-quality)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## Overview

HappyPoor is a Telegram bot application built in Go that uses PostgreSQL for data persistence. The application follows a clean architecture pattern with a focus on maintainability and testability.

## Features

- Telegram bot integration
- PostgreSQL database storage
- Robust migration system
- Clean architecture
- Live reloading for development
- Comprehensive testing support

## Project Structure

```
happypoor/
├── cmd/
│   ├── main.go            # Main application entry point
│   └── migrate/           # Migration command tool
│       └── main.go
├── internal/
│   ├── db/                # Database package
│   │   ├── db.go          # Database connection
│   │   └── models.go      # Database models
│   ├── migrations/        # Migration system
│   │   ├── migrations.go  # Migration infrastructure
│   │   └── versions/      # Individual migration files
│   │       ├── 001_create_users_table.go
│   │       └── ...
│   ├── handlers/          # Telegram bot handlers
│   └── services/          # Business logic services
├── scripts/               # Helper scripts
│   └── migration_template.go
├── .gitignore
├── Makefile
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- [golangci-lint](https://golangci-lint.run/) for linting

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/happypoor.git
cd happypoor
```

2. Install dependencies:

```bash
go mod download
```

### Environment Setup

Create a `.env` file in the project root (or set environment variables):

```
DATABASE_URL=postgresql://user:password@localhost:5432/happypoor?sslmode=disable
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
```

## Database Management

HappyPoor uses a version-based migration system to manage database schema changes.

### Working with Migrations

The application uses a custom migration system that:

- Tracks migrations in a database table
- Applies migrations in order
- Allows for rollbacks (when supported)
- Provides status information

### Migration Commands

Use the following commands to manage database migrations:

```bash
# Run all pending migrations
make db/migrate

# Check the status of migrations
make db/status

# Roll back the last migration (requires confirmation)
make db/rollback

# Create a new migration file
make db/migrate/new name="add email column"
```

### Creating a New Migration

1. Create a new migration file:

```bash
make db/migrate/new name="add email column"
```

2. Edit the generated file in `internal/migrations/versions/`:

```go
func add_email_column(tx *gorm.DB) error {
    return tx.Exec(`
        ALTER TABLE users
        ADD COLUMN email VARCHAR(255) UNIQUE
    `).Error
}
```

3. For migrations that support rollback, use the `RegisterMigrationWithRollback` function:

```go
func init() {
    migrations.RegisterMigrationWithRollback("003", "Add email column",
        add_email_column, rollback_email_column)
}

func rollback_email_column(tx *gorm.DB) error {
    return tx.Exec(`ALTER TABLE users DROP COLUMN email`).Error
}
```

## Development

### Building and Running

To build and run the application:

```bash
# Build the application
make build

# Run the application
make run
```

### Live Reloading

For development with live reloading:

```bash
make run/live
```

This uses [Air](https://github.com/cosmtrek/air) to automatically rebuild and restart the application when files change.

### Code Quality

Maintain code quality with:

```bash
# Run linters
make lint

# Format code and tidy dependencies
make tidy

# Attempt to fix lint errors
make lint-fix
```

## Testing

Run tests with:

```bash
# Run all tests
make test

# Run tests with coverage report
make test/cover

# Run tests for CI environments
make test-ci
```

## Deployment

1. Build the binary:

```bash
go build -o happypoor ./cmd/main.go
```

2. Run migrations:

```bash
go run ./cmd/migrate/main.go -command up
```

3. Start the application:

```bash
./happypoor
```

For production deployments, consider using a process manager like systemd or Docker.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
