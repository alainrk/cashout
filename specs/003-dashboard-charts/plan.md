# Implementation Plan: Dashboard Visual Analytics

**Branch**: `003-dashboard-charts` | **Date**: 2026-05-16 | **Spec**: [spec.md](spec.md)

## Summary

Add three new aggregated analytics endpoints (monthly, trend, year) and a charting layer (Chart.js v4, ~70 KB gzipped, MIT) to the existing dashboard. New "Trends" and "Year" tabs; donut + income/expense bar added to the existing month view. No schema changes, no new dependencies on the Go side.

## Technical Context

**Language/Version**: Go 1.26 (backend), vanilla JS ES2020 (frontend)
**Primary Dependencies**: existing GORM repo layer; **Chart.js 4.x** added as a single vendored file in `web/static/js/vendor/chart.umd.min.js`
**Storage**: none new — read-only aggregations off existing `transactions` table.
**Target Platform**: same as today (Linux server + browser).
**Constraints**: NFR-001 (≤100 KB JS), NFR-002 (≤200 ms p95), NFR-003 (server-side aggregation only).

## Choice of charting library

- **Chart.js 4** chosen for: tiny footprint (~70 KB gzipped UMD), MIT, no build step required (vendored UMD), covers donut/bar/line/stacked-bar.
- Alternatives considered & rejected:
  - **ApexCharts**: bigger (~200 KB), heavier theming overhead.
  - **D3**: too low-level for this scope.
  - **ECharts**: powerful but ~340 KB; overkill.
- Loaded as a plain `<script src="...">` (no npm/bundler), keeping with the project's no-build frontend.

## Project Structure

```text
specs/003-dashboard-charts/
├── plan.md
├── spec.md
└── contracts/
    └── analytics-api.md    # JSON shapes

internal/
├── db/
│   └── analytics.go        # NEW: SQL aggregations (monthly, trend, year)
├── repository/
│   └── analytics.go        # NEW: thin wrapper
└── web/
    ├── analytics.go        # NEW: 3 handlers + DTOs
    └── web.go              # MODIFY: register 3 routes

web/
├── templates/
│   ├── dashboard.html      # MODIFY: add chart canvases + Trends/Year tabs
│   └── partials/           # OPTIONAL: split if file gets too big
└── static/
    ├── js/
    │   ├── vendor/
    │   │   └── chart.umd.min.js   # NEW: vendored Chart.js 4
    │   ├── charts.js              # NEW: chart rendering glue
    │   └── analytics.js           # NEW: fetch + cache analytics responses
    └── css/
        └── dashboard.css   # MODIFY: chart container styles
```

## Phase 0 — Research / Open Questions

1. **Library**: Chart.js 4 confirmed. Pin exact version (e.g. `4.4.3`) and SRI-hash it in the template.
2. **Aggregation in SQL vs Go**: prefer SQL `GROUP BY` with `date_trunc('month', date)` for trend and year endpoints — pushes work to Postgres, leverages existing indexes on `(tg_id, date)`.
3. **Caching**: skip HTTP caching for v1 (single-user scale); revisit if NFR-002 misses.
4. **Time zones**: aggregation in UTC (matches current `Transaction.Date` type `date`). Document in the API contract.

## Phase 1 — Analytics Endpoints

All under `/web/api/analytics/*`, auth-required (use existing `requireAuth` middleware). User scope is enforced by injecting `user.TgID` into the WHERE clause.

### `GET /web/api/analytics/monthly?month=YYYY-MM`

Response:

```json
{
  "month": "2026-05",
  "totalIncome": 2400.00,
  "totalExpenses": 1850.50,
  "balance": 549.50,
  "byCategory": {
    "Expense": [
      { "category": "Grocery", "amount": 380.00, "count": 12, "pct": 20.53 },
      { "category": "EatingOut", "amount": 220.00, "count": 8, "pct": 11.89 }
    ],
    "Income":  [
      { "category": "Salary", "amount": 2400.00, "count": 1, "pct": 100.0 }
    ]
  }
}
```

Powered by existing `GetMonthCategorizedTotals` + a single SQL query for counts.

### `GET /web/api/analytics/trend?months=12`

Response:

```json
{
  "from": "2025-06",
  "to":   "2026-05",
  "points": [
    { "month": "2025-06", "income": 2400, "expense": 1800, "balance": 600 },
    ...
  ]
}
```

SQL (single round-trip):

```sql
SELECT
  to_char(date, 'YYYY-MM') AS ym,
  type,
  SUM(amount)::float AS total
FROM transactions
WHERE tg_id = $1 AND date >= $2 AND date <  $3
GROUP BY ym, type
ORDER BY ym;
```

