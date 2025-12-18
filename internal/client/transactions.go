package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cashout/internal/model"
	"cashout/internal/utils"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// AddTransactionExpense handles the add expense intent at top-level
func (c *Client) AddTransactionExpense(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Reset user state and start delete search flow
	user.Session.State = model.StateInsertingExpense
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "Sure! To add a new <b>expense</b>:\nTell me category, amount and description. You can also specify a date and change it later, today is default.\n\n<i>Examples:</i>\n<code>Irish Pub 3.4</code>\n<code>January salary 3k 10/01</code>", &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Cancel", CallbackData: "transactions.cancel"}},
			},
		},
	})
	if err != nil {
		return err
	}

	return ext.ContinueGroups
}

// AddTransactionIncome handles the add income intent at top-level
func (c *Client) AddTransactionIncome(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Reset user state and start delete search flow
	user.Session.State = model.StateInsertingIncome
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "Sure! To add a new <b>income</b>:\nTell me category, amount and description. You can also specify a date and change it later, today is default.\n\n<i>Examples:</i>\n<code>Irish Pub 3.4</code>\n<code>January salary 3k 10/01</code>", &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Cancel", CallbackData: "transactions.cancel"}},
			},
		},
	})
	if err != nil {
		return err
	}

	return ext.ContinueGroups
}

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
		_, err = c.SendAddTransactionKeyboard(b, ctx, "Invalid action, add a transaction.")
		return err
	}

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	_, _, err = msg.EditText(b, fmt.Sprintf("Sure! To add a new <b>%s</b>:\nTell me category, amount and description. You can also specify a date and change it later, today is default.\n\n<i>Examples:</i>\n<code>Irish Pub 3.4</code>\n<code>January salary 3k 10/01</code>", action), &gotgbot.EditMessageTextOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return err
	}

	_, _, err = msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Cancel", CallbackData: "transactions.cancel"}},
			},
		},
	})
	if err != nil {
		return err
	}

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
		_, err := c.SendAddTransactionKeyboard(b, ctx, "Invalid action, add a transaction.")
		return err
	}

	extractedTransaction, err := c.LLM.ExtractTransaction(ctx.Message.Text, transactionType)
	if err != nil {
		msg := "I'm sorry, I couldn't understand your transaction!"
		_, errm := ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		err = errors.Join(err, errm)
		return err
	}

	if extractedTransaction.Amount == 0 {
		msg := "I'm sorry, I couldn't understand your transaction!"
		_, errm := ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		err = errors.Join(err, errm)
		return err
	}

	// Convert to model.Transaction and save immediately
	transaction := model.Transaction{
		TgID:        user.TgID,
		Type:        extractedTransaction.Type,
		Category:    model.TransactionCategory(extractedTransaction.Category),
		Amount:      extractedTransaction.Amount,
		Description: extractedTransaction.Description,
		Date:        extractedTransaction.Date,
		Currency:    model.CurrencyEUR,
	}

	err = c.Repositories.Transactions.Add(&transaction)
	if err != nil {
		err = errors.Join(err, SendMessage(ctx, b, "There has been an error saving your transaction, please retry", nil))
		c.Logger.Errorln("failed to add transaction", err)
		return fmt.Errorf("failed to add transaction: %w", err)
	}

	// Store the transaction ID in session for potential edits
	user.Session.State = model.StateEditingNewTransaction
	user.Session.Body = strconv.FormatInt(transaction.ID, 10)

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	msg := fmt.Sprintf("%s <b>Transaction saved!</b>\n\n%s (€ %.2f), %s on %s", emoji, transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
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
						Text:         "Delete",
						CallbackData: fmt.Sprintf("transactions.delete.%d", transaction.ID),
					},
					{
						Text:         "Home",
						CallbackData: "transactions.home",
					},
				},
			},
		},
	})
	if err != nil {
		c.Logger.Errorln("failed to send saved message", err)
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

	// Load transaction from DB using the ID stored in session
	transactionID, err := strconv.ParseInt(user.Session.Body, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse transaction ID from session: %w", err)
	}

	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction from database: %w", err)
	}

	query := ctx.CallbackQuery

	field := strings.Split(query.Data, ".")[2]

	var text string
	var opts *gotgbot.SendMessageOpts

	switch field {
	case "description":
		user.Session.State = model.StateEditingTransactionDescription

		keyboard := [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Cancel",
					CallbackData: "transactions.editcancel",
				},
			},
		}
		opts = &gotgbot.SendMessageOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: keyboard}}
		text = fmt.Sprintf("Enter a new description for the transaction:\n\nCurrent: %s ", transaction.Description)
	case "amount":
		user.Session.State = model.StateEditingTransactionAmount

		keyboard := [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Cancel",
					CallbackData: "transactions.editcancel",
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
					CallbackData: "transactions.editcancel",
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
				{{Text: "Pets"}},
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

	err = c.CleanupKeyboard(b, ctx)
	if err != nil {
		return err
	}

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, text, opts)

	return err
}

