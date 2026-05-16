# Implementation Plan: Budgets & Alerts

**Branch**: `002-budgets-alerts` | **Date**: 2026-05-16 | **Spec**: [spec.md](spec.md)

## Summary

Add per-category monthly budgets for expense categories. Persist budgets and per-month per-threshold alert state. Hook into the existing add-transaction confirmation flow to display progress and fire one-shot 80%/100% threshold alerts. Surface in `/budgets` command, home keyboard, monthly recap, and a new Budgets tab on the web dashboard.

## Technical Context

**Language/Version**: Go 1.26
**Primary Dependencies**: gotgbot/v2, GORM, gocron, OpenAI-compatible LLM
**Storage**: PostgreSQL — two new tables (`budgets`, `budget_alerts`)
**Target Platform**: Linux server (Docker)
**Constraints**: Telegram callback data 64-byte limit; keep aggregation EUR-only for now to match existing assumptions.
**Scale**: Single-user personal-finance scale.

## Project Structure

### Documentation

```text
specs/002-budgets-alerts/
├── plan.md          # this file
├── spec.md
├── data-model.md
├── research.md      # OPTIONAL — open questions only
└── contracts/
    └── web-api.md   # REST contract for /web/api/budgets
```

### Source Code Changes

```text
internal/
├── migrations/versions/
│   └── 011_create_budgets.go          # NEW
├── model/
│   └── budgets.go                     # NEW: Budget, BudgetAlert structs + helpers
├── db/
│   └── budgets.go                     # NEW: CRUD + GetSpentForMonth, GetAlertFired, SetAlertFired
├── repository/
│   ├── repository.go                  # MODIFY: register Budgets repo
│   └── budgets.go                     # NEW: thin wrapper
├── client/
│   ├── budgets.go                     # NEW: /budgets, /budget set|delete, list keyboard, set wizard
│   ├── send.go                        # MODIFY: add "Budgets" home button
│   ├── setup.go                       # MODIFY: register handlers + commands
│   ├── router.go                      # MODIFY: route new states + intent
│   ├── transactions.go                # MODIFY: after Expense save → call budget.Evaluate + append message
│   └── helpers.go                     # MODIFY: shared format helper for progress lines
├── scheduler/
│   └── monthly.go                     # MODIFY: append budgets section to recap
├── ai/
│   ├── prompt.go                      # MODIFY: add "budget" to intent classifier
│   └── deepseek.go                    # MODIFY: IntentBudget constant
├── model/
│   └── users.go                       # MODIFY: add BudgetSetWaitAmount state
└── web/
    ├── budgets.go                     # NEW: handleAPIBudgets, Create, Delete + Budgets page
    ├── web.go                         # MODIFY: register routes
    └── templates/dashboard.html       # MODIFY: add Budgets tab + JS

web/static/
├── css/dashboard.css                  # MODIFY: progress bar styles
└── js/budgets.js                      # NEW
```

## Phase 0 — Research / Open Questions

1. **Currency**: lock to EUR-only in v1 (consistent with existing aggregation). Multi-currency deferred.
2. **Where to call `budget.Evaluate`**: after successful insert in the **confirm transaction** path only (not during inline edits before save). Single call site in `client/transactions.go`.
3. **Threshold model**: hard-coded `[80, 100]` for v1. User-customisable thresholds deferred.
4. **LLM intent**: extend the intent enum but keep parsing of `/budget set X N` deterministic (regex), not LLM-routed. LLM only routes free-text "set my grocery budget to 400" → command.

## Phase 1 — Data Model

### Migration `011_create_budgets.go`

```sql
CREATE TABLE budgets (
  id          BIGSERIAL PRIMARY KEY,
  tg_id       BIGINT NOT NULL,
  category    transaction_category NOT NULL,
  amount      DECIMAL(15,2) NOT NULL CHECK (amount > 0),
  currency    currency_type NOT NULL DEFAULT 'EUR',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tg_id, category)
);
CREATE INDEX idx_budgets_tg ON budgets(tg_id);

CREATE TABLE budget_alerts (
  id          BIGSERIAL PRIMARY KEY,
  tg_id       BIGINT NOT NULL,
  category    transaction_category NOT NULL,
  year_month  CHAR(7) NOT NULL,           -- '2026-05'
  threshold   SMALLINT NOT NULL CHECK (threshold IN (80, 100)),
  fired_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tg_id, category, year_month, threshold)
);
CREATE INDEX idx_budget_alerts_tg_month ON budget_alerts(tg_id, year_month);
```

Rollback drops both tables.

### Go structs (`internal/model/budgets.go`)

```go
type Budget struct {
    ID        int64
    TgID      int64
    Category  TransactionCategory
    Amount    float64
    Currency  CurrencyType
    CreatedAt time.Time
    UpdatedAt time.Time
}
type BudgetAlert struct {
    ID        int64
    TgID      int64
    Category  TransactionCategory
    YearMonth string  // "2026-05"
    Threshold int16
    FiredAt   time.Time
}
```

