# Research: Duplicate Transaction

## Decision: Clone Flow Entry Point

- **Decision**: Show 10 most recent expenses as default view, with "Search More" to enter the wizard.
- **Rationale**: Matches the user's stated preference ("first of all it provides the latest expenses"). Expenses are cloned far more frequently than incomes (recurring daily spending). Showing recent expenses immediately minimizes taps for the common case.
- **Alternatives considered**: Showing all recent transactions (both types) — rejected because expenses dominate usage and mixing types adds clutter. Starting with wizard directly — rejected because it adds unnecessary steps for the common case.

## Decision: Wizard Flow Design

- **Decision**: Three-step wizard: Type selection → Category selection → Free text search / Show All. Reuse existing search patterns from edit/delete flows.
- **Rationale**: Consistent with existing edit/delete search UX. Users already know the category → search pattern. Adding type (Income/Expense) as a first step is needed because clone spans both types unlike edit/delete which show all.
- **Alternatives considered**: Single flat search across all transactions — rejected because too many results without filtering. Reuse edit search directly — almost works but needs the type filter prepended.

## Decision: Clone Confirmation UI

- **Decision**: Reuse the same "saved transaction" confirmation pattern from `addTransaction()` — show transaction summary with edit buttons for each field, plus Confirm and Cancel.
- **Rationale**: The add transaction flow already has the exact UX the user wants: "the user must be able to edit it right away, for every field". The only difference is the clone shows a pre-filled confirmation rather than saving immediately.
- **Alternatives considered**: Save immediately then show edit options (like current add flow) — actually this IS the simpler approach and matches existing behavior. Clone saves immediately with today's date, then shows the same edit UI as add transaction. This means less new code and consistent UX.

## Decision: Save Strategy

- **Decision**: Save the cloned transaction immediately (with today's date), then show edit options. Same pattern as `addTransaction()`.
- **Rationale**: The existing add flow saves first, then lets users edit. This is simpler (no temporary state to manage) and consistent. If the user doesn't want to edit, they're done. If they do, the edit-in-place flow already works.
- **Alternatives considered**: Hold in memory until confirm — rejected because it requires new state management and breaks consistency with the add flow.

## Decision: New Session States

- **Decision**: Add minimal new states: `StateSelectingCloneTransaction` (for recent list + wizard), `StateSelectingCloneSearchCategory` and `StateEnteringCloneSearchQuery` (for wizard search).
- **Rationale**: Follows the exact same state pattern as edit/delete search flows. The `FreeTextRouter` needs to handle the search query state for clone.
- **Alternatives considered**: Reusing edit states — rejected because it would create ambiguity in the router.

## Decision: Home Keyboard Layout

- **Decision**: Add "Clone" button alongside Edit and Delete in row 2, making it a 3-button row: `✏️ Edit | 📋 Clone | 🗑 Delete`.
- **Rationale**: Clone is a transaction operation like Edit and Delete, so it belongs in the same row. The user specifically asked for it in the "quick keyboard in the home".
- **Alternatives considered**: New row for Clone — rejected because it makes the keyboard taller and Clone logically groups with Edit/Delete.

## Decision: LLM Intent

- **Decision**: Add `clone` intent to the LLM classifier prompt with keywords: clone, duplicate, repeat, copy, same again.
- **Rationale**: Consistent with how all other flows are accessible via natural language. Low effort — just add to the prompt template and switch case.
- **Alternatives considered**: Not adding LLM intent — rejected because it breaks the pattern where every flow is reachable via natural language.
