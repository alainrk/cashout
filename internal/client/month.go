package client

import (
	"fmt"
	"happypoor/internal/model"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// MonthRecap returns to the user the breakdown and the total for the expenses and income of the current month.
func (c *Client) MonthRecap(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.LastMessage = ctx.Message.Text

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	// Get current month
	month := time.Now().Month()

	res, err := c.Repositories.Transactions.GetMonthlyTotalsCurrentYear(user.TgID)
	if err != nil {
		return err
	}

	t, ok := res[int(month)]
	if !ok {
		ctx.EffectiveMessage.Reply(b, "No transactions for this month", &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	var msg string
	var monthtot float64

	msg += fmt.Sprintf("🗓 %s\n", time.Month(month).String())

	if ex, ok := t[model.TypeExpense]; ok {
		msg += fmt.Sprintf("-%.2f€\n", ex)
		monthtot -= ex
	}

	if in, ok := t[model.TypeIncome]; ok {
		msg += fmt.Sprintf("+%.2f€\n", in)
		monthtot += in
	}

	msg += fmt.Sprintf("Total: %.2f€\n\n", monthtot)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
