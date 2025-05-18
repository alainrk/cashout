package client

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (c *Client) AddTransactionIntent(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
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

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg.EditText(b, fmt.Sprintf("Sure! To add a new <b>%s</b>:\nTell me category, amount and description. You can also specify a date and change it later, today is default.\n\n<i>Examples:</i>\n<code>Irish Pub 3.4</code>\n<code>January salary 3k 10/01</code>", action), &gotgbot.EditMessageTextOpts{
		ParseMode: "HTML",
	})
	msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Cancel", CallbackData: "transactions.cancel"}},
			},
		},
	})

	return ext.ContinueGroups
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
		msg := "I'm sorry, I couldn't understand your transaction!"
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return err
	}

	if transaction.Amount == 0 {
		msg := "I'm sorry, I couldn't understand your transaction!"
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return nil
	}

	// Store the transaction in the session
	user.Session.State = model.StateWaitingConfirm
	s, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to stringify the body: %w", err)
	}
	user.Session.Body = string(s)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg := fmt.Sprintf("%s (€ %.2f), %s on %s. Confirm?", transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, msg, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit category",
						CallbackData: "transactions.edit.category",
					},
				},
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
				},
				{
					{
						Text:         "Edit amount",
						CallbackData: "transactions.edit.amount",
					},
				},
				{
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
				},
			},
		},
	})
	if err != nil {
		c.Logger.Errorln("failed to send confirm message", err)
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

	var transaction model.Transaction
	err = json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	query := ctx.CallbackQuery

	field := strings.Split(query.Data, ".")[2]

	var text string
	var opts *gotgbot.SendMessageOpts

	switch field {
	case "amount":
		user.Session.State = model.StateEditingTransactionAmount

		keyboard := [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Cancel",
					CallbackData: "transactions.cancel",
				},
			},
		}
		opts = &gotgbot.SendMessageOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard}}
		text = fmt.Sprintf("Enter a new amount for the transaction:\n\nCurrent: %.2f€ ", transaction.Amount)
	case "date":
		user.Session.State = model.StateEditingTransactionDate
		text = "Add your date (e.g. dd mm, dd-mm, dd-mm-yyyy)."
		keyboard := [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Cancel",
					CallbackData: "transactions.cancel",
				},
			},
		}
		opts = &gotgbot.SendMessageOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard}}
	case "category":
		user.Session.State = model.StateEditingTransactionCategory
		text = "Choose your category among the following ones."

		keyboard := [][]gotgbot.KeyboardButton{
			{{Text: "Cancel"}},
			{{Text: "Salary"}},
			{{Text: "OtherIncomes"}},
		}

		if transaction.Type == model.TypeExpense {
			keyboard = [][]gotgbot.KeyboardButton{
				{{Text: "Cancel"}},
				{{Text: "Car"}},
				{{Text: "Clothes"}},
				{{Text: "Grocery"}},
				{{Text: "House"}},
				{{Text: "Bills"}},
				{{Text: "Entertainment"}},
				{{Text: "Sport"}},
				{{Text: "EatingOut"}},
				{{Text: "Transport"}},
				{{Text: "Learning"}},
				{{Text: "Toiletry"}},
				{{Text: "Health"}},
				{{Text: "Tech"}},
				{{Text: "Gifts"}},
				{{Text: "Travel"}},
				{{Text: "OtherExpenses"}},
			}
		}

		opts = &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
				Keyboard:        keyboard,
				OneTimeKeyboard: true,
				IsPersistent:    false,
				ResizeKeyboard:  true,
			},
		}
	default:
		return fmt.Errorf("unknown field: %s", field)
	}

	c.CleanupKeyboard(b, ctx)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	b.SendMessage(ctx.EffectiveSender.ChatId, text, opts)

	return nil
}

