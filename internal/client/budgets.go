package client

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"cashout/internal/model"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gorm.io/gorm"
)

// budgetThresholds are the percentages at which alerts fire (each at most once per month).
var budgetThresholds = []int16{80, 100}

// BudgetProgress is the result of evaluating a user's spending against their budget.
type BudgetProgress struct {
	Limit     float64
	Spent     float64
	Pct       int
	NewAlerts []int16 // subset of budgetThresholds that just crossed on this insert
}

// EvaluateAfterExpenseInsert computes budget progress and fires any newly-crossed
// threshold alerts for the calendar month of the inserted transaction.
// Returns nil if the user has no budget set or the transaction is not an Expense.
//
// Alert semantics:
//   - "80% approaching" is a one-shot transition warning (dedup'd in DB per month).
//   - ">=100% over budget" is an ongoing condition; surfaced on every expense
//     while the user remains over, so they don't sleepwalk past the limit.
func (c *Client) EvaluateAfterExpenseInsert(tx model.Transaction) (*BudgetProgress, error) {
	if tx.Type != model.TypeExpense {
		return nil, nil
	}

	budget, err := c.Repositories.Budgets.Get(tx.TgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	year := tx.Date.Year()
	month := int(tx.Date.Month())

	spentAfter, err := c.Repositories.Budgets.TotalExpensesForMonth(tx.TgID, year, month)
	if err != nil {
		return nil, err
	}
	spentBefore := spentAfter - tx.Amount
	if spentBefore < 0 {
		spentBefore = 0
	}

	yearMonth := fmt.Sprintf("%04d-%02d", year, month)
	progress := &BudgetProgress{
		Limit: budget.Amount,
		Spent: spentAfter,
		Pct:   int(math.Floor(spentAfter / budget.Amount * 100)),
	}

	// 80% approaching: one-shot transition warning (dedup'd in DB).
	// Only surface if the user is in the [80,100) band — once they cross 100,
	// the over-budget message below supersedes it.
	if spentAfter < budget.Amount {
		cutoff80 := budget.Amount * 0.80
		if spentBefore < cutoff80 && spentAfter >= cutoff80 {
			fired, err := c.Repositories.Budgets.TryMarkAlertFired(tx.TgID, yearMonth, 80)
			if err != nil {
				c.Logger.Warnf("failed to mark budget alert fired: %v", err)
			} else if fired {
				progress.NewAlerts = append(progress.NewAlerts, 80)
			}
		}
	}

	// Over budget: ongoing condition — always surface while pct >= 100.
	// Mark fired on the first crossing for record-keeping, but don't gate display.
	if spentAfter >= budget.Amount {
		if spentBefore < budget.Amount {
			if _, err := c.Repositories.Budgets.TryMarkAlertFired(tx.TgID, yearMonth, 100); err != nil {
				c.Logger.Warnf("failed to mark budget alert fired: %v", err)
			}
		}
		progress.NewAlerts = append(progress.NewAlerts, 100)
	}

	return progress, nil
}

// BudgetSuffixForTx returns the FormatBudgetSuffix string for the month of the
// given transaction (the relevant month for any edit/delete impact on a tx),
// swallowing internal errors with a log line. Empty string if no budget.
func (c *Client) BudgetSuffixForTx(tx model.Transaction) string {
	progress, err := c.BudgetStatusForMonth(tx.TgID, tx.Date.Year(), int(tx.Date.Month()))
	if err != nil {
		c.Logger.Warnf("budget status lookup failed: %v", err)
		return ""
	}
	return FormatBudgetSuffix(progress)
}

// BudgetStatusForMonth returns the current budget status for a given month
// without firing any alerts. Used after edit/delete confirmations to keep the
// over-budget warning visible whenever the user is still over.
// Returns nil if no budget is set.
func (c *Client) BudgetStatusForMonth(tgID int64, year, month int) (*BudgetProgress, error) {
	budget, err := c.Repositories.Budgets.Get(tgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	spent, err := c.Repositories.Budgets.TotalExpensesForMonth(tgID, year, month)
	if err != nil {
		return nil, err
	}

	progress := &BudgetProgress{
		Limit: budget.Amount,
		Spent: spent,
		Pct:   int(math.Floor(spent / budget.Amount * 100)),
	}
	// Surface the over-budget nag on every refresh while user remains over.
	if spent >= budget.Amount {
		progress.NewAlerts = []int16{100}
	}
	return progress, nil
}

// FormatBudgetSuffix builds the trailing message lines appended to a transaction
// confirmation when a budget exists.
func FormatBudgetSuffix(p *BudgetProgress) string {
	if p == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n\n📊 Budget: %.2f / %.2f € (%d%%)", p.Spent, p.Limit, p.Pct))
	for _, t := range p.NewAlerts {
		switch t {
		case 80:
			b.WriteString("\n⚠️ Approaching monthly budget (80% used)")
		case 100:
			over := p.Spent - p.Limit
			b.WriteString(fmt.Sprintf("\n🚨 Over budget by %.2f €", over))
		}
	}
	return b.String()
}

// ShowBudget renders the current budget and progress for the current month.
func (c *Client) ShowBudget(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.Body = ""
	if err := c.Repositories.Users.Update(&user); err != nil {
		return fmt.Errorf("failed to update user state: %w", err)
	}

	budget, err := c.Repositories.Budgets.Get(user.TgID)
	keyboard := [][]gotgbot.InlineKeyboardButton{
		{{Text: "💰 Set / Update Budget", CallbackData: "budget.setprompt"}},
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			text := "📊 <b>Monthly Budget</b>\n\nYou haven't set a monthly budget yet.\n\nUse <code>/budget set &lt;amount&gt;</code> or tap below."
			keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{{Text: "🏠 Home", CallbackData: "transactions.home"}})
			return SendMessage(ctx, b, text, keyboard)
		}
		return fmt.Errorf("failed to get budget: %w", err)
	}

	now := time.Now()
	spent, err := c.Repositories.Budgets.TotalExpensesForMonth(user.TgID, now.Year(), int(now.Month()))
	if err != nil {
		return fmt.Errorf("failed to compute month total: %w", err)
	}
	pct := int(math.Floor(spent / budget.Amount * 100))

	indicator := "✅"
	if pct >= 100 {
		indicator = "🚨"
	} else if pct >= 80 {
		indicator = "⚠️"
	}

	text := fmt.Sprintf(
		"📊 <b>Monthly Budget</b>\n\n%s %.2f / %.2f € (%d%%)\nMonth: %s",
		indicator, spent, budget.Amount, pct, now.Format("January 2006"),
	)

	keyboard = append(keyboard,
		[]gotgbot.InlineKeyboardButton{{Text: "🗑 Remove Budget", CallbackData: "budget.delete"}},
		[]gotgbot.InlineKeyboardButton{{Text: "🏠 Home", CallbackData: "transactions.home"}},
	)

	return SendMessage(ctx, b, text, keyboard)
}