func (c *Client) editTransactionDate(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	// Load transaction from DB using the ID stored in session
	transactionID, err := strconv.ParseInt(user.Session.Body, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse transaction ID from session: %w", err)
	}

	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction from database: %w", err)
	}

	// Get date from DD-MM-YYYY to date
	date, err := utils.ParseDate(ctx.Message.Text)
	if err != nil {
		fmt.Printf("failed to parse date: %v\n", err)
		_, errm := b.SendMessage(ctx.EffectiveSender.ChatId, "Invalid date, please try again.", nil)
		err = errors.Join(err, errm)
		return err
	}

	if date.After(time.Now()) {
		_, err := b.SendMessage(ctx.EffectiveSender.ChatId, "I don't support future dates, please try again.", nil)
		return errors.Join(err, fmt.Errorf("invalid date: %s", ctx.Message.Text))
	}

	// Update the transaction in DB
	transaction.Date = date
	err = c.Repositories.Transactions.Update(&transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Update session state
	user.Session.State = model.StateEditingNewTransaction
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	m := fmt.Sprintf("%s <b>Transaction updated!</b>\n\n%s (€ %.2f), %s on %s", emoji, transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
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
						Text:         "Delete",
						CallbackData: fmt.Sprintf("transactions.delete.%d", transaction.ID),
					},
					{
						Text:         "Home",
						CallbackData: "transactions.home",
					},
				},
			},
		},
	})

	return err
}

func (c *Client) editTransactionAmount(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	// Load transaction from DB using the ID stored in session
	transactionID, err := strconv.ParseInt(user.Session.Body, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse transaction ID from session: %w", err)
	}

	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction from database: %w", err)
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

	// Update the transaction in DB
	transaction.Amount = newAmount
	err = c.Repositories.Transactions.Update(&transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Update session state
	user.Session.State = model.StateEditingNewTransaction
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	m := fmt.Sprintf("%s <b>Transaction updated!</b>\n\n%s (€ %.2f), %s on %s", emoji, transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
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
						Text:         "Delete",
						CallbackData: fmt.Sprintf("transactions.delete.%d", transaction.ID),
					},
					{
						Text:         "Home",
						CallbackData: "transactions.home",
					},
				},
			},
		},
	})

	return err
}

func (c *Client) editTransactionDescription(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	// Load transaction from DB using the ID stored in session
	transactionID, err := strconv.ParseInt(user.Session.Body, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse transaction ID from session: %w", err)
	}

	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction from database: %w", err)
	}

	newDescription := strings.TrimSpace(ctx.Message.Text)
	if newDescription == "" {
		_, err = b.SendMessage(
			ctx.EffectiveSender.ChatId,
			"Description cannot be empty.",
			nil,
		)
		return err
	}

	// Update the transaction in DB
	transaction.Description = newDescription
	err = c.Repositories.Transactions.Update(&transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Update session state
	user.Session.State = model.StateEditingNewTransaction
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	m := fmt.Sprintf("%s <b>Transaction updated!</b>\n\n%s (€ %.2f), %s on %s", emoji, transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
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
						Text:         "Delete",
						CallbackData: fmt.Sprintf("transactions.delete.%d", transaction.ID),
					},
					{
						Text:         "Home",
						CallbackData: "transactions.home",
					},
				},
			},
		},
	})

	return err
}

