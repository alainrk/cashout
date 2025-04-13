package client

import (
	"fmt"
	"happypoor/internal/model"

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
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	ctx.EffectiveMessage.Reply(b, "Add a transaction", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			Keyboard: [][]gotgbot.KeyboardButton{
				{
					{
						Text: "Income",
					},
					{
						Text: "Expense",
					},
				},
			},
			IsPersistent:   true,
			ResizeKeyboard: true,
		},
	})

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

	user.Session.LastCommand = model.CommandAddExpense
	if transactionType == model.TypeIncome {
		user.Session.LastCommand = model.CommandAddIncome
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
	case model.CommandAddIncome:
		transactionType = model.TypeIncome
	case model.CommandAddExpense:
		transactionType = model.TypeExpense
	default:
		// answer the user that they should chose a valid command first and send the keyboard
		ctx.EffectiveMessage.Reply(b, "Add a transaction", &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
				Keyboard: [][]gotgbot.KeyboardButton{
					{
						{
							Text: "Income",
						},
						{
							Text: "Expense",
						},
					},
				},
				IsPersistent:   true,
				ResizeKeyboard: true,
			},
		})
		return nil
	}

	user.Session.State = model.StateWaiting
	user.Session.LastMessage = ctx.Message.Text

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	transaction, err := c.LLM.ExtractTransaction(ctx.Message.Text, transactionType)
	if err != nil {
		msg := fmt.Sprintf("I'm sorry, I can't understand your transaction '%s', %s!", user.Session.LastMessage, user.Name)
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return err
	}

	// TODO: Save the transaction and send back a proper keyboard/command/text

	msg := fmt.Sprintf("%s (€ %.2f), %s", transaction.Category, transaction.Amount, transaction.Description)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

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

	ctx.EffectiveMessage.Reply(b, "Add a transaction", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			Keyboard: [][]gotgbot.KeyboardButton{
				{
					{
						Text: "Income",
					},
					{
						Text: "Expense",
					},
				},
			},
			IsPersistent:   true,
			ResizeKeyboard: true,
		},
	})

	return nil
}
