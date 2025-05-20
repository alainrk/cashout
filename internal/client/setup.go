package client

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// Create a matcher which only matches text which is not a command.
func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}

func confirmCommand(msg *gotgbot.Message) bool {
	return message.Text(msg) && strings.Trim(msg.Text, " ") == "Confirm"
}

func cancelText(msg *gotgbot.Message) bool {
	return message.Text(msg) && strings.ToLower(strings.Trim(msg.Text, " ")) == "cancel"
}

func SetupHandlers(dispatcher *ext.Dispatcher, c *Client) {
	// Top-level message for LLM goes into AddTransaction and gets the expense/income intent from user session state.
	dispatcher.AddHandler(handlers.NewMessage(noCommands, c.FreeTextRouter))
	dispatcher.AddHandler(handlers.NewMessage(cancelText, c.Cancel))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("transactions.new."), c.AddTransactionIntent))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("transactions.edit."), c.EditTransactionIntent))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("transactions.cancel"), c.Cancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("transactions.confirm"), c.Confirm))

	dispatcher.AddHandler(handlers.NewCommand("list", func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "/list")))
		return c.ListTransactions(b, ctx)
	}))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("list.cancel"), c.Cancel)) // Not instrumenting callbacks yet
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.year."), c.ListYearNavigation)) // Not instrumenting callbacks yet
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.month."), c.ListMonthTransactions)) // Not instrumenting callbacks yet
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.page."), c.ListTransactionPage))   // Not instrumenting callbacks yet

	dispatcher.AddHandler(handlers.NewCommand("edit", func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "/edit")))
		return c.EditTransactions(b, ctx)
	}))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("edit.page."), c.EditTransactionPage))     // Not instrumenting callbacks yet
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("edit.select."), c.EditTransactionSelect)) // Not instrumenting callbacks yet
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("edit.field."), c.EditTransactionField)) // Not instrumenting callbacks yet

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.page."), c.DeleteTransactionPage))       // Not instrumenting callbacks yet
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.confirm."), c.DeleteTransactionConfirm)) // Not instrumenting callbacks yet

	dispatcher.AddHandler(handlers.NewCommand("cancel", func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "/cancel")))
		return c.Cancel(b, ctx)
	}))
	dispatcher.AddHandler(handlers.NewCommand("delete", func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "/delete")))
		return c.DeleteTransactions(b, ctx)
	}))
	dispatcher.AddHandler(handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "/start")))
		return c.Start(b, ctx)
	}))
	// "/new" is an alias for "/start"
	dispatcher.AddHandler(handlers.NewCommand("new", func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "/new")))
		return c.Start(b, ctx) // It calls c.Start, so it's like an alias
	}))
	dispatcher.AddHandler(handlers.NewCommand("month", func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "/month")))
		return c.MonthRecap(b, ctx)
	}))
	dispatcher.AddHandler(handlers.NewCommand("year", func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "/year")))
		return c.YearRecap(b, ctx)
	}))

	// Callbacks that effectively act as commands for metrics
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.month"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		// For callbacks that are essentially commands, we can log them with a "command.name" like attribute.
		// Using a "callback." prefix to distinguish from direct command entries.
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "callback_home.month")))
		return c.MonthRecap(b, ctx)
	}))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.year"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "callback_home.year")))
		return c.YearRecap(b, ctx)
	}))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.list"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "callback_home.list")))
		return c.ListTransactions(b, ctx)
	}))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.delete"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "callback_home.delete")))
		return c.DeleteTransactions(b, ctx)
	}))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.edit"), func(b *gotgbot.Bot, ctx *ext.Context) error {
		reqCtx := ctx.Request.Context()
		if reqCtx == nil {
			reqCtx = context.Background()
		}
		c.CommandCounter.Add(reqCtx, 1, metric.WithAttributes(attribute.String("command.name", "callback_home.edit")))
		return c.EditTransactions(b, ctx)
	}))
}
