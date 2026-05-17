# Feature Specification: Dashboard Visual Analytics

**Feature Branch**: `003-dashboard-charts`
**Created**: 2026-05-16
**Status**: Draft
**Input**: "Add real charts to the web dashboard: monthly income vs expenses, category breakdown donut, 12-month trend, and (if 002 is shipped) budget progress bars."

## Background

The README advertises "Visual Analytics: Category breakdowns and trends" but `internal/web/dashboard.go` only returns flat lists and aggregate numbers. There are no charts in the current dashboard. This feature closes the gap.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Monthly category donut (Priority: P1)

When viewing a month, the user sees a donut chart of expense distribution by category, with percentages, hover tooltips, and a legend.

**Why this priority**: Single most useful at-a-glance view of where money goes; complements the existing stats numbers.

**Independent Test**: Open dashboard for a month with transactions → donut renders showing each expense category proportionally; clicking a slice filters the transaction list below by that category.

**Acceptance Scenarios**:

1. **Given** a month has expense transactions, **When** the dashboard loads, **Then** an "Expenses by Category" donut chart renders with one slice per category present.
2. **Given** the chart is rendered, **When** user hovers a slice, **Then** a tooltip shows `Category · €X.XX · YY%`.
3. **Given** the chart is rendered, **When** user clicks a slice, **Then** the transaction list filters to that category and the chart slice is highlighted.
4. **Given** a month has no expenses, **When** the dashboard loads, **Then** the chart area shows "No data for this month".

---

### User Story 2 - Income vs Expenses bar (Priority: P1)

A small bar chart shows Income vs Expenses for the selected month with the resulting balance.

**Independent Test**: Month with 2000 income, 1500 expense → two bars labelled accordingly with a "Balance: +€500" caption.

**Acceptance Scenarios**:

1. **Given** the month has data, **When** the dashboard loads, **Then** two stacked or side-by-side bars render for Income and Expenses with currency labels.
2. **Given** balance is negative, **When** rendered, **Then** the balance label is red; otherwise green.

---

### User Story 3 - 12-month trend (Priority: P2)

A line chart on a dedicated "Trends" tab shows monthly income, monthly expenses, and monthly balance over the trailing 12 months.

**Independent Test**: Open Trends tab → 3 lines render over 12 month labels, with the latest month being the currently selected one.

**Acceptance Scenarios**:

1. **Given** at least 1 month of data, **When** the Trends tab is opened, **Then** a line chart renders with up to 12 monthly points.
2. **Given** user hovers a point, **When** the tooltip displays, **Then** it shows month + income + expense + balance.
3. **Given** the user clicks a month label, **When** clicked, **Then** the dashboard navigates to that month's transactions view.

---

### User Story 4 - Top categories table (Priority: P2)

Below the donut, a sortable table lists the top 5 expense categories for the month with: category, amount, count, % of total, and (if budgets exist) budget usage.

**Acceptance Scenarios**:

1. **Given** the month has expense data, **When** the dashboard loads, **Then** a top-5 table renders sorted by amount desc.
2. **Given** budgets exist for some categories, **When** the table renders, **Then** a "Budget" column shows `spent/limit` for those categories and a dash for others.

---

### User Story 5 - Year overview tab (Priority: P3)

A "Year" tab shows a year-to-date summary: total income, total expenses, balance, monthly stacked bar chart for the calendar year, and yearly category donut.

**Acceptance Scenarios**:

1. **Given** the Year tab is opened, **When** the page loads, **Then** YTD totals are shown alongside a monthly stacked bar (income green / expense red) for Jan→Dec.
2. **Given** the year selector is changed, **When** changed, **Then** charts re-render for the chosen year.

---

### Edge Cases

- Empty months / new users: every chart shows a friendly empty state, never an error.
- Very long category names: legend truncates with ellipsis but full name on hover.
- Mobile width (<640px): charts stack vertically and remain readable; donut shrinks but legend stays usable.
- Currency: all charts display the user's existing default (EUR) — out of scope to mix currencies.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide an aggregated JSON endpoint `GET /web/api/analytics/monthly?month=YYYY-MM` returning category totals (Expense + Income), totals, and balance.
- **FR-002**: System MUST provide `GET /web/api/analytics/trend?months=12` returning per-month income/expense/balance arrays.
- **FR-003**: System MUST provide `GET /web/api/analytics/year?year=YYYY` returning per-month totals and yearly category breakdown.
- **FR-004**: Dashboard MUST render a category donut for the currently selected month above the transaction list.
- **FR-005**: Dashboard MUST render an Income vs Expenses bar in the stats area.
- **FR-006**: Dashboard MUST add a "Trends" nav tab with a 12-month line chart.
- **FR-007**: Dashboard MUST add a "Year" nav tab with stacked bar + yearly donut.
- **FR-008**: Charts MUST be interactive: hover tooltips and click-through behaviour for donut→category filter and trend→month navigation.
- **FR-009**: All chart endpoints MUST require auth and scope to the authenticated user.
- **FR-010**: Frontend MUST use a single client-side chart library; no server-rendered images.

### Non-Functional Requirements

- **NFR-001**: Total added JS payload (chart lib + dashboard glue) ≤ 100 KB gzipped.
- **NFR-002**: `analytics/monthly` p95 latency ≤ 200 ms on a user with 5 years of data.
- **NFR-003**: No client-side iteration over per-transaction arrays just to compute aggregates — server returns aggregates ready to plot.

### Key Entities

- No new persistent entities. New read-only DTOs in `internal/web/analytics.go`.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: User can see the category donut for any month in ≤1 click.
- **SC-002**: 12-month trend renders in ≤500 ms after tab click on a user with full seed data.
- **SC-003**: Clicking a donut slice filters the transaction list in ≤100 ms client-side (no extra round-trip).
- **SC-004**: README "Visual Analytics" claim is satisfied — every bullet is implemented.
