package client

import (
	"fmt"
	"happypoor/internal/model"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (c *Client) YearRecap(b *gotgbot.Bot, ctx *ext.Context) error {
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

	res, err := c.Repositories.Transactions.GetMonthlyTotalsCurrentYear(user.TgID)
	if err != nil {
		return err
	}

	var msg string
	var yeartot float64
	var yearex float64
	var yearin float64

	currMonth := time.Now().Month()

	for m := 1; m <= int(currMonth); m++ {
		t, ok := res[m]

		msg += fmt.Sprintf("🗓 %s\n", time.Month(m).String())
		if !ok {
			msg += "No entries\n\n"
			continue
		}

		var monthtot float64

		if ex, ok := t[model.TypeExpense]; ok {
			msg += fmt.Sprintf("-%.2f€\n", ex)
			monthtot -= ex
			yearex += ex
		}

		if in, ok := t[model.TypeIncome]; ok {
			msg += fmt.Sprintf("+%.2f€\n", in)
			monthtot += in
			yearin += in
		}

		yeartot += monthtot

		msg += fmt.Sprintf("Total: %.2f€\n\n", monthtot)
	}

	msg += "\n💰 Year to Date\n"

	if yearex > 0 {
		msg += fmt.Sprintf("-%.2f€\n", yearex)
	}
	if yearin > 0 {
		msg += fmt.Sprintf("+%.2f€\n", yearin)
	}
	msg += fmt.Sprintf("Total: %.2f€", yeartot)

	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
