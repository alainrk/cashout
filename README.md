# Cashout AI

Telegram **AI Agent** for Income and Expense Management with a comprehensive **Web Dashboard** featuring **Passkey/WebAuthn** support for secure, passwordless authentication.

You can self-host it following the Developer section down below.

<p align="center">
  <img src="/assets/demo-tg.gif" alt="Demo" height="550px">
  <img src="/assets/demo-web.gif" alt="Demo" height="550px">
</p>

## Features

Cashout is an intelligent Telegram bot that leverages AI to make expense tracking effortless. Simply send a message in natural language, and the bot will understand and categorize your transactions automatically.

### 🤖 AI-Powered Transaction Processing

- **Intelligent Intent Routing**: Just type naturally - no commands needed. The AI understands whether you want to add a transaction, check your weekly summary, search, edit, delete, or export. Simply say "show me this week" or "delete coffee" and the bot figures out the rest
- **Natural Language Understanding**: Just type "coffee 3.50" or "salary 3000 yesterday" - no complex commands needed
- **Smart Categorization**: Automatically assigns the right category based on your description
- **Flexible Date Recognition**: Understands various date formats (dd/mm, dd-mm-yyyy, "yesterday", etc.)
- **Multi-language Support**: Works with transaction descriptions in any language

### 💰 Transaction Management

- **Quick Entry**: Add expenses and income with a single message
- **Inline Editing**: Modify amount, category, description, or date before confirming
- **Bulk Operations**: Edit or delete existing transactions with paginated navigation
- **Transaction Types**: Track both expenses (18 categories) and income (2 categories)
- **Search and Full Listing**: Find transactions by full text search and category or full listing
- **Export Functionality**: Download all your transactions as CSV files

### 📊 Financial Insights

- **Weekly Recap**: Get detailed breakdowns of your current week's spending
- **Monthly Summary**: View month-by-month financial performance with category breakdowns
- **Yearly Overview**: See annual trends and top spending categories
- **Balance Tracking**: Instant calculation of income vs expenses for any period
- **Category Analysis**: Understand where your money goes with percentage breakdowns

### 🌐 Web Dashboard

- **Multiple Authentication Methods**:
  - Telegram-based login with verification codes
  - Email-based passwordless authentication
  - Passkey/WebAuthn support for passwordless biometric login
- **Transaction Management**:
  - Add new transactions directly from the web interface
  - View detailed transaction history with search and filtering
  - Monthly navigation with intuitive controls
- **Visual Analytics**:
  - Real-time balance, income, and expense statistics
  - Category breakdowns and trends
  - Transaction counts and summaries
- **Security Features**:
  - Rate-limited authentication endpoints
  - Secure session management with configurable duration
  - Support for multiple passkeys per user
  - Passkey management (register, list, delete)

### 🔔 Smart Reminders

- **Automated Weekly Recaps**: Receive your previous week's summary every Monday
- **Automated Monthly Recaps**: Receive your previous month's summary on the 1st of each month
- **Intelligent Scheduling**: Only sends reminders to active users
- **Reliable Delivery**: Built-in retry mechanism for failed notifications

### 💻 Available Commands

- `/start` - Initialize the bot and see the main menu
- `/edit` - Edit an existing transaction
- `/delete` - Delete a transaction
- `/list` - View all transactions (paginated)
- `/search` - Search transactions by description
- `/week` - Get current week's financial summary
- `/month` - Get current month's financial summary
- `/year` - Get current year's financial summary
- `/export` - Export all transactions to CSV

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

- Go 1.24 or higher
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

