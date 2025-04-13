# HappyPoor

Telegram bot for expense tracking AI-Powered (for real!).

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL Database
- Access to an OpenAI-compatible API model, with its API Key and Endpoint (e.g. DeepSeek, OpenAI, etc.)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/alainrk/happypoor.git
cd happypoor
```

2. Install dependencies:

```bash
go mod download
```

### Environment Setup

Copy the example `.env` file in the project root (or set environment variables) and edit it accordingly:

```bash
cp .env.example .env
```

## Database Management

HappyPoor uses a version-based migration system to manage database schema changes.

Use the following commands to manage database migrations:

```bash
# Run all pending migrations
go run ./cmd/migrate/main.go -command up

# Create a new migration just by copy-pasting a previous one and editing it accordingly
cp internal/migrations/versions/001*.go internal/migrations/versions/00X_your_migration.go
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

# Run the application with live reloading (requires Air)
make run/live
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