func (c *Client) editTransactionDate(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	// Get date from DD-MM-YYYY to date
	date, err := utils.ParseDate(ctx.Message.Text)
	if err != nil {
		fmt.Printf("failed to parse date: %v\n", err)
		b.SendMessage(ctx.EffectiveSender.ChatId, "Invalid date, please try again.", nil)
		return err
	}

	if date.After(time.Now()) {
		b.SendMessage(ctx.EffectiveSender.ChatId, "I don't support future dates, please try again.", nil)
		return fmt.Errorf("invalid date: %s", ctx.Message.Text)
	}

	transaction.Date = date

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
						Text:         "Edit category",
						CallbackData: "transactions.edit.category",
					},
				},
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
				},
				{
					{
						Text:         "Edit amount",
						CallbackData: "transactions.edit.amount",
					},
				},
				{
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
				},
			},
		},
	})

	return nil
}

func (c *Client) editTransactionAmount(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	// Parse new amount from message
	amountStr := strings.TrimSpace(ctx.Message.Text)
	amountStr = strings.ReplaceAll(amountStr, ",", ".")
	newAmount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Invalid amount. Please enter a valid number.",
			nil,
		)
		return err
	}

	if newAmount <= 0 {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Amount must be greater than zero.",
			nil,
		)
		return err
	}

	// Update the transaction
	transaction.Amount = newAmount

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
						Text:         "Edit category",
						CallbackData: "transactions.edit.category",
					},
				},
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
				},
				{
					{
						Text:         "Edit amount",
						CallbackData: "transactions.edit.amount",
					},
				},
				{
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
					},
				},
			},
		},
	})

	return nil
}

func (c *Client) editTransactionCategory(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	var transaction model.Transaction
	err := json.Unmarshal([]byte(user.Session.Body), &transaction)
	if err != nil {
		return fmt.Errorf("failed to extract transaction from the session: %w", err)
	}

	if !model.IsValidTransactionCategory(ctx.Message.Text) {
		b.SendMessage(ctx.EffectiveSender.ChatId, "Invalid category, please try again.", nil)
		return fmt.Errorf("invalid category: %s", ctx.Message.Text)
	}
	transaction.Category = model.TransactionCategory(ctx.Message.Text)

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
						Text:         "Edit category",
						CallbackData: "transactions.edit.category",
					},
				},
				{
					{
						Text:         "Edit date",
						CallbackData: "transactions.edit.date",
					},
				},
				{
					{
						Text:         "Edit amount",
						CallbackData: "transactions.edit.amount",
					},
				},
				{
					{
						Text:         "Cancel",
						CallbackData: "transactions.cancel",
					},
					{
						Text:         "Confirm",
						CallbackData: "transactions.confirm",
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

	transaction.TgID = user.TgID
	transaction.Currency = model.CurrencyEUR

	err = c.Repositories.Transactions.Add(transaction)
	if err != nil {
		SendMessage(ctx, b, "There has been an error saving your transaction, please retry", nil)
		c.Logger.Errorln("failed to add transaction", err)

		// Reset the state
		user.Session.State = model.StateInsertingIncome
		if transaction.Type == model.TypeExpense {
			user.Session.State = model.StateInsertingExpense
		}
		err = c.Repositories.Users.Update(&user)
		if err != nil {
			return fmt.Errorf("failed to set user data to reset the state: %w", err)
		}

		return fmt.Errorf("failed to add transaction: %w", err)
	}

	user.Session.State = model.StateNormal
	user.Session.Body = ""

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
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
	return c.SendHomeKeyboard(b, ctx, fmt.Sprintf("%s Your transaction has been saved!", emoji))
}

// Cancel returns to normal state.
func (c *Client) Cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	c.CleanupKeyboard(b, ctx)
	c.SendHomeKeyboard(b, ctx, "Your operation has been canceled!\nWhat else can I do for you?\n\n/edit - Edit a transaction\n/delete - Delete a transaction\n/list - List your transactions\n/month Month Recap\n/year Year Recap")

	return ext.EndGroups
}
