# Telegram Bot Callback Contracts: Clone Flow

## New Callback Data Patterns

### Home Keyboard Entry
| Callback Data | Handler | Description |
|---------------|---------|-------------|
| `home.clone` | `CloneTransactions` | Entry point from home keyboard |

### Recent Transactions List
| Callback Data | Handler | Description |
|---------------|---------|-------------|
| `clone.select.{ID}` | `CloneTransactionSelected` | User selected a transaction to clone |
| `clone.page.{OFFSET}` | `CloneTransactionPage` | Pagination for recent list |
| `clone.noop` | `CloneNoop` | No-op for page indicator |
| `clone.searchmore` | `CloneSearchMore` | Enter wizard search |

### Wizard Search
| Callback Data | Handler | Description |
|---------------|---------|-------------|
| `clone.search.type.{expense\|income\|all}` | `CloneSearchTypeSelected` | Type filter selection |
| `clone.search.category.{CAT}` | `CloneSearchCategorySelected` | Category filter selection |
| `clone.search.page.{CAT}.{OFFSET}.{QUERY}` | `CloneSearchResultsPage` | Paginated search results |
| `clone.search.select.{ID}` | `CloneSearchTransactionSelected` | Select from search results |
| `clone.search.showall` | `CloneSearchShowAll` | Show all in selected category |
| `clone.search.noop` | `CloneSearchNoop` | No-op for page indicator |
| `clone.search.new` | `CloneSearchNew` | Start new search |
| `clone.search.cancel` | `CloneSearchCancel` | Cancel search, reset state |
| `clone.search.home` | `CloneSearchHome` | Return to home |

### Post-Clone (Reuses Existing)
After clone saves, the flow enters `StateEditingNewTransaction` and reuses existing callbacks:
- `transactions.edit.{field}` — Edit cloned transaction field
- `transactions.delete.{ID}` — Delete the clone
- `transactions.home` — Return to home

## New Bot Command
| Command | Handler | Description |
|---------|---------|-------------|
| `/clone` | `CloneTransactions` | Start clone flow (same as home.clone) |

## New LLM Intent
| Intent | Keywords | Routes To |
|--------|----------|-----------|
| `clone` | clone, duplicate, repeat, copy, same again | `CloneTransactions` |
