# Feature Specification: Duplicate Transaction

**Feature Branch**: `001-duplicate-transaction`  
**Created**: 2026-04-03  
**Status**: Draft  
**Input**: User description: "Implement transaction duplication flow with clone command, home keyboard shortcut, smart transaction picker (latest expenses, wizard with type/category/search filtering), inline edit before save, today's date default"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Quick Clone from Home Keyboard (Priority: P1)

The user taps "Clone" on the home keyboard. The bot shows the 10 most recent expenses. The user taps one, reviews the cloned transaction (with today's date), optionally edits any field, then confirms to save.

**Why this priority**: Most common use case — users frequently re-enter recurring expenses (coffee, lunch, groceries). This provides the fastest path to duplicate a transaction.

**Independent Test**: Can be fully tested by tapping "Clone" on home, selecting a recent transaction, and confirming. Delivers immediate value for repeat expenses.

**Acceptance Scenarios**:

1. **Given** user has existing transactions, **When** they tap "Clone" on home keyboard, **Then** bot shows the 10 most recent expenses sorted by date descending with numbered selection buttons.
2. **Given** the recent expenses list is shown, **When** user taps a transaction number, **Then** bot creates a clone with today's date and shows the confirmation UI with edit options (date, category, amount, description, type) and confirm/cancel buttons.
3. **Given** the clone confirmation is shown, **When** user taps "Confirm", **Then** the transaction is saved to the database with today's date and all other fields copied, and user returns to home.
4. **Given** the clone confirmation is shown, **When** user taps "Edit date" and enters a new date, **Then** the date is updated and the confirmation UI refreshes.
5. **Given** user has no transactions, **When** they tap "Clone", **Then** bot shows "No transactions to clone" with a home button.

---

### User Story 2 - Clone via Wizard Search (Priority: P2)

When the recent expenses list doesn't contain the desired transaction, the user taps "Search More" to enter a wizard: first filtering by type (Income/Expense), then by category or free text search, to find and clone any historical transaction.

**Why this priority**: Covers the case where the transaction to clone is older or an income. Builds on P1's clone confirmation flow.

**Independent Test**: Can be tested by tapping Clone → Search More → selecting type → category → selecting a transaction → confirming clone.

**Acceptance Scenarios**:

1. **Given** user is on the recent expenses list, **When** they tap "Search More", **Then** bot shows type selection: "Expenses" / "Incomes" / "All".
2. **Given** user selected a type, **When** bot shows category selection, **Then** it shows relevant categories (expense categories for Expense, income categories for Income, all for All) plus "All Categories".
3. **Given** user selected a category, **When** bot shows search prompt, **Then** user can enter free-text search or tap "Show All" to see all transactions in that category.
4. **Given** search results are shown, **When** user selects a transaction, **Then** bot shows the same clone confirmation UI as P1 (today's date, edit options, confirm).
5. **Given** search results are paginated, **When** user navigates pages, **Then** pagination works correctly with Previous/Next buttons.

---

### User Story 3 - Clone via /clone Command (Priority: P3)

User can type `/clone` to enter the clone flow directly. Behaves identically to tapping "Clone" on the home keyboard.

**Why this priority**: Provides command-line equivalent for power users who prefer typing. Minimal additional work since it reuses P1's entry point.

**Independent Test**: Can be tested by typing `/clone` and verifying the same recent expenses list appears.

**Acceptance Scenarios**:

1. **Given** user types `/clone`, **When** bot processes the command, **Then** it shows the same recent expenses list as the home keyboard "Clone" button.

---

### User Story 4 - Clone Intent via Natural Language (Priority: P3)

User can type natural language like "clone", "duplicate", "repeat last transaction" and the LLM intent classifier routes them to the clone flow.

**Why this priority**: Consistent with existing bot behavior where natural language is classified into intents. Low effort since it only requires adding a new intent to the classifier.

**Independent Test**: Can be tested by typing "clone my last expense" and verifying the clone flow starts.

**Acceptance Scenarios**:

1. **Given** user types "duplicate a transaction", **When** the LLM classifies intent, **Then** it returns "clone" intent and routes to the clone flow.

---

### Edge Cases

- What happens when the cloned transaction has a currency other than EUR? → Currency is copied as-is from the original.
- What happens when user cancels mid-flow? → State resets to Normal, home keyboard shown.
- What happens when the source transaction was deleted between selection and confirmation? → Show error "Transaction no longer exists" with home button.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST add a "Clone" button to the home keyboard.
- **FR-002**: System MUST show the 10 most recent expenses when clone flow starts.
- **FR-003**: System MUST provide a "Search More" wizard with type → category → text search filtering.
- **FR-004**: System MUST create a clone with today's date and all other fields copied from the original.
- **FR-005**: System MUST allow inline editing of all fields (date, category, amount, description) before confirming the clone.
- **FR-006**: System MUST register `/clone` as a bot command.
- **FR-007**: System MUST add "clone" as a recognizable intent in the LLM intent classifier.
- **FR-008**: System MUST handle the case where the source transaction no longer exists.

### Key Entities

- **Transaction**: Existing entity — the clone creates a new Transaction with a new ID, today's date, and all other fields copied from the source.
- **UserSession**: Existing entity — new states needed for clone flow (selecting, searching, editing clone).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: User can duplicate a recent expense in 3 taps (Clone → select → confirm).
- **SC-002**: User can find and clone any historical transaction via the wizard.
- **SC-003**: Clone flow is consistent with existing edit/delete search patterns (same category keyboard, same pagination, same search UX).
- **SC-004**: All cloned transactions default to today's date.
