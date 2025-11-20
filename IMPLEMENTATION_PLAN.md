# Web App Implementation Plan

## Overview

Transform Cashout from a Telegram-only app to a dual-interface system with both Telegram bot and standalone web app. This requires creating a proper REST API layer and a pure HTML/CSS/JS web client.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Client Layer                                │
├──────────────────────────────┬──────────────────────────────────────┤
│   Telegram Bot               │   Web App (HTML/CSS/JS)              │
│   (existing)                 │   (new - zero frameworks)            │
│   - Chat interface           │   - Login/Register                   │
│   - Bot commands             │   - Dashboard                        │
│   - Inline keyboards         │   - Transaction CRUD                 │
│   - Auth via Telegram        │   - Statistics & Reports             │
│                              │   - Account Linking UI               │
└──────────────────────────────┴──────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Server Layer (Go)                              │
├──────────────────────────────┬──────────────────────────────────────┤
│   internal/client/           │   internal/api/                      │
│   (existing Telegram logic)  │   (new REST API)                     │
│   - Bot handlers             │   - JWT authentication               │
│   - State management         │   - RESTful endpoints                │
│   - Telegram-specific UI     │   - JSON responses                   │
│                              │   - CORS & middleware                │
└──────────────────────────────┴──────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Business Logic Layer                             │
│                    internal/repository/                             │
│                    (shared by both interfaces)                      │
│   - User management (with reconciliation)                           │
│   - Transaction CRUD                                                │
│   - Statistics calculations                                         │
│   - Data access abstraction                                         │
└─────────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Data Layer                                  │
│                       PostgreSQL Database                           │
│   - users (with user_id, tg_id, email)                             │
│   - transactions (linked via user_id)                               │
│   - web_sessions                                                    │
└─────────────────────────────────────────────────────────────────────┘

User Reconciliation:
  Same user can authenticate via:
    - Telegram (tg_id) → internal/client/ → repository → database
    - Web (email/password) → JWT → internal/api/ → repository → database
  Both paths access the same user record via user_id
```

## Current Architecture Analysis

### Existing Components

- **Telegram Bot** (`internal/client/`): Handles all user interactions via Telegram
- **Web Dashboard** (`internal/web/`): Basic read-only dashboard with authentication
- **Repository Layer** (`internal/repository/`): Data access layer (already exists)
- **Model Layer** (`internal/model/`): Data structures (already exists)
- **AI Integration** (`internal/ai/`): Transaction parsing via LLM

### Current Web Features

- Session-based authentication (Telegram-verified)
- Read-only dashboard showing:
  - Monthly statistics (balance, income, expenses, transaction count)
  - Transaction list (list and clustered views)
  - Month navigation
- API endpoints: `/web/api/stats`, `/web/api/transactions` (read-only)

### Missing for Full Web App

- Transaction creation/editing/deletion via API
- User registration/authentication independent of Telegram
- Full CRUD operations API
- Search functionality
- Export functionality
- Reporting (weekly/monthly/yearly recaps)

---

## User Reconciliation Strategy

### Overview

The system will support users accessing their data from both Telegram and Web interfaces. User reconciliation (linking Telegram and email accounts) will be done **manually by administrators**, not through automated account linking features.

Two types of users exist:

1. **Telegram-only users**: Existing users who only use the bot (have `tg_id` only)
2. **Web-only users**: New users who register via web app with **passwordless authentication** (have `email` only)

**Note**: To allow a user to access their data from both interfaces, an administrator will manually update the database to add email to a Telegram user's record (or vice versa)

**Authentication Method**: Web users authenticate via **magic links** sent to their email. Links expire after 30 days.

### User Types Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Users Table                                     │
├──────────────┬──────────┬──────────┬──────────────────┬─────────────────┤
│ user_id      │ tg_id    │ email    │ magic_link_token │ magic_link_exp  │
│ (UUID)       │ (int64)  │ (string) │ (string)         │ (timestamp)     │
├──────────────┼──────────┼──────────┼──────────────────┼─────────────────┤
│                                                                          │
│ Type 1: Telegram-only User                                              │
│ abc-123-def  │ 98765432 │ NULL     │ NULL             │ NULL            │
│              │          │          │                  │                 │
│ Type 2: Web-only User (passwordless)                                    │
│ xyz-456-ghi  │ NULL     │ user@ex  │ token123...      │ 2025-12-20      │
│              │          │          │                  │                 │
│ Type 3: Manually Linked (both Telegram + Web)                          │
│ qwe-789-rty  │ 12345678 │ linked@  │ token456...      │ 2025-12-20      │
│ (Admin manually added email to Telegram user's record)                 │
└──────────────┴──────────┴──────────┴──────────────────┴─────────────────┘

Authentication Paths:
  Telegram Bot → tg_id lookup → User
  Web Login    → email lookup → User
  API Request  → JWT (user_id) → User
```

