package client

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (c *Client) SendAddTransactionKeyboard(b *gotgbot.Bot, ctx *ext.Context, text string) (*gotgbot.Message, error) {
	return b.SendMessage(ctx.EffectiveSender.ChatId, text, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: "💰 Add Income", CallbackData: "transactions.new.income"},
					{Text: "💸 Add Expense", CallbackData: "transactions.new.expense"},
				},
			},
		},
	})
}

func (c *Client) SendHomeKeyboard(b *gotgbot.Bot, ctx *ext.Context, text string) error {
	var err error
	keyboard := [][]gotgbot.InlineKeyboardButton{
		{
			{Text: "💰 Add Income", CallbackData: "transactions.new.income"},
			{Text: "💸 Add Expense", CallbackData: "transactions.new.expense"},
		},
		{
			{Text: "✏️ Edit", CallbackData: "home.edit"},
			{Text: "🗑 Delete", CallbackData: "home.delete"},
		},
		{
			{Text: "Week Recap", CallbackData: "home.week"},
			{Text: "Month Recap", CallbackData: "home.month"},
		},
		{
			{Text: "Year Recap", CallbackData: "home.year"},
			{Text: "📄 Full List", CallbackData: "home.list"},
		},
	}

	// Send or update message
	if ctx.CallbackQuery != nil {
		_, _, err = ctx.CallbackQuery.Message.EditText(b, text, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
	} else {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, text, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
	}

	return err
}

func (c *Client) CleanupKeyboard(b *gotgbot.Bot, ctx *ext.Context) error {
	var err error
	// Cleanup inline keyboard if exists
	if ctx.CallbackQuery != nil {
		_, _, err = ctx.CallbackQuery.Message.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{},
			},
		})
	}
	// Cleanup markup keyboard if exists
	if ctx.Message != nil {
		_, err = ctx.Message.Reply(b, "", &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.ReplyKeyboardRemove{},
		})
	}
	return err
}

// SendMessage abstracts the sending of a message regardless it's a callback from inline keyboard or a "top level" message
func SendMessage(ctx *ext.Context, b *gotgbot.Bot, message string, keyboard [][]gotgbot.InlineKeyboardButton) error {
	var err error

	// Send or update message
	if ctx.CallbackQuery != nil {
		_, _, err = ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
	} else {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
	}

	return err
}
