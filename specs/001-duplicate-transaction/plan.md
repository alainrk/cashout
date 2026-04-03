# Implementation Plan: Duplicate Transaction

**Branch**: `001-duplicate-transaction` | **Date**: 2026-04-03 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-duplicate-transaction/spec.md`

## Summary

Add a "Clone" flow to the Telegram bot allowing users to duplicate existing transactions with today's date. Two paths: quick-clone from recent expenses (3 taps), or wizard search with type/category/text filtering. After saving, reuse existing edit-in-place UI. Adds `/clone` command, home keyboard button, and LLM intent.

## Technical Context

**Language/Version**: Go 1.24  
**Primary Dependencies**: gotgbot/v2 (Telegram bot framework), GORM (ORM), OpenAI-compatible LLM  
**Storage**: PostgreSQL with JSONB session column  
**Testing**: Manual Telegram bot testing (no test framework currently in use)  
**Target Platform**: Linux server (Docker)  
**Project Type**: Telegram bot + web dashboard  
**Constraints**: Telegram inline keyboard callback data max 64 bytes  
**Scale/Scope**: Single-user personal finance bot

## Constitution Check

*GATE: Constitution is a template (not project-specific). No violations to check.*

No project-specific constitution has been defined. Proceeding with standard engineering practices.

## Project Structure

### Documentation (this feature)

```text
specs/001-duplicate-transaction/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── telegram-callbacks.md
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
internal/
├── client/
│   ├── clone.go         # NEW: All clone flow handlers
│   ├── setup.go         # MODIFY: Register clone handlers
│   ├── send.go          # MODIFY: Add Clone to home keyboard
│   └── router.go        # MODIFY: Handle clone states + intent
├── model/
│   └── users.go         # MODIFY: Add clone state constants
└── ai/
    ├── prompt.go        # MODIFY: Add clone intent to classifier
    └── deepseek.go      # MODIFY: Add IntentClone constant
```

**Structure Decision**: Single project structure (already established). One new file (`clone.go`) following the pattern of `edit.go` and `delete.go`. All other changes are modifications to existing files.