### User Identification

- **Primary Key**: `tg_id` (int64) - Remains for backwards compatibility
- **Universal ID**: `user_id` (UUID) - New field for API-first identification
- **Email**: `email` (string) - Optional, unique when present
- **Telegram ID**: `tg_id` (int64) - Optional for new web users, required for Telegram users

### Authentication Flows

#### Telegram User Authenticating

1. User interacts with bot
2. Bot uses `tg_id` to identify user
3. User data is fetched using `tg_id`

#### Web User Authenticating (Email/Password)

1. User logs in with email/password
2. API validates credentials
3. JWT token issued with `user_id` claim
4. User data is fetched using `user_id` or `email`

#### Manually Linked User Scenario

1. Administrator identifies a user who wants both Telegram and Web access
2. Administrator manually updates the database to add:
   - Email and password_hash to an existing Telegram user, OR
   - tg_id to an existing Web user
3. Once linked, user can access the same data from both interfaces
4. All transactions, settings, and data are shared

### Manual Reconciliation Process

**For Telegram users wanting Web access:**
```sql
-- Admin manually adds email to enable passwordless web login
UPDATE users
SET email = 'user@example.com'
WHERE tg_id = 98765432;

-- User will then request a magic link via the web app
-- Magic link token and expiry will be set automatically when requested
```

**For Web users wanting Telegram access:**
```sql
-- Admin manually adds Telegram ID
UPDATE users
SET tg_id = 98765432
WHERE email = 'user@example.com';
```

### Data Access Patterns

#### Repository Layer Changes

```go
// GetUserByTgID - Existing method (unchanged)
func (r *UserRepository) GetUserByTgID(tgID int64) (*User, error)

// GetUserByEmail - New method
func (r *UserRepository) GetUserByEmail(email string) (*User, error)

// GetUserByUserID - New method
func (r *UserRepository) GetUserByUserID(userID string) (*User, error)

// GetUserByIdentifier - New unified method
func (r *UserRepository) GetUserByIdentifier(identifier string) (*User, error)
// Tries to match by user_id, email, or tg_id
```

#### Transaction Access

- Transactions are currently linked to users via `tg_id` (foreign key)
- **Option A - Keep tg_id as FK**:
  - For web-only users, generate a unique synthetic `tg_id` (e.g., using snowflake algorithm or random large negative number)
  - Pro: No schema changes to transactions table
  - Con: `tg_id` becomes less meaningful
- **Option B - Migrate to user_id as FK** (RECOMMENDED):
  - Add migration to change transactions foreign key from `tg_id` to `user_id`
  - Update existing transactions to use `user_id` instead of `tg_id`
  - Pro: Cleaner schema, `user_id` is universal identifier
  - Con: Requires migration of existing data
- **Decision**: Use Option B for cleaner architecture

### Constraints & Validation

