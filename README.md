# Cashout AI

Telegram **AI Agent** for Income and Expense Management.

You can self-host it following the Developer section down below.

<p align="center">
  <img src="/assets/demo.gif" alt="Demo" height="550px">
</p>

## Features

Cashout is an intelligent Telegram bot that leverages AI to make expense tracking effortless. Simply send a message in natural language, and the bot will understand and categorize your transactions automatically.

### 🤖 AI-Powered Transaction Processing

- **Natural Language Understanding**: Just type "coffee 3.50" or "salary 3000 yesterday" - no complex commands needed
- **Smart Categorization**: Automatically assigns the right category based on your description
- **Flexible Date Recognition**: Understands various date formats (dd/mm, dd-mm-yyyy, "yesterday", etc.)
- **Multi-language Support**: Works with transaction descriptions in any language

### 💰 Transaction Management

- **Quick Entry**: Add expenses and income with a single message
- **Inline Editing**: Modify amount, category, description, or date before confirming
- **Bulk Operations**: Edit or delete existing transactions with paginated navigation
- **Transaction Types**: Track both expenses (18 categories) and income (2 categories)
- **Export Functionality**: Download all your transactions as CSV files

### 📊 Financial Insights

- **Weekly Recap**: Get detailed breakdowns of your current week's spending
- **Monthly Summary**: View month-by-month financial performance with category breakdowns
- **Yearly Overview**: See annual trends and top spending categories
- **Balance Tracking**: Instant calculation of income vs expenses for any period
- **Category Analysis**: Understand where your money goes with percentage breakdowns

### 🔔 Smart Reminders

- **Automated Weekly Recaps**: Receive your previous week's summary every Monday
- **Intelligent Scheduling**: Only sends reminders to active users
- **Reliable Delivery**: Built-in retry mechanism for failed notifications

### 🎯 User Experience

- **Intuitive Interface**: Clean inline keyboards for all operations
- **Smart Navigation**: Year/month selectors for browsing historical data
- **Pagination**: Handle large transaction lists with ease
- **Quick Actions**: Home screen with instant access to all major functions
- **Cancel Anytime**: Every operation can be cancelled mid-flow

### 🛠️ Technical Features

- **Database Migrations**: Version-controlled schema management
- **Webhook & Polling Support**: Flexible deployment options
- **Development Tools**: Built-in database seeder for testing
- **Modular Architecture**: Clean separation of concerns for easy maintenance
- **Configurable Access**: Optional user whitelist for private deployments

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

### Database Seeding

The Dev DB Seeder generates test transaction data for development:

```bash
# Set the user's Telegram ID you want to seed data for
export SEED_USER_TG_ID=123456789

# Seed the database with random transactions
make db/seed
```

The seeder will:

- Generate 5 years of transaction history
- Create 90% expenses and 10% income transactions
- Distribute transactions across all categories
- Ensure at least one salary per month
- Delete existing transactions before seeding (idempotent)

## License

This project is licensed under the MIT License - see the LICENSE file for details.
