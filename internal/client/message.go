package client

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

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

func GetMessageFromContext(ctx *ext.Context) string {
	var msg string

	if ctx.Message != nil {
		return ctx.Message.Text
	}

	if ctx.CallbackQuery != nil {
		msg = ctx.CallbackQuery.Data
	}

	return msg
}