- A user must have at least one authentication method (`tg_id` OR `email` + `password_hash`)
- Cannot unlink the last remaining authentication method
- `email` must be unique across all users
- `tg_id` must be unique across all users
- `user_id` must be unique across all users

### Migration Path for Existing Users

1. Add new fields to users table via migration
2. Generate `user_id` (UUID) for all existing users
3. Existing Telegram users have `email`, `password_hash` as NULL
4. They can optionally link email/password later

---

## Implementation Plan

### Stage 1: API Foundation & Authentication

**Goal**: Create a complete REST API layer with proper authentication that works independently of Telegram, with user reconciliation allowing both Telegram ID and email as user identifiers

**Components to Build**:

1. **API Package** (`internal/api/`)
   - `server.go` - API server setup with middleware
   - `router.go` - API route definitions
   - `middleware.go` - Common middleware (auth, CORS, rate limiting, logging)
   - `response.go` - Standard JSON response helpers

2. **User Reconciliation System**
   - **Unified User Model**: Single user can be accessed via both `tg_id` and `email`
   - **Manual Reconciliation**: Admins manually link accounts via database updates
   - **Dual Authentication Modes**:
     - Telegram users: Authenticated via Telegram widget/bot verification
     - Email users: Authenticated via passwordless magic link + JWT
   - **User Lookup**: Repository methods to find user by `tg_id` OR `email`

3. **Authentication System** (extend `internal/repository/auth.go`)
   - JWT-based authentication for API
   - **Passwordless email authentication** (magic link)
   - User registration (email only - no password)
   - Login via email magic link
   - Token refresh mechanism
   - Magic link expiration (30 days)

4. **API Endpoints - Authentication**
   - `POST /api/v1/auth/register` - Create new user account (email only)
   - `POST /api/v1/auth/login` - Request magic link via email
   - `GET /api/v1/auth/verify?token=<token>` - Verify magic link and issue JWT
   - `POST /api/v1/auth/refresh` - Refresh access token
   - `POST /api/v1/auth/logout` - Invalidate token
   - `GET /api/v1/auth/me` - Get current user info

**Database Changes**:

_Migration 1 - User Authentication Fields_:

- Add to `users` table:
  - `user_id` (UUID, unique, indexed, not null) - Internal universal identifier
  - `email` (VARCHAR, nullable, unique) - For passwordless email authentication
  - `magic_link_token` (VARCHAR, nullable) - Token for magic link authentication
  - `magic_link_expires` (TIMESTAMP, nullable) - Magic link expiration (30 days from issue)
  - `last_login` (TIMESTAMP, nullable) - Track last successful login
- Generate `user_id` (UUID) for all existing users
- Add check constraint: `tg_id IS NOT NULL OR email IS NOT NULL`

_Migration 2 - Transaction Foreign Key Update_:

- Add `user_id` (UUID, indexed) column to `transactions` table
- Populate `transactions.user_id` by joining with `users.tg_id`
- Add foreign key constraint: `transactions.user_id` → `users.user_id`
- Keep `tg_id` column in transactions for backwards compatibility (can remove later)
- Update indexes to use `user_id` instead of `tg_id`

**Note on User Identification**:

- `tg_id` remains the primary key for backwards compatibility
- `user_id` (UUID) is a new universal identifier for API usage
- JWT tokens will contain `user_id`
- Repository layer will support lookups by `tg_id`, `email`, or `user_id`
- A user can have:
  - Only `tg_id` (Telegram-only user)
  - Only `email` (Web-only user with passwordless auth)
  - Both `tg_id` AND `email` (Linked account - manually reconciled)

**Tests**:

- [ ] User registration with email only (no password)
- [ ] Request magic link via email
- [ ] Verify magic link and receive JWT
- [ ] Magic link expiration after 30 days
- [ ] Rate limiting on magic link requests
- [ ] Token refresh flow
- [ ] JWT validation middleware
- [ ] User lookup by tg_id
- [ ] User lookup by email
- [ ] User lookup by user_id (UUID)
- [ ] Transaction access via both Telegram and web with same user (after manual linking)

