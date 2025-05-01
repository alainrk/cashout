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

func (c *Client) SendHomeKeyboard(b *gotgbot.Bot, ctx *ext.Context, text string) (*gotgbot.Message, error) {
	return b.SendMessage(ctx.EffectiveSender.ChatId, text, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: "💰 Add Income", CallbackData: "transactions.new.income"},
					{Text: "💸 Add Expense", CallbackData: "transactions.new.expense"},
				},
				// TODO: (see /cmd/server/main.go)
				// {
				// 	{Text: "🗑 Delete Transaction", CallbackData: "home.delete"},
				// 	{Text: "📄 Full List", CallbackData: "home.list"},
				// },
				// {
				// 	{Text: "Month Recap", CallbackData: "home.month"},
				// 	{Text: "Year Recap", CallbackData: "home.year"},
				// },
			},
		},
	})
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
