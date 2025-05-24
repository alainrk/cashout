package client

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

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