**Success Criteria**:

- Can register a new user via API (email only, no password)
- Can request magic link and receive it via email
- Magic link successfully authenticates user and issues JWT
- Magic links expire after 30 days
- Token authentication works on protected endpoints
- Both Telegram users and email users can coexist
- Manually linked users can access their data from both Telegram bot and web app
- Repository methods support lookups by tg_id, email, and user_id

**Status**: Not Started

---

### Stage 2: Transaction CRUD API

**Goal**: Implement full CRUD operations for transactions via REST API

**Components to Build**:

1. **Transaction API Handlers** (`internal/api/transactions.go`)
   - Create transaction
   - Update transaction
   - Delete transaction
   - Get transaction by ID
   - List transactions with pagination
   - Search transactions (full-text + filters)
   - Bulk import transactions

2. **API Endpoints - Transactions**
   - `POST /api/v1/transactions` - Create transaction
   - `GET /api/v1/transactions/:id` - Get single transaction
   - `PUT /api/v1/transactions/:id` - Update transaction
   - `DELETE /api/v1/transactions/:id` - Delete transaction
   - `GET /api/v1/transactions` - List transactions (paginated, filterable)
     - Query params: `page`, `limit`, `type`, `category`, `start_date`, `end_date`, `search`
   - `POST /api/v1/transactions/search` - Advanced search with full-text
   - `POST /api/v1/transactions/import` - Bulk import (CSV)

3. **Request/Response Models** (`internal/api/models.go`)
   - `CreateTransactionRequest`
   - `UpdateTransactionRequest`
   - `TransactionResponse`
   - `TransactionListResponse` (with pagination metadata)
   - `SearchTransactionRequest`

4. **Validation Layer**
   - Input validation for all transaction fields
   - Date validation
   - Amount validation (positive, reasonable precision)
   - Category validation (enum)
   - Currency validation (enum)

**Reuse Existing**:

- LLM transaction parsing from `internal/ai/` for natural language input
- Repository methods from `internal/repository/transactions.go`

**Tests**:

- [ ] Create transaction with valid data
- [ ] Create transaction with AI parsing
- [ ] Update transaction fields
- [ ] Delete transaction
- [ ] List transactions with pagination
- [ ] Filter transactions by date range
- [ ] Filter by category and type
- [ ] Search transactions by description

**Success Criteria**:

- Can perform all CRUD operations via API
- AI-powered transaction creation works
- Proper error handling and validation
- Pagination works correctly

**Status**: Not Started

---

### Stage 3: Statistics & Reporting API

**Goal**: Expose all reporting and statistics functionality via API

**Components to Build**:

1. **Statistics API Handlers** (`internal/api/statistics.go`)
   - Monthly summary
   - Weekly summary
   - Yearly summary
   - Category breakdown
   - Spending trends

2. **API Endpoints - Statistics**
   - `GET /api/v1/stats/monthly?month=YYYY-MM` - Monthly statistics
   - `GET /api/v1/stats/weekly?week=YYYY-Www` - Weekly statistics
   - `GET /api/v1/stats/yearly?year=YYYY` - Yearly statistics
   - `GET /api/v1/stats/categories?start_date=&end_date=` - Category breakdown
   - `GET /api/v1/stats/trends?period=month&count=6` - Spending trends

3. **Export API Handlers** (`internal/api/export.go`)
   - CSV export
   - PDF export (optional for later)

4. **API Endpoints - Export**
   - `GET /api/v1/export/csv?start_date=&end_date=` - Export to CSV
   - `POST /api/v1/export/csv` - Export with advanced filters

**Reuse Existing**:

- Calculation logic from `internal/client/week.go`, `month.go`, `year.go`
- Export logic from `internal/client/export.go`

**Tests**:

