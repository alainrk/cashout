# Data Model: Duplicate Transaction

## Existing Entities (No Changes)

### Transaction
No schema changes. Clone creates a new `Transaction` row with:
- New auto-incremented `ID`
- Same `TgID` (user ownership)
- `Date` set to `time.Now()` (today) — overridden from source
- Same `Type`, `Category`, `Amount`, `Currency`, `Description` — copied from source
- New `CreatedAt` / `UpdatedAt` timestamps (auto-generated)

### User
No schema changes to `User` struct.

## Modified Entities

### UserSession States

New state constants added to `model.StateType`:

| State | Value | Purpose |
|-------|-------|---------|
| `StateSelectingCloneTransaction` | `"selecting_clone_transaction"` | User is viewing recent expenses list or wizard results |
| `StateSelectingCloneSearchCategory` | `"selecting_clone_search_category"` | User is selecting category in clone wizard |
| `StateEnteringCloneSearchQuery` | `"entering_clone_search_query"` | User is typing free-text search in clone wizard |

**Session.Body usage during clone flow:**
- During `StateSelectingCloneSearchCategory`: empty
- During `StateEnteringCloneSearchQuery`: selected category (e.g., `"Grocery"` or `"all"`)
- After clone save: transaction ID stored (same as add flow, enters `StateEditingNewTransaction`)

## State Transitions

```
StateNormal
  ↓ (user taps "Clone" on home / types /clone / LLM intent)
StateSelectingCloneTransaction
  ↓ (user selects from recent list)
  → Save clone → StateEditingNewTransaction (reuses add flow's edit UI)
  
StateSelectingCloneTransaction
  ↓ (user taps "Search More")
StateSelectingCloneSearchCategory
  ↓ (user selects type filter, then category)
StateEnteringCloneSearchQuery
  ↓ (user types search or taps "Show All")
StateSelectingCloneTransaction (back to selection, now with search results)
  ↓ (user selects transaction)
  → Save clone → StateEditingNewTransaction (reuses add flow's edit UI)
```

## No Database Migration Required

All changes are application-level (new Go constants, new handlers). The `UserSession.State` is stored as a JSON string in the `session` JSONB column, so new state values work without schema changes.
