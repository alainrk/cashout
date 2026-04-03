# Quickstart: Duplicate Transaction

## What This Feature Does

Adds a "Clone" flow to the Telegram bot that lets users duplicate an existing transaction with today's date. The flow provides two paths:
1. **Quick clone**: Shows 10 most recent expenses → tap to clone → edit if needed → done
2. **Wizard search**: Filter by type (income/expense) → category → text search → select → clone

After cloning, the transaction is saved immediately and the user gets the same edit UI as after adding a new transaction.

## Files to Create

| File | Purpose |
|------|---------|
| `internal/client/clone.go` | All clone flow handlers (entry, recent list, wizard, selection, save) |

## Files to Modify

| File | Change |
|------|--------|
| `internal/model/users.go` | Add 3 new `StateType` constants for clone flow |
| `internal/client/setup.go` | Register `/clone` command + all `clone.*` callback handlers |
| `internal/client/send.go` | Add "Clone" button to home keyboard (row 2) |
| `internal/client/router.go` | Add clone search query state handling in `FreeTextRouter` + clone intent routing |
| `internal/ai/prompt.go` | Add `clone` intent to `LLMIntentClassificationPromptTemplate` |
| `internal/ai/deepseek.go` | Add `IntentClone` constant |

## How to Test

```bash
# Build and run
make run-server

# In Telegram:
# 1. Send /start → verify "Clone" button appears in home keyboard
# 2. Tap "Clone" → verify recent expenses shown
# 3. Select a transaction → verify clone saved with today's date
# 4. Edit a field → verify it updates
# 5. Send /clone → verify same flow as tapping Clone button
# 6. Tap "Search More" → verify wizard works (type → category → search)
# 7. Type "clone" as free text → verify LLM routes to clone flow
```

## Architecture Notes

- The clone flow follows the exact same patterns as edit/delete flows
- Clone saves immediately then shows edit UI (same as add transaction)
- No database migration needed — all changes are application-level
- Callback data format: `clone.{action}.{value}` (consistent with edit/delete)