- [ ] Monthly stats calculation
- [ ] Weekly stats calculation
- [ ] Yearly stats calculation
- [ ] Category breakdown accuracy
- [ ] CSV export format

**Success Criteria**:

- All statistics match existing Telegram bot calculations
- CSV export works with date ranges
- Proper data aggregation and formatting

**Status**: Not Started

---

### Stage 4: Pure HTML/CSS/JS Web Client

**Goal**: Build a zero-framework web client using only vanilla HTML, CSS, and JavaScript

**File Structure**:

```
web/
├── index.html           # Landing/login page
├── register.html        # Registration page
├── verify.html          # Magic link verification page
├── dashboard.html       # Main dashboard (replaces current)
├── transactions.html    # Transaction management page
├── stats.html           # Statistics & reports page
├── settings.html        # User settings
├── css/
│   ├── reset.css       # CSS reset
│   ├── variables.css   # CSS custom properties
│   ├── layout.css      # Layout components
│   └── components.css  # Reusable components
└── js/
    ├── api.js          # API client wrapper
    ├── auth.js         # Authentication logic
    ├── storage.js      # LocalStorage utilities
    ├── components.js   # Reusable UI components
    ├── dashboard.js    # Dashboard page logic
    ├── transactions.js # Transaction page logic
    ├── stats.js        # Statistics page logic
    └── settings.js     # Settings page logic (including account linking)
```

**Components to Build**:

1. **Authentication Pages**
   - `index.html` - Login page
     - Email input field only
     - "Send Magic Link" button
     - Message: "Check your email for login link"
   - `register.html` - Registration page
     - Email and name fields only (no password)
     - "Register" button
     - Notice: "You'll receive a login link via email"
   - `verify.html` - Magic link verification page
     - Auto-verifies token from URL
     - Shows success/error message
     - Auto-redirects to dashboard on success

2. **Dashboard Page** (`dashboard.html`)
   - Month navigation
   - Statistics cards (balance, income, expenses, count)
   - Quick transaction entry with AI parsing
   - Recent transactions list
   - Quick filters (income/expense, category)

- **Settings Page** (`settings.html`)
  - User profile information
  - Display linked accounts (Telegram and/or Email) - read-only
  - Request new magic link (if current one expired)
  - Account deletion

3. **Transactions Page** (`transactions.html`)
   - Advanced search form
     - Text search with AI parsing
     - Date range picker
     - Category filter (multi-select)
     - Type filter (income/expense)
   - Transaction list with pagination
   - Inline editing
   - Delete with confirmation
   - Bulk actions (delete, export)
   - Sort by date, amount, category

4. **Statistics Page** (`stats.html`)
   - Period selector (weekly/monthly/yearly)
   - Date range picker
   - Charts (using Canvas API or simple bar charts)
     - Income vs Expenses over time
     - Category breakdown pie chart
     - Spending trends
   - Export button for CSV

5. **Shared UI Components**
   - Navigation bar
   - Transaction form (create/edit)
   - Modal dialog
   - Loading spinner
   - Toast notifications
   - Confirmation dialog
   - Date picker (native or simple custom)
   - Dropdown/select
   - Pagination controls

6. **JavaScript Modules**
   - `api.js`:
     - Axios-like fetch wrapper
     - JWT token management
     - Request/response interceptors
     - Error handling
   - `auth.js`:
     - Login/logout logic
     - Token storage
     - Auto-redirect on 401
     - Token refresh
   - `storage.js`:
     - LocalStorage wrapper
     - Session management
     - Cache helpers
   - `components.js`:
     - Modal system
     - Toast notifications
     - Reusable form components
   - `utils.js`:
     - Date formatting
     - Currency formatting
     - Validation helpers

**CSS Approach**:

- Mobile-first responsive design
- CSS Grid and Flexbox for layout
- CSS custom properties for theming
- Minimal, clean design (inspired by current dashboard)
- No CSS framework (pure CSS)
- Support for dark mode (media query)

