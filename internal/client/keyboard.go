package client

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (c *Client) SendTransactionKeyboard(b *gotgbot.Bot, ctx *ext.Context, text string) (*gotgbot.Message, error) {
	return ctx.EffectiveMessage.Reply(b, text, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			Keyboard: [][]gotgbot.KeyboardButton{
				{
					{
						Text: "Add Income",
					},
					{
						Text: "Add Expense",
					},
				},
			},
			IsPersistent:   true,
			ResizeKeyboard: true,
		},
	})
}

func (c *Client) SendConfirmKeyboard(b *gotgbot.Bot, ctx *ext.Context, text string) (*gotgbot.Message, error) {
	return ctx.EffectiveMessage.Reply(b, text, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			Keyboard: [][]gotgbot.KeyboardButton{
				{
					{
						Text: "Confirm",
					},
					{
						Text: "Cancel",
					},
				},
			},
			IsPersistent:    false,
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
	})
}
