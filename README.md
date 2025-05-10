# Cashout

Telegram bot for expense tracking **AI-Powered** (_for real!_).

<p align="center">
  <img src="/assets/cashout.png" alt="Logo" width="250px">
</p>

## Demo

<p align="center">
  <img src="/assets/demo.gif" alt="Demo" height="550px">
</p>

## Features

<details open>

<p align="center">
  <img src="/assets/7.png" alt="Cashout app welcome screen" width="450px">
</p>

Navigate your transactions in a simple way.

<p align="center">
  <img src="/assets/8.png" alt="List command" width="450px">
</p>

<p align="center">
  <img src="/assets/9.png" alt="List output" width="450px">
</p>

Expense and income creation with options to edit or confirm. It tries to automatically categorize your transactions and fix some common errors. Also, it tries to extract the correct amount despite given in a natural language or anyway, not in a standard format.

<p align="center">
  <img src="/assets/3.png" alt="Transaction confirmation interface" width="450px">
</p>

<p align="center">
  <img src="/assets/6.png" alt="Income entry interface" width="450px">
</p>

Change the date or category of a transaction with any intelligible format.

<p align="center">
  <img src="/assets/4.png" alt="Date entry interface" width="450px">
</p>

<p align="center">
  <img src="/assets/2.png" alt="Category selection interface" width="450px">
</p>

Select and delete transactions from your records.

<p align="center">
  <img src="/assets/1.png" alt="Delete Transaction interface" width="450px">
</p>

Monthly and yearly financial summary.

<p align="center">
  <img src="/assets/5.png" alt="Financial summary displays" width="450px">
</p>

</details>

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL Database
- Access to an OpenAI-compatible API model, with its API Key and Endpoint (e.g. DeepSeek, OpenAI, etc.)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/alainrk/cashout.git
cd cashout
```

2. Install dependencies:

```bash
go mod download
```

### Environment Setup

```bash
cp .env.example .env
```

Copy the example `.env` file in the project root (or set environment variables) and edit it accordingly:
Spin up local infrastructure

```bash
docker compose up -d
```

## Database Management

Cashout uses a version-based migration system to manage database schema changes.

Use the following commands to manage database migrations:

```bash
# Run all pending migrations
go run ./cmd/migrate/main.go -command up

# Run all pending migrations with another .env file
go run ./cmd/migrate/main.go -command up -env .prod.env

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

## Run Mode

The bot can run both in `webhook` and `polling` mode.

### Webhook Mode

You need to set the relevant environment variables:

```
RUN_MODE='webhook'
WEBHOOK_DOMAIN='https://your-domain.com'
WEBHOOK_SECRET='xxxyyyzzz'
WEBHOOK_PORT='8080'
```

### Polling Mode

You need to set the relevant environment variable:

```
RUN_MODE='polling'
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