**Features**:

- [ ] Passwordless login flow (magic link)
- [ ] User registration (email only)
- [ ] Magic link verification
- [ ] Dashboard with monthly stats
- [ ] Quick transaction entry with AI
- [ ] Transaction list with pagination
- [ ] Transaction search and filters
- [ ] Transaction CRUD operations
- [ ] Category-based filtering
- [ ] Date range filtering
- [ ] CSV export
- [ ] Statistics visualization
- [ ] Mobile-responsive design
- [ ] Offline state handling
- [ ] Loading states
- [ ] Error handling

**Success Criteria**:

- Works on all modern browsers (Chrome, Firefox, Safari, Edge)
- Fully responsive (mobile, tablet, desktop)
- No framework dependencies
- Fast load time (<2s on 3G)
- Accessible (keyboard navigation, screen readers)
- Works offline for viewing cached data

**Status**: Not Started

---

### Stage 5: Integration & Polish

**Goal**: Integrate all components, ensure backwards compatibility, and add finishing touches

**Tasks**:

1. **Server Configuration**
   - Update `cmd/web/main.go` to serve both API and static files
   - Configure CORS properly
   - Set up different ports or path prefixes for API vs old web dashboard
   - Environment variables for API configuration

2. **Migration Path**
   - Ensure existing Telegram users continue working
   - Data migration scripts for user_id generation (run once during deployment)

3. **Documentation**
   - API documentation (OpenAPI/Swagger spec)
   - Web app user guide
   - Developer setup guide
   - Deployment guide

4. **Error Handling & Logging**
   - Structured logging for API requests
   - Error codes and messages
   - Client-side error reporting

5. **Security**
   - HTTPS enforcement
   - CSRF protection
   - Rate limiting per user
   - Input sanitization
   - SQL injection prevention (via GORM)
   - XSS prevention

6. **Performance**
   - API response caching
   - Database query optimization
   - Asset minification (optional)
   - Gzip compression

7. **Testing**
   - Integration tests for API endpoints
   - E2E tests for web app flows
   - Load testing for API

**Success Criteria**:

- Telegram bot continues to work unchanged
- Web app is fully functional
- API is documented
- All tests pass
- Security audit complete
- Performance benchmarks met

**Status**: Not Started

---

## Technical Decisions

### API Design

- **REST over GraphQL**: Simpler, better browser support, easier to debug
- **JWT over sessions**: Stateless, better for scaling, works well with SPA
- **API versioning**: `/api/v1/...` to allow future changes
- **JSON only**: No XML or other formats needed

### Authentication

- **Passwordless magic links**: No password management, better UX and security
- **Magic link expiry**: 30 days (configurable)
- **JWT with refresh tokens**: Balance of security and UX
- **Dual auth modes**: Support both Telegram and passwordless email
- **Token expiry**: Access token 15min, refresh token 7 days
- **Email sending**: Use existing email service (SMTP) for magic links

### Database

- **Extend existing schema**: Minimal changes to existing tables
- **GORM migrations**: Use existing migration system
- **Indexes**: Add indexes for common queries (user_id, date ranges)

### Frontend

- **Pure HTML/CSS/JS**: No build step, no framework lock-in
- **ES6 modules**: Modern JavaScript with modules
- **Progressive enhancement**: Works without JS for basic views
- **LocalStorage**: For JWT and basic caching

### Code Organization

- **Separation of concerns**: API, web, client (Telegram) are separate
- **Shared business logic**: Repository layer used by all
- **Middleware pattern**: For auth, logging, CORS, etc.

---

## Migration Strategy

### Phase 1: API-first

1. Build API without breaking existing functionality
2. Telegram bot continues to use existing code
3. New API runs alongside old web dashboard

### Phase 2: Web client

1. Build new web client consuming API
2. Deploy to `/app` path (old dashboard at `/web`)
3. Gradual user migration