Pivot in Go to per-month income/expense/balance, padding empty months with zero so the line is continuous.

### `GET /web/api/analytics/year?year=YYYY`

Response:

```json
{
  "year": 2026,
  "totalIncome": 26000,
  "totalExpenses": 21000,
  "balance": 5000,
  "byMonth": [
    { "month": 1, "income": 2400, "expense": 1700, "balance": 700 },
    ... 12 entries
  ],
  "byCategory": {
    "Expense": [ {"category":"Grocery","amount":3500,"count":120,"pct":16.67}, ... ],
    "Income":  [ {"category":"Salary","amount":24000,"count":12,"pct":92.3}, ... ]
  }
}
```

Reuses `GetYearCategorizedTotals` + same SQL pattern as trend, filtered to one year.

### Contracts file

`contracts/analytics-api.md` documents these JSON shapes and the error envelope (existing `sendJSONError` pattern).

## Phase 2 — Frontend

### Vendored library

Add `web/static/js/vendor/chart.umd.min.js` (pinned 4.4.3). Referenced once in `dashboard.html` with `integrity="sha384-..."` and `crossorigin="anonymous"`.

### `analytics.js` — data layer

- `async function fetchMonthly(ym)` / `fetchTrend(n)` / `fetchYear(y)` — each returns the JSON payload.
- Per-tab in-memory cache `{ ym → payload }` to avoid refetching when user toggles tabs.

### `charts.js` — render layer

Three reusable factories:

- `renderCategoryDonut(ctx, byCategoryExpense, { onSliceClick })`
- `renderIncomeExpenseBar(ctx, income, expense)`
- `renderTrendLine(ctx, points, { onPointClick })`
- `renderYearStacked(ctx, byMonth)`

Each function:
1. Destroys the previous chart instance on that canvas (Chart.js requires this on re-render).
2. Applies a shared palette (`web/static/js/palette.js` or inline in `charts.js`) so categories use the same colour everywhere.
3. Returns the new Chart instance.

### Template changes

In `web/templates/dashboard.html`:

- Add canvases on the Transactions page:
  - `<canvas id="incomeExpenseBar" height="120">` in the stats row.
  - `<canvas id="categoryDonut" height="240">` above the transaction list, with adjacent legend.
- New nav tabs `Trends` and `Year`, each with:
  - Trends page: `<canvas id="trendLine">` + a `<select>` for window size (3 / 6 / 12 months).
  - Year page: `<select>` for year + `<canvas id="yearStacked">` + `<canvas id="yearDonut">`.

### Donut → list filter

On slice click, set a `categoryFilter` state and re-render the existing transaction list filtering on `tx.category === selected`. No new round-trip — list is already fetched.

### Trend → month navigation

On point click, set `window.location = '/web/dashboard?month=' + clickedMonth`.

## Phase 3 — Testing

- **Backend**: table-driven Go tests for the aggregation pivot logic (especially empty-month padding and unknown-category handling). Mock the DB layer.
- **DB**: a small integration test against a seeded DB asserting query shape and totals match a hand-computed expectation.
- **Frontend**: manual cross-browser QA (Safari, Chrome, Firefox) + mobile width DevTools at 360, 640, 1024px. Empty-state checks.
- **Perf**: time `analytics/monthly` and `analytics/trend` on the seeded 5-year DB; assert p95 ≤ NFR-002 with `wrk` or simple repeat curl.

## Phase 4 — Rollout

1. Land backend endpoints behind no flag (read-only, low risk).
2. Land the vendored Chart.js commit separately so it can be reverted cleanly if needed.
3. Land template + JS changes.
4. Update README "Visual Analytics" section to match what's actually shipped.

## Risk & Mitigations

- **JS bloat**: monitor gzipped size in CI (`make assets-size` target, optional). Cap at 100 KB per NFR-001.
- **Library SRI break on upgrade**: pin version + SRI hash; document upgrade procedure.
- **Slow queries on huge histories**: indices already exist on `(tg_id, date)` and `(tg_id, type, date)` via GORM tags; verify with `EXPLAIN`.
- **Empty-state UX**: ensure every chart factory handles the empty array path explicitly to avoid a broken canvas.

## Out of Scope

- Custom dashboards / user-configurable charts.
- Drill-down beyond category filter on the donut.
- Sankey or cash-flow diagrams.
- Multi-currency aggregation.
- Server-side image rendering / PDF export of charts.

## Rough Effort

- Backend endpoints + pivot logic + tests: 1 day
- Vendoring + chart factories: 0.5 day
- Template + tabs + interactions: 1 day
- QA + perf check + README update: 0.5 day
**Total: ~3 days**
