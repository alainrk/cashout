package client

import (
	"encoding/json"
	"fmt"
	"happypoor/internal/model"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Start introduces the bot.
func (c *Client) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Welcome to HappyPoor, %s!", user.Name)
	c.SendTransactionKeyboard(b, ctx, msg)

	return nil
}

func (c *Client) AddIncomeIntent(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.addTransactionIntent(b, ctx, model.TypeIncome)
}

func (c *Client) AddExpenseIntent(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.addTransactionIntent(b, ctx, model.TypeExpense)
}

func (c *Client) addTransactionIntent(b *gotgbot.Bot, ctx *ext.Context, transactionType model.TransactionType) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	if user.Session.Iterations == 0 {
		return c.Start(b, ctx)
	}

	user.Session.LastCommand = model.CommandAddExpenseIntent
	if transactionType == model.TypeIncome {
		user.Session.LastCommand = model.CommandAddIncomeIntent
	}
	user.Session.State = model.StateWaiting
	user.Session.LastMessage = ctx.Message.Text

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	ctx.EffectiveMessage.Reply(b, "Sure, just tell me category, amount and decription.", &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			Keyboard: [][]gotgbot.KeyboardButton{
				{
					{
						Text: "Cancel",
					},
				},
			},
			IsPersistent:    false,
			OneTimeKeyboard: true,
			ResizeKeyboard:  true,
		},
	})

	return nil
}

func (c *Client) AddTransaction(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	if user.Session.Iterations == 0 {
		return c.Start(b, ctx)
	}

	var transactionType model.TransactionType

	switch user.Session.LastCommand {
	case model.CommandAddIncomeIntent:
		transactionType = model.TypeIncome
	case model.CommandAddExpenseIntent:
		transactionType = model.TypeExpense
	default:
		// answer the user that they should chose a valid command first and send the keyboard
		c.SendTransactionKeyboard(b, ctx, "Add a transaction")
		return nil
	}

	transaction, err := c.LLM.ExtractTransaction(ctx.Message.Text, transactionType)
	if err != nil {
		msg := fmt.Sprintf("I'm sorry, I can't understand your transaction '%s', %s!", user.Session.LastMessage, user.Name)
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return err
	}

	// Store the transaction in the session
	user.Session.State = model.StateWaiting
	user.Session.LastCommand = model.CommandAddTransaction
	user.Session.LastMessage = ctx.Message.Text
	s, err := json.Marshal(transaction)
	user.Session.Body = string(s)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg := fmt.Sprintf("%s (€ %.2f), %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description)
	c.SendConfirmKeyboard(b, ctx, msg, []string{"Edit"})

	return nil
}

// Edit edits the transaction previously inserted, basically returns to the same add Income/Expense state, cleaning up the session.
func (c *Client) AmendTransaction(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	var transaction model.Transaction
	err = json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	user.Session.State = model.StateEdit
	user.Session.LastMessage = ctx.Message.Text

	// Trick to put the session state back to the correct add transaction
	user.Session.LastCommand = model.CommandAddExpenseIntent
	if transaction.Type == model.TypeIncome {
		user.Session.LastCommand = model.CommandAddIncomeIntent
	}
	// Clean up the previously inserted session
	user.Session.Body = ""

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	ctx.EffectiveMessage.Reply(b, "Sure, to better edit the transaction, specify the category, amount, and description as best you can.", &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			Keyboard: [][]gotgbot.KeyboardButton{
				{
					{
						Text: "Cancel",
					},
				},
			},
			IsPersistent:    false,
			OneTimeKeyboard: true,
			ResizeKeyboard:  true,
		},
	})

	return nil
}

// Confirm confirms the previous action after the user been prompted.
func (c *Client) Confirm(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	var transaction model.Transaction
	err = json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	user.Session.State = model.StateNormal
	user.Session.LastCommand = model.CommandConfirm
	user.Session.LastMessage = ctx.Message.Text
	user.Session.Body = ""

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	transaction.TgID = user.TgID
	transaction.Currency = model.CurrencyEUR

	c.Repositories.Transactions.Add(transaction)

	c.SendTransactionKeyboard(b, ctx, "Your transaction has been saved!")

	return nil
}

// Cancel returns to normal state.
func (c *Client) Cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.LastCommand = model.CommandCancel
	user.Session.LastMessage = ctx.Message.Text

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	c.SendTransactionKeyboard(b, ctx, "Add a transaction")

	return nil
}

// MonthRecap returns to the user the breakdown and the total for the expenses and income of the current month.
func (c *Client) MonthRecap(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.LastCommand = model.CommandMonthRecap
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

	if cmonth, ok := res[int(month)]; ok {
		fmt.Printf("%+v\n", cmonth)
		var msg string
		var total float64
		if ex, ok := cmonth[model.TypeExpense]; ok {
			msg += fmt.Sprintf("Expenses: € %.2f \n", ex)
			total -= ex
		}
		if in, ok := cmonth[model.TypeIncome]; ok {
			msg += fmt.Sprintf("Income: € %.2f\n", in)
			total += in
		}

		msg += fmt.Sprintf("Total: € %.2f", total)
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
	}

	return nil
}

func (c *Client) YearRecap(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.LastCommand = model.CommandMonthRecap
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

	// TODO:
	// - Create and put everything into an array
	// - Take month and expense/income from the map in O(1)
	// - Loop on month 0->11 there so months gets sorted and take ex/in from inner map
	for m, t := range res {
		var monthtot float64

		msg += fmt.Sprintf("Month: %s\n", time.Month(m).String())

		if ex, ok := t[model.TypeExpense]; ok {
			msg += fmt.Sprintf("Expenses: € %.2f \n", ex)
			monthtot -= ex
			yearex += ex
		}

		if in, ok := t[model.TypeIncome]; ok {
			msg += fmt.Sprintf("Income: € %.2f\n", in)
			monthtot += in
			yearin += in
		}

		yeartot += monthtot

		msg += fmt.Sprintf("Total: € %.2f", monthtot)

		msg += "\n---\n"
	}

	msg += fmt.Sprintf("Year Expenses: € %.2f\n", yearex)
	msg += fmt.Sprintf("Year Incomes: € %.2f\n", yearin)
	msg += fmt.Sprintf("Year Total: € %.2f", yeartot)

	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