// BudgetCommand handles /budget and its subcommands: "set <amount>", "delete", or "" (show).
func (c *Client) BudgetCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.Message == nil {
		return c.ShowBudget(b, ctx)
	}

	parts := strings.Fields(ctx.Message.Text)
	// parts[0] == "/budget"
	if len(parts) < 2 {
		return c.ShowBudget(b, ctx)
	}

	switch strings.ToLower(parts[1]) {
	case "set":
		if len(parts) < 3 {
			_, err := b.SendMessage(ctx.EffectiveSender.ChatId,
				"Usage: <code>/budget set &lt;amount&gt;</code>",
				&gotgbot.SendMessageOpts{ParseMode: "HTML"})
			return err
		}
		return c.budgetSet(b, ctx, parts[2])
	case "delete", "remove", "clear":
		return c.budgetDelete(b, ctx)
	default:
		// Treat "/budget 400" as "/budget set 400"
		return c.budgetSet(b, ctx, parts[1])
	}
}

// BudgetSetPrompt enters a wizard waiting for the amount.
func (c *Client) BudgetSetPrompt(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateBudgetSetWaitAmount
	user.Session.Body = ""
	if err := c.Repositories.Users.Update(&user); err != nil {
		return fmt.Errorf("failed to update user state: %w", err)
	}

	keyboard := [][]gotgbot.InlineKeyboardButton{
		{{Text: "Cancel", CallbackData: "budget.cancel"}},
	}
	return SendMessage(ctx, b, "Enter your monthly budget amount in € (e.g. <code>1500</code>):", keyboard)
}

// BudgetSetFromMessage receives the amount typed by the user after BudgetSetPrompt.
func (c *Client) BudgetSetFromMessage(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	if err := c.Repositories.Users.Update(&user); err != nil {
		return fmt.Errorf("failed to reset user state: %w", err)
	}
	return c.budgetSet(b, ctx, strings.TrimSpace(ctx.Message.Text))
}

// BudgetDeleteCallback handles the inline-keyboard delete button.
func (c *Client) BudgetDeleteCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.budgetDelete(b, ctx)
}

// BudgetCancel resets state and returns to home.
func (c *Client) BudgetCancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	if err := c.Repositories.Users.Update(&user); err != nil {
		return fmt.Errorf("failed to reset user state: %w", err)
	}
	return c.SendHomeKeyboard(b, ctx, "Operation cancelled.")
}

func (c *Client) budgetSet(b *gotgbot.Bot, ctx *ext.Context, amountStr string) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	amountStr = strings.ReplaceAll(amountStr, ",", ".")
	var amount float64
	if _, err := fmt.Sscanf(amountStr, "%f", &amount); err != nil || amount <= 0 {
		_, sendErr := b.SendMessage(ctx.EffectiveSender.ChatId,
			"Invalid amount. Please enter a positive number, e.g. <code>1500</code>.",
			&gotgbot.SendMessageOpts{ParseMode: "HTML"})
		return sendErr
	}

	budget := model.Budget{
		TgID:     user.TgID,
		Amount:   amount,
		Currency: model.CurrencyEUR,
	}
	if err := c.Repositories.Budgets.Upsert(&budget); err != nil {
		return fmt.Errorf("failed to upsert budget: %w", err)
	}

	now := time.Now()
	spent, err := c.Repositories.Budgets.TotalExpensesForMonth(user.TgID, now.Year(), int(now.Month()))
	if err != nil {
		return fmt.Errorf("failed to compute month total: %w", err)
	}
	pct := int(math.Floor(spent / amount * 100))

	text := fmt.Sprintf(
		"✅ Monthly budget set to <b>%.2f €</b>.\n\nThis month so far: %.2f € (%d%%).",
		amount, spent, pct,
	)
	return c.SendHomeKeyboard(b, ctx, text)
}

func (c *Client) budgetDelete(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	if err := c.Repositories.Budgets.Delete(user.TgID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.SendHomeKeyboard(b, ctx, "No budget to remove.")
		}
		return fmt.Errorf("failed to delete budget: %w", err)
	}
	return c.SendHomeKeyboard(b, ctx, "🗑 Budget removed.")
}