### Phase 3: Consolidation

1. Consider migrating Telegram bot to use API internally (optional)
2. Deprecate old web dashboard
3. Unified codebase

---

## Risks & Mitigations

### Risk: Breaking existing Telegram functionality

**Mitigation**: Keep `internal/client/` unchanged, only share repository layer

### Risk: Authentication complexity

**Mitigation**: Start simple (JWT only), add complexity incrementally

### Risk: Frontend complexity without framework

**Mitigation**: Keep UI simple, use web standards, progressive enhancement

### Risk: API performance issues

**Mitigation**: Add caching early, database indexes, pagination

### Risk: Security vulnerabilities

**Mitigation**: Security review at each stage, rate limiting, input validation

---

## Timeline Estimates

- **Stage 1**: 3-5 days (API foundation + auth)
- **Stage 2**: 3-4 days (Transaction CRUD API)
- **Stage 3**: 2-3 days (Statistics & reporting API)
- **Stage 4**: 5-7 days (Web client)
- **Stage 5**: 2-3 days (Integration & polish)

**Total**: 15-22 days (single developer, full-time)

---

## Definition of Done

### For API

- [ ] All endpoints documented (OpenAPI spec)
- [ ] Tests written and passing
- [ ] Error handling implemented
- [ ] Rate limiting active
- [ ] CORS configured
- [ ] Logging in place

### For Web Client

- [ ] All pages functional
- [ ] Mobile responsive
- [ ] Cross-browser tested
- [ ] Error states handled
- [ ] Loading states implemented
- [ ] Accessible (WCAG 2.1 AA)

### For Overall Project

- [ ] Telegram bot still works
- [ ] Web app fully functional
- [ ] Documentation complete
- [ ] Security review done
- [ ] Performance benchmarks met
- [ ] Deployment guide ready

### For User Reconciliation

- [ ] Existing Telegram users retain full functionality
- [ ] New web users can register and use app independently
- [ ] Manually linked users (via admin DB updates) can access same data from both interfaces
- [ ] All user lookups work correctly (by tg_id, email, user_id)
- [ ] Transactions are properly associated with user via user_id
- [ ] Database check constraint prevents users without any auth method

---

## User Reconciliation Examples

### Example 1: Existing Telegram User Gets Web Access (Manual)

```
Initial State:
  user_id: abc-123-def
  tg_id: 98765432
  email: NULL
  magic_link_token: NULL
  magic_link_expires: NULL

Admin Action: Manually adds email to enable web access
  → Updates database:
    UPDATE users
    SET email = 'telegram_user@example.com'
    WHERE user_id = 'abc-123-def';

User Action: Visits web app and requests magic link
  → Enters email: telegram_user@example.com
  → System sends magic link to email
  → magic_link_token and magic_link_expires are set

Final State:
  user_id: abc-123-def
  tg_id: 98765432
  email: telegram_user@example.com
  magic_link_token: abc123xyz...
  magic_link_expires: 2025-12-20

Result: Can now access via both Telegram bot AND web login (passwordless)
```

### Example 2: New Web User Gets Telegram Access (Manual)

```
Initial State:
  user_id: xyz-456-ghi
  tg_id: NULL
  email: web_user@example.com
  magic_link_token: token789...
  magic_link_expires: 2025-12-20

Admin Action: Manually adds Telegram ID
  → Identifies user's Telegram ID (98765432)
  → Updates database:
    UPDATE users
    SET tg_id = 98765432
    WHERE user_id = 'xyz-456-ghi';

Final State:
  user_id: xyz-456-ghi
  tg_id: 98765432
  email: web_user@example.com
  magic_link_token: token789...
  magic_link_expires: 2025-12-20

Result: Can now access via both web login (magic link) AND Telegram bot
```

### Example 3: Web-Only User (Never Links Telegram)

