# Feature Specification: Monthly Budget & Alerts

**Feature Branch**: `002-budgets-alerts`
**Created**: 2026-05-16
**Status**: Implemented (v1, single user-wide budget)
**Input**: "Single monthly budget per user with progress tracking and threshold alerts via Telegram, and dashboard editor."

## Scope (v1)

A single **user-wide** monthly expense budget (one row per user, no per-category granularity). Per-category budgets are explicitly out of scope for v1; revisit only if user feedback demands it.

## User Scenarios & Testing

### User Story 1 - Set the monthly budget (P1)

User defines a single monthly cap for total expenses (e.g. 1500€). From then on, every recap and every new Expense transaction shows progress against this budget.

**Acceptance Scenarios**:

1. `/budget set 1500` stores the budget and confirms.
2. `/budget` shows current budget with current-month progress (spent/limit, percentage, indicator).
3. After upsert, `/budget set 1800` updates the same row (no duplicates).
4. `/budget delete` removes the budget.
5. Adding any Expense (any category) with an active budget appends `Budget: spent/limit (pct%)` to the confirmation message.

### User Story 2 - Threshold alerts (P1)

When the user's cumulative monthly expense crosses 80% or 100% of the budget, the bot proactively warns. Each threshold fires at most once per calendar month.

**Acceptance Scenarios**:

1. Crossing 80% (but below 100%) appends "⚠️ Approaching monthly budget (80% used)".
2. Crossing 100% appends "🚨 Over budget by X €".
3. Crossing the same threshold a second time within the same month does NOT re-fire.
4. New month → state resets and alerts can fire again.

### User Story 3 - Budget in monthly recap (P2)

The automated monthly recap (1st of each month) includes a single line `📊 Budget: spent/limit (pct%) ✅/⚠️/🚨` for the just-closed month.

### User Story 4 - Dashboard Budget tab (P2)

The web dashboard gains a "Budget" tab where the user can view the current month's progress (with a progress bar) and create/update/remove their budget.

### Edge Cases

- **Income categories**: budgets only consider Expense transactions.
- **Currency**: stored in EUR only.
- **Mid-month creation**: progress reflects current-month-to-date spending. The 80%/100% alerts will fire on the next crossing transaction in that month if applicable.
- **Edit/delete of past txs**: do NOT re-fire past alerts; only forward crossings.

## Requirements

- **FR-001**: A single budget row per user (`UNIQUE(tg_id)`).
- **FR-002**: Telegram: `/budget`, `/budget set <amount>`, `/budget delete`, plus a home keyboard "📊 Budget" button.
- **FR-003**: On Expense insert, append progress and any newly-crossed alerts to the bot reply.
- **FR-004**: Alert state persisted; each `(tg_id, year_month, threshold)` fires at most once.
- **FR-005**: Monthly recap message includes a Budget line when a budget exists.
- **FR-006**: REST endpoints under `/web/api/budget` (GET/POST/DELETE) gated by auth.
- **FR-007**: Dashboard "Budget" tab with view + create/update/delete and a progress bar.

## Entities

- **Budget**: `id, tg_id (unique), amount (decimal 15,2), currency (EUR), timestamps`.
- **BudgetAlert**: `id, tg_id, year_month (CHAR(7)), threshold (smallint, 80|100), fired_at`. Unique on `(tg_id, year_month, threshold)`.

## Success Criteria

- **SC-001**: Set/update budget in ≤2 messages.
- **SC-002**: Every Expense reply with an active budget shows `spent/limit (pct%)`.
- **SC-003**: 80% and 100% each fire exactly one alert per month.
- **SC-004**: Monthly recap shows budget outcome for the closed month.
- **SC-005**: Web Budget tab CRUDs work without page reload errors.
