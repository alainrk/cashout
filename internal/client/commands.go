package client

import (
	"encoding/json"
	"fmt"
	"happypoor/internal/model"
	"happypoor/internal/utils"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Start introduces the bot.
func (c *Client) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Welcome to HappyPoor, %s!", user.Name)
	c.SendAddTransactionKeyboard(b, ctx, msg)

	return nil
}

func (c *Client) AddTransactionIntent(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	if user.Session.Iterations == 0 {
		return c.Start(b, ctx)
	}

	query := ctx.CallbackQuery
	msg := query.Message

	action := strings.Split(query.Data, ".")[2]

	switch action {
	case "income":
		user.Session.State = model.StateInsertingIncome
	case "expense":
		user.Session.State = model.StateInsertingExpense
	default:
		c.SendAddTransactionKeyboard(b, ctx, "Invalid action, add a transaction.")
		return nil
	}

	user.Session.LastMessage = action

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg.EditText(b, fmt.Sprintf("Sure, to add a new %s, just tell me category, amount and decription.", action), nil)
	msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Cancel", CallbackData: "transactions.cancel"}},
			},
		},
	})

	return ext.ContinueGroups
}

func (c *Client) FreeTextRouter(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	if user.Session.Iterations == 0 {
		return c.Start(b, ctx)
	}

	if user.Session.State == model.StateInsertingIncome || user.Session.State == model.StateInsertingExpense {
		return c.addTransaction(b, ctx, user)
	}

	if user.Session.State == model.StateEditingTransaction {
		return c.editTransaction(b, ctx, user)
	}

	return fmt.Errorf("invalid top-level state")
}

func (c *Client) addTransaction(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transactionType model.TransactionType

	switch user.Session.State {
	case model.StateInsertingIncome:
		transactionType = model.TypeIncome
	case model.StateInsertingExpense:
		transactionType = model.TypeExpense
	default:
		c.SendAddTransactionKeyboard(b, ctx, "Invalid action, add a transaction.")
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

	if transaction.Amount == 0 {
		msg := fmt.Sprintf("I'm sorry, I can't understand your transaction '%s', %s!", user.Session.LastMessage, user.Name)
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	// Store the transaction in the session
	user.Session.State = model.StateNormal
	user.Session.LastMessage = ctx.Message.Text
	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg := fmt.Sprintf("%s (€ %.2f), %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description)
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, msg, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("failed to send confirm message: %v\n", err)
		return err
	}

	return nil
}

func (c *Client) EditTransactionIntent(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	c.CleanupInlineKeyboard(b, ctx)

	user.Session.State = model.StateEditingTransaction
	user.Session.LastMessage = "edit"

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	b.SendMessage(ctx.EffectiveSender.ChatId, "Add your date (e.g. dd mm, dd-mm, dd-mm-yyyy).", &gotgbot.SendMessageOpts{})

	return nil
}

// EditTransaction edits a transaction.
func (c *Client) editTransaction(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	// TODO: Handle this with LLM to support different formats automatically and also yesterday, 2 days ago etc.

	// Get date from DD-MM-YYYY to date
	date, err := utils.ParseDate(ctx.Message.Text)
	if err != nil {
		// TODO: Handle invalid date
		fmt.Printf("failed to parse date: %v\n", err)
		return err
	}

	// TODO: Reject future date

	transaction.Date = date

	user.Session.LastMessage = ctx.Message.Text
	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	m := fmt.Sprintf("%s (€ %.2f), %s on %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
				},
			},
		},
	})

	return nil
}

// Confirm confirms the previous action after the user been prompted.
func (c *Client) Confirm(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	var transaction model.Transaction
	err = json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	user.Session.State = model.StateNormal
	user.Session.LastMessage = "confirm"
	user.Session.Body = ""

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	transaction.TgID = user.TgID
	transaction.Currency = model.CurrencyEUR

	err = c.Repositories.Transactions.Add(transaction)
	if err != nil {
		// TODO: Handle send failure message to the user
		fmt.Printf("failed to add transaction: %v\n", err)
		return fmt.Errorf("failed to add transaction: %w", err)
	}

	// Remove the keyboard from the previous message
	query := ctx.CallbackQuery
	msg := query.Message
	msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{},
		},
	})

	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}
	c.SendAddTransactionKeyboard(b, ctx, fmt.Sprintf("%s Your transaction has been saved!", emoji))

	return nil
}

// Cancel returns to normal state.
func (c *Client) Cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery
	msg := query.Message

	user.Session.State = model.StateNormal
	user.Session.LastMessage = "cancel"

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{},
	})

	c.SendAddTransactionKeyboard(b, ctx, "Add a transaction")

	return ext.EndGroups
}

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
