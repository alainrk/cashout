# Tasks: Duplicate Transaction

**Input**: Design documents from `/specs/001-duplicate-transaction/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: No automated tests requested (project uses manual Telegram bot testing).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add session state constants and AI intent that all clone stories depend on

- [ ] T001 [P] Add clone session state constants (`StateSelectingCloneTransaction`, `StateSelectingCloneSearchCategory`, `StateEnteringCloneSearchQuery`) in `internal/model/users.go`
- [ ] T002 [P] Add `IntentClone` constant in `internal/ai/deepseek.go` and add `"clone"` intent (with keywords: clone, duplicate, repeat, copy, same again) to `LLMIntentClassificationPromptTemplate` in `internal/ai/prompt.go`

**Checkpoint**: Foundation constants ready — clone flow handlers can now be implemented.

---

## Phase 2: User Story 1 — Quick Clone from Home Keyboard (Priority: P1) 🎯 MVP

**Goal**: User taps "Clone" on home keyboard → sees 10 most recent expenses → taps one → clone saved with today's date → edit UI shown.

**Independent Test**: Send `/start` → verify "Clone" button → tap "Clone" → see recent expenses → select one → verify clone saved with today's date and edit options shown.

### Implementation for User Story 1

- [ ] T003 [US1] Create `internal/client/clone.go` with `CloneTransactions` entry handler: authenticate user, set state to `StateSelectingCloneTransaction`, call `showRecentExpensesForClone`. Implement `showRecentExpensesForClone` to fetch 10 most recent expenses via `Repositories.Transactions.GetUserTransactionsPaginated` (filtered to expenses), format numbered list (same style as `formatDeletableTransactions` in `delete.go`), create keyboard with numbered selection buttons (`clone.select.{ID}`), pagination (`clone.page.{OFFSET}`), page indicator (`clone.noop`), "Search More" (`clone.searchmore`), and "Cancel" (`transactions.cancel`).
- [ ] T004 [US1] Implement `CloneTransactionSelected` handler in `internal/client/clone.go`: fetch transaction by ID from callback data (`clone.select.{ID}`), verify ownership, create clone with `time.Now()` date and all other fields copied, save via `Repositories.Transactions.Add`, set session state to `StateEditingNewTransaction` with new transaction ID in `Session.Body`, send confirmation message with edit keyboard (same format as `addTransaction` in `transactions.go`: edit description/category/date/amount buttons, delete, home).
- [ ] T005 [US1] Implement `CloneTransactionPage` and `CloneNoop` handlers in `internal/client/clone.go`: pagination for recent expenses list (same pattern as `DeleteTransactionPage`/`DeleteNoop` in `delete.go`).
- [ ] T006 [US1] Add "📋 Clone" button to home keyboard in `internal/client/send.go`: modify row 2 of `SendHomeKeyboard` from `[✏️ Edit, 🗑 Delete]` to `[✏️ Edit, 📋 Clone, 🗑 Delete]` with callback data `home.clone`. Also add same button to all other keyboard layouts that include the home keyboard (in `sendRecapWithNavigation`).
- [ ] T007 [US1] Register clone handlers in `internal/client/setup.go`: add `handlers.NewCallback(callbackquery.Equal("home.clone"), c.CloneTransactions)`, `handlers.NewCallback(callbackquery.Prefix("clone.select."), c.CloneTransactionSelected)`, `handlers.NewCallback(callbackquery.Prefix("clone.page."), c.CloneTransactionPage)`, `handlers.NewCallback(callbackquery.Equal("clone.noop"), c.CloneNoop)`, `handlers.NewCallback(callbackquery.Equal("clone.searchmore"), c.CloneSearchMore)`.

**Checkpoint**: User Story 1 complete — quick clone from home keyboard works end-to-end. Clone saves immediately, edit UI reuses existing `StateEditingNewTransaction` flow.

---

## Phase 3: User Story 2 — Clone via Wizard Search (Priority: P2)

**Goal**: From recent expenses list, user taps "Search More" → type filter (Expense/Income/All) → category selection → free text search → select transaction → clone.

**Independent Test**: Tap Clone → "Search More" → select "Incomes" → select "Salary" → "Show All" → select a transaction → verify clone saved with today's date.

### Implementation for User Story 2

- [ ] T008 [US2] Implement `CloneSearchMore` handler in `internal/client/clone.go`: show type selection keyboard with 3 buttons: "💸 Expenses" (`clone.search.type.expense`), "💰 Incomes" (`clone.search.type.income`), "All" (`clone.search.type.all`), plus cancel button.
- [ ] T009 [US2] Implement `CloneSearchTypeSelected` handler in `internal/client/clone.go`: parse type from callback data, store type context in session body, set state to `StateSelectingCloneSearchCategory`, show category selection keyboard. For "expense" show expense categories, for "income" show income categories, for "all" show all categories. Include "All Categories" option. Use same category keyboard layout as `showDeleteSearchCategorySelection` in `delete.go` but with `clone.search.category.{CAT}` callback data.
- [ ] T010 [US2] Implement `CloneSearchCategorySelected` handler in `internal/client/clone.go`: parse category from callback data (`clone.search.category.{CAT}`), update session state to `StateEnteringCloneSearchQuery` with category in `Session.Body`, prompt user for search text with "Show All" (`clone.search.showall`) and "Cancel" (`clone.search.cancel`) buttons.
- [ ] T011 [US2] Implement `CloneSearchQueryEntered` handler in `internal/client/clone.go`: get search query from message text, get category from session body, reset state to `StateNormal`, call `showCloneSearchResults` with category, query, offset=0.
- [ ] T012 [US2] Implement `showCloneSearchResults` in `internal/client/clone.go`: use `Repositories.Transactions.SearchUserTransactions` to fetch results, format paginated list (same style as `formatDeleteSearchResults`), create keyboard with numbered selection buttons (`clone.search.select.{ID}`), pagination (`clone.search.page.{CAT}.{OFFSET}.{QUERY}`), page indicator (`clone.search.noop`), "New Search" (`clone.search.new`), and "Home" (`clone.search.home`) buttons. When user selects a transaction, reuse `CloneTransactionSelected` logic (clone + save + edit UI).
- [ ] T013 [US2] Implement remaining wizard handlers in `internal/client/clone.go`: `CloneSearchTransactionSelected` (same clone logic as `CloneTransactionSelected`), `CloneSearchResultsPage` (parse `clone.search.page.{CAT}.{OFFSET}.{QUERY}`), `CloneSearchShowAll` (search with `%` wildcard), `CloneSearchNoop`, `CloneSearchNew` (restart clone flow), `CloneSearchCancel` (reset state, cancel), `CloneSearchHome` (send home keyboard).
- [ ] T014 [US2] Add clone search state handling in `FreeTextRouter` in `internal/client/router.go`: add `StateEnteringCloneSearchQuery` case that calls `CloneSearchQueryEntered`.
- [ ] T015 [US2] Register wizard search handlers in `internal/client/setup.go`: add handlers for `clone.search.type.`, `clone.search.category.`, `clone.search.page.`, `clone.search.select.`, `clone.search.showall`, `clone.search.noop`, `clone.search.new`, `clone.search.cancel`, `clone.search.home`.

**Checkpoint**: User Story 2 complete — wizard search flow works. Users can find any transaction via type/category/text filtering and clone it.

---

## Phase 4: User Stories 3 & 4 — /clone Command + LLM Intent (Priority: P3)

**Goal**: `/clone` command and natural language "clone" intent both route to the clone flow.

**Independent Test**: Type `/clone` → verify recent expenses list. Type "duplicate a transaction" → verify clone flow starts.

### Implementation for User Stories 3 & 4

- [ ] T016 [P] [US3] Register `/clone` command in `internal/client/setup.go`: add `handlers.NewCommand("clone", c.CloneTransactions)`.
- [ ] T017 [P] [US4] Add clone intent routing in `internal/client/router.go`: add `case ai.IntentClone:` to `classifyAndRouteIntent` switch that calls `c.CloneTransactions(b, ctx)`.

**Checkpoint**: All user stories complete — clone accessible via home button, /clone command, and natural language.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Edge case handling and verification

- [ ] T018 Handle deleted source transaction in `CloneTransactionSelected` and `CloneSearchTransactionSelected` in `internal/client/clone.go`: if `GetByID` returns error, show "Transaction no longer exists" with home button
- [ ] T019 Handle empty state in `showRecentExpensesForClone` in `internal/client/clone.go`: if user has 0 expenses, show "No transactions to clone. Add your first transaction!" with Add Income/Expense and Home buttons
- [ ] T020 Verify `go build ./...` compiles successfully and manually test full flow per quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **US1 (Phase 2)**: Depends on Phase 1 (needs state constants)
- **US2 (Phase 3)**: Depends on Phase 2 (reuses clone save logic from US1)
- **US3+4 (Phase 4)**: Depends on Phase 1 (needs intent constant); can run in parallel with US1/US2 since it only adds routing
- **Polish (Phase 5)**: Depends on all phases complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends on Phase 1 only — independently testable
- **User Story 2 (P2)**: Depends on US1's `CloneTransactionSelected` for clone+save logic
- **User Story 3 (P3)**: No dependency on other stories — just registers command to existing handler
- **User Story 4 (P3)**: No dependency on other stories — just adds intent routing

### Within Each User Story

- Handlers before registration (write code, then wire it up)
- Entry handlers before detail handlers
- Home keyboard change (T006) can run in parallel with handler implementation

### Parallel Opportunities

- T001 and T002 can run in parallel (different files)
- T003, T004, T005, T006 can run in parallel (T003-T005 in clone.go, T006 in send.go — but T003-T005 are same file so sequential)
- T016 and T017 can run in parallel (different files)
- US3 and US4 can run in parallel with US2

---

## Parallel Example: Phase 1

```
Task T001: Add state constants in internal/model/users.go
Task T002: Add IntentClone + prompt update in internal/ai/deepseek.go + internal/ai/prompt.go
→ Both can execute simultaneously (different files)
```

## Parallel Example: Phase 4

```
Task T016: Register /clone command in internal/client/setup.go
Task T017: Add clone intent routing in internal/client/router.go
→ Both can execute simultaneously (different files)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001, T002)
2. Complete Phase 2: US1 (T003–T007)
3. **STOP and VALIDATE**: Test quick clone from home keyboard
4. Deploy if ready — users get 80% of the value from just this

### Incremental Delivery

1. Phase 1 → Setup complete
2. Phase 2 → US1: Quick clone works (MVP!)
3. Phase 3 → US2: Wizard search works (full feature)
4. Phase 4 → US3+4: Command + NLP (polish)
5. Phase 5 → Edge cases handled

---

## Notes

- All clone handlers go in one new file: `internal/client/clone.go`
- Clone saves immediately (same as add transaction) — no temporary state needed
- Post-clone edit reuses existing `StateEditingNewTransaction` + `transactions.edit.*` handlers
- No database migration required
- Callback data must stay under 64 bytes (Telegram limit) — verified all patterns fit
