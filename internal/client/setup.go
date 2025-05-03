package client

import (
	"strings"

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

	dispatcher.AddHandler(handlers.NewCommand("list", c.ListTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.year."), c.ListYearNavigation))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.month."), c.ListMonthTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.page."), c.ListTransactionPage))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.page."), c.DeleteTransactionPage))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.confirm."), c.DeleteTransactionConfirm))

	dispatcher.AddHandler(handlers.NewCommand("cancel", c.Cancel))
	dispatcher.AddHandler(handlers.NewCommand("delete", c.DeleteTransactions))
	dispatcher.AddHandler(handlers.NewCommand("start", c.Start))
	dispatcher.AddHandler(handlers.NewCommand("new", c.Start))
	dispatcher.AddHandler(handlers.NewCommand("month", c.MonthRecap))
	dispatcher.AddHandler(handlers.NewCommand("year", c.YearRecap))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.month"), c.MonthRecap))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.year"), c.YearRecap))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.list"), c.ListTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.delete"), c.DeleteTransactions))
}