```
State:
  user_id: qwe-789-rty
  tg_id: NULL
  email: web_only@example.com
  magic_link_token: tokenABC...
  magic_link_expires: 2025-12-20

Access:
  - Can only login via web (passwordless magic link)
  - Cannot use Telegram bot
  - All transactions linked via user_id
  - Fully functional web experience
  - Must request new magic link every 30 days
```

---

## Quick Reference: User Reconciliation Key Points

### Database Schema Changes Summary

```sql
-- Users table additions
ALTER TABLE users ADD COLUMN user_id UUID NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE users ADD COLUMN email VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN magic_link_token VARCHAR(255);
ALTER TABLE users ADD COLUMN magic_link_expires TIMESTAMP;
ALTER TABLE users ADD COLUMN last_login TIMESTAMP;
ALTER TABLE users ADD CONSTRAINT check_auth_method
  CHECK (tg_id IS NOT NULL OR email IS NOT NULL);
CREATE UNIQUE INDEX idx_users_user_id ON users(user_id);
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;

-- Transactions table additions
ALTER TABLE transactions ADD COLUMN user_id UUID;
UPDATE transactions SET user_id = (SELECT user_id FROM users WHERE users.tg_id = transactions.tg_id);
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_user_id
  FOREIGN KEY (user_id) REFERENCES users(user_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
```

### API Endpoints Summary

```
Authentication (Passwordless):
POST   /api/v1/auth/register          - Create new user (email only)
POST   /api/v1/auth/login             - Request magic link via email
GET    /api/v1/auth/verify?token=...  - Verify magic link and get JWT
POST   /api/v1/auth/refresh           - Refresh access token
POST   /api/v1/auth/logout            - Invalidate token
GET    /api/v1/auth/me                - Get current user info

Transactions:
POST   /api/v1/transactions           - Create transaction
GET    /api/v1/transactions/:id       - Get transaction by ID
PUT    /api/v1/transactions/:id       - Update transaction
DELETE /api/v1/transactions/:id       - Delete transaction
GET    /api/v1/transactions           - List transactions (paginated, filterable)
POST   /api/v1/transactions/search    - Advanced search
POST   /api/v1/transactions/import    - Bulk import (CSV)

Statistics:
GET    /api/v1/stats/monthly          - Monthly statistics
GET    /api/v1/stats/weekly           - Weekly statistics
GET    /api/v1/stats/yearly           - Yearly statistics
GET    /api/v1/stats/categories       - Category breakdown

Export:
GET    /api/v1/export/csv             - Export to CSV
```

### Repository Method Additions

```go
// User repository extensions
GetUserByEmail(email string) (*User, error)
GetUserByUserID(userID string) (*User, error)
GetUserByIdentifier(identifier string) (*User, error) // Tries user_id, email, or tg_id

// Transaction repository - update existing methods to use user_id
GetUserTransactionsByDateRange(userID string, start, end time.Time) ([]Transaction, error)
CreateTransaction(userID string, transaction *Transaction) error
UpdateTransaction(userID string, transactionID int64, updates *Transaction) error
DeleteTransaction(userID string, transactionID int64) error
// ... other methods updated to use user_id instead of tg_id
```

### Key Implementation Notes

1. **Backward Compatibility**: Existing Telegram users continue working without any changes
2. **Passwordless Authentication**: Web users use magic links (no passwords to manage)
3. **Magic Link Expiry**: Links valid for 30 days, users must request new link after expiration
4. **Flexible Authentication**: Users can have Telegram-only, email-only, or both (manually linked)
5. **Data Unification**: Single user record accessed via multiple identifiers (tg_id, email, user_id)
6. **Security**: Database constraint ensures at least one authentication method at all times
7. **Migration**: Generate `user_id` for all existing users, update FK from `tg_id` to `user_id`
8. **JWT Claims**: Use `user_id` in JWT tokens for universal identification
9. **Manual Reconciliation**: Account linking is done by admins via direct database updates, not through API or bot commands