## Phase 2 — Repository / DB

`internal/db/budgets.go`:

- `UpsertBudget(b *Budget)`
- `DeleteBudget(tgID int64, cat TransactionCategory)`
- `ListBudgets(tgID int64) ([]Budget, error)`
- `GetSpentInCategoryForMonth(tgID, year, month, cat) (float64, error)` — sum of expense amounts.
- `IsAlertFired(tgID, cat, ym, threshold) (bool, error)`
- `MarkAlertFired(tgID, cat, ym, threshold) error`

Repository wrapper at `internal/repository/budgets.go` mirrors the existing pattern.

## Phase 3 — Budget Evaluation Engine

New helper (in `internal/client/budgets.go`):

```go
type BudgetProgress struct {
    Category    model.TransactionCategory
    Spent       float64
    Limit       float64
    Pct         int      // floored
    NewAlerts   []int    // subset of {80,100} that crossed on THIS insertion
}

func EvaluateAfterInsert(repo Repositories, tx model.Transaction) (*BudgetProgress, error)
```

Algorithm:
1. If `tx.Type != Expense` → return nil.
2. Load budget for `(tg_id, category)`. If absent → nil.
3. Compute `spentAfter = GetSpentInCategoryForMonth(...)` (already includes the just-inserted row).
4. Compute `spentBefore = spentAfter - tx.Amount`.
5. For each threshold `t ∈ {80, 100}`:
   - `crossedNow = spentBefore < limit*t/100 && spentAfter >= limit*t/100`
   - If `crossedNow && !IsAlertFired(...)`: `MarkAlertFired(...)` and append to `NewAlerts`.
6. Return progress struct.

Caller in `client/transactions.go` appends to the post-save confirmation message:
- always: `Budget: 75/400 (19%)` if a budget exists,
- additionally: warning lines for any `NewAlerts`.

## Phase 4 — Telegram Commands & Keyboards

- `/budgets` → list with progress lines, "Set budget" and "Back" inline buttons.
- `/budget set <Category> <Amount>` → upsert, confirm.
- `/budget delete <Category>` → delete, confirm.
- Home keyboard: add "💰 Budgets" button → opens `/budgets` view.
- New user state `BudgetSetWaitAmount` for the interactive "Set budget" wizard (category keyboard reused → amount text input).

## Phase 5 — Monthly Recap

`internal/scheduler/monthly.go`: after the existing category breakdown, append:

```
📊 Budgets
- Grocery: 380/400 (95%) ✅
- EatingOut: 220/150 (147%) 🚨
```

Driven by `ListBudgets` + `GetSpentInCategoryForMonth` for the closed month.

## Phase 6 — Web Dashboard

New endpoints in `internal/web/budgets.go`, registered in `internal/web/web.go`:

- `GET  /web/api/budgets?month=YYYY-MM` → `[{category, limit, spent, pct, currency}]`
- `POST /web/api/budgets` — body `{category, amount}` → upsert
- `DELETE /web/api/budgets` — body `{category}` → delete

Frontend:
- New nav tab "Budgets" in `dashboard.html`.
- `web/static/js/budgets.js` — fetch list, render progress bars, handle form submit.
- CSS: `.budget-bar { background: green/amber/red based on pct }`.

## Phase 7 — Testing

- Unit tests for `EvaluateAfterInsert` covering: no budget, under 80, crossing 80, crossing 100, already-fired-not-repeated, new-month-reset.
- DB tests for upsert uniqueness + alert uniqueness.
- Manual Telegram QA per spec acceptance scenarios.
- Manual web QA: create → list → delete → progress colors.

## Risk & Mitigations

- **Race condition** on alert dedup: two transactions inserted in the same second could double-fire. Mitigation: rely on `UNIQUE(tg_id, category, year_month, threshold)` — `MarkAlertFired` uses `INSERT ... ON CONFLICT DO NOTHING` and only appends to `NewAlerts` if rows affected = 1.
- **Recap drift**: monthly recap runs at calendar boundary; ensure it queries the *closed* month, not current.
- **Edit/delete of past txs** can quietly push past months over budget without re-alert — accepted per spec (FR considers only forward crossings).

## Out of Scope

- Multi-currency budgets and FX.
- Weekly/yearly budget periods.
- Budget rollover.
- Sub-category or tag-level budgets.
- LLM-driven natural-language budget editing.

## Rough Effort

- Migration + model + repo: 0.5 day
- Evaluation engine + transaction hook + unit tests: 0.5 day
- Telegram commands + keyboards + state: 0.5 day
- Monthly recap section: 0.25 day
- Web API + dashboard tab + JS/CSS: 0.75 day
- QA pass: 0.25 day
**Total: ~2.75 days**