func (c *Client) editTransactionCategory(b *gotgbot.Bot, ctx *ext.Context, user model.User) error {
	// Load transaction from DB using the ID stored in session
	transactionID, err := strconv.ParseInt(user.Session.Body, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse transaction ID from session: %w", err)
	}

	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction from database: %w", err)
	}

	if !model.IsValidTransactionCategory(ctx.Message.Text) {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "Invalid category, please try again.", nil)
		return errors.Join(err, fmt.Errorf("invalid category: %s", ctx.Message.Text))
	}

	// Update the transaction in DB
	transaction.Category = model.TransactionCategory(ctx.Message.Text)
	err = c.Repositories.Transactions.Update(&transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Update session state
	user.Session.State = model.StateEditingNewTransaction
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	m := fmt.Sprintf("%s <b>Transaction updated!</b>\n\n%s (€ %.2f), %s on %s", emoji, transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
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
						Text:         "Delete",
						CallbackData: fmt.Sprintf("transactions.delete.%d", transaction.ID),
					},
					{
						Text:         "Home",
						CallbackData: "transactions.home",
					},
				},
			},
		},
	})

	return err
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

	err = c.Repositories.Transactions.Add(&transaction)
	if err != nil {
		err = errors.Join(err, SendMessage(ctx, b, "There has been an error saving your transaction, please retry", nil))
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
	_, _, err = msg.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{},
		},
	})
	if err != nil {
		c.Logger.Errorln("failed to remove the keyboard from the previous message", err)
	}

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

	err = c.CleanupKeyboard(b, ctx)
	err = errors.Join(err, c.SendHomeKeyboard(b, ctx, "Your operation has been canceled!\nWhat else can I do for you?\n\n/edit - Edit a transaction\n/delete - Delete a transaction\n/search - Search transactions\n/list - List your transactions\n/week Week Recap\n/month Month Recap\n/year Year Recap\n/export - Export all transactions to CSV"))

	return err
}

// DeleteNewTransaction deletes a newly created transaction and returns to home.
func (c *Client) DeleteNewTransaction(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: transactions.delete.TRANSACTION_ID)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	transactionID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %v", err)
	}

	// Get the transaction before deleting it (for confirmation message)
	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	// Verify ownership
	if transaction.TgID != user.TgID {
		_, _, err = ctx.CallbackQuery.Message.EditText(
			b,
			"This transaction doesn't belong to you.",
			&gotgbot.EditMessageTextOpts{},
		)
		return err
	}

	// Delete the transaction
	err = c.Repositories.Transactions.Delete(transactionID, user.TgID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	// Send success message and return to home
	return c.SendHomeKeyboard(b, ctx, "Transaction deleted!")
}

// TransactionHome returns to home after adding/editing a transaction.
func (c *Client) TransactionHome(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	return c.SendHomeKeyboard(b, ctx, "What can I do for you?")
}

// EditCancel cancels the current edit and returns to the transaction view.
func (c *Client) EditCancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Load transaction from DB using the ID stored in session
	transactionID, err := strconv.ParseInt(user.Session.Body, 10, 64)
	if err != nil {
		// If we can't parse the ID, just go home
		user.Session.State = model.StateNormal
		user.Session.Body = ""
		err = c.Repositories.Users.Update(&user)
		if err != nil {
			return fmt.Errorf("failed to set user data: %w", err)
		}
		return c.SendHomeKeyboard(b, ctx, "What can I do for you?")
	}

	transaction, err := c.Repositories.Transactions.GetByID(transactionID)
	if err != nil {
		// If we can't find the transaction, just go home
		user.Session.State = model.StateNormal
		user.Session.Body = ""
		err = c.Repositories.Users.Update(&user)
		if err != nil {
			return fmt.Errorf("failed to set user data: %w", err)
		}
		return c.SendHomeKeyboard(b, ctx, "What can I do for you?")
	}

	// Update session state
	user.Session.State = model.StateEditingNewTransaction
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	err = c.CleanupKeyboard(b, ctx)
	if err != nil {
		c.Logger.Warnln("failed to cleanup keyboard", err)
	}

	emoji := "💰"
	if transaction.Type == model.TypeExpense {
		emoji = "💸"
	}

	m := fmt.Sprintf("%s <b>Transaction saved!</b>\n\n%s (€ %.2f), %s on %s", emoji, transaction.Category, transaction.Amount, transaction.Description, transaction.Date.Format("02-01-2006"))
	_, err = b.SendMessage(ctx.EffectiveSender.ChatId, m, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Edit description",
						CallbackData: "transactions.edit.description",
					},
				},
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
						Text:         "Delete",
						CallbackData: fmt.Sprintf("transactions.delete.%d", transaction.ID),
					},
					{
						Text:         "Home",
						CallbackData: "transactions.home",
					},
				},
			},
		},
	})

	return err
}