```env
TELEGRAM_BOT_API_TOKEN='XXXXXXXXXX:AAAA_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'
DATABASE_URL='postgres://postgres:postgres@localhost:5432/postgres'
OPENAI_API_KEY='sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
OPENAI_BASE_URL='https://api.deepseek.com/v1'
LLM_MODEL='deepseek-chat'
RUN_MODE='webhook' # webhook or polling
WEBHOOK_DOMAIN=''
WEBHOOK_SECRET=''
WEBHOOK_HOST='localhost'
WEBHOOK_PORT='8080'
LOG_LEVEL='info'
# Dev purpose, comma separated. Keep it empty to allow all
ALLOWED_USERS=''
# Seed purpose - set the Telegram ID of the user to seed transactions for
SEED_USER_TG_ID=''
# Web Server Configuration
WEB_HOST=localhost
WEB_PORT=8081
# Session Configuration (optional)
SESSION_SECRET=your-random-session-secret-here
SESSION_DURATION=24h
# Email Service Configuration (for passwordless email login)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=your-email@gmail.com
# WebAuthn Configuration (for passkey support)
WEBAUTHN_RP_NAME=Cashout
WEBAUTHN_RP_ID=localhost
WEBAUTHN_ORIGIN=http://localhost:8081
```

Spin up local infrastructure:

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

For migrations that support rollback, use the `RegisterMigrationWithRollback` function:

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

#### Telegram Bot

```bash
# Build the bot
make build

# Run the bot
make run

# Run the bot with live reloading (requires Air)
make run/live
```

#### Web Server

```bash
# Build the web server
make build-web

# Run the web server
make run-web

# Run the web server with live reloading
make run/live-web
```

#### Both Services

```bash
# Build both applications
make build-all

# Build both for Linux
make build-linux-all

# Note: Running both requires two terminals
# Terminal 1: make run
# Terminal 2: make run-web
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

## Deployment

### Docker Compose (Recommended)

The project includes a complete Docker Compose setup:

```bash
# Start all services (database, bot, web server)
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down
```

This will start:

- PostgreSQL database on port 5432
- Telegram bot (webhook mode on port 8080)
- Web dashboard on port 8081
- Automatic database migrations

### Manual Deployment

#### Telegram Bot

The bot can run both in `webhook` and `polling` mode.

**Webhook Mode:**

```env
RUN_MODE='webhook'
WEBHOOK_DOMAIN='https://your-domain.com'
WEBHOOK_SECRET='xxxyyyzzz'
WEBHOOK_PORT='8080'
```

**Polling Mode:**

```env
RUN_MODE='polling'
```

#### Web Server

The web server runs independently and can be configured:

```env
WEB_HOST=0.0.0.0  # For production
WEB_PORT=8081
SESSION_SECRET=your-random-session-secret-here
SESSION_DURATION=24h
```

### LLM Setup

Any OpenAI compatible API LLM can be used:

**Example with DeepSeek:**

```env
OPENAI_API_KEY='sk-xxx'
OPENAI_BASE_URL='https://api.deepseek.com/v1'
LLM_MODEL='deepseek-chat'
```

**Example with OpenAI:**

```env
OPENAI_API_KEY='sk-xxx'
OPENAI_BASE_URL='https://api.openai.com/v1'
LLM_MODEL='gpt-4'
```

## Web Dashboard Usage

### Authentication Options

The web dashboard supports three authentication methods:

1. **Telegram Login** (code-based):
   - Enter your Telegram username
   - Check Telegram for a 6-digit verification code
   - Enter the code to access your dashboard

2. **Email Login** (passwordless):
   - Enter your registered email address
   - Check your email for a 6-digit verification code
   - Enter the code to access your dashboard

3. **Passkey Login** (WebAuthn):
   - Register a passkey from your dashboard settings after initial login
   - Use biometric authentication (fingerprint, face recognition) on subsequent logins
   - No codes needed - instant secure access

### Dashboard Features

1. **Access**: Navigate to `http://localhost:8081` (or your configured domain)
2. **Login**: Choose your preferred authentication method
3. **Dashboard**: View your financial data with month navigation
4. **Statistics**: See real-time balance, income, expenses, and transaction counts
5. **Transactions**:
   - Add new transactions directly from the web interface
   - Browse detailed transaction history with search and filtering
   - View transactions by category
6. **Passkey Management**: Register, view, and delete passkeys for your account

The web dashboard provides a complementary interface to the Telegram bot, offering:

- Better visualization for large datasets
- Month-by-month navigation
- Desktop and mobile-friendly transaction management
- Multiple secure authentication options
- Direct transaction creation without needing Telegram

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests in CI mode
make test-ci

# Run security checks
make sec
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
