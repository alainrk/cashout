package client

import (
	"strings"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// Create a matcher which only matches text which is not a command.
func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
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
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("transactions.delete."), c.DeleteNewTransaction))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("transactions.cancel"), c.Cancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("transactions.editcancel"), c.EditCancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("transactions.home"), c.TransactionHome))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("transactions.confirm"), c.Confirm))

	dispatcher.AddHandler(handlers.NewCommand("list", c.ListTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("list.cancel"), c.Cancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.year."), c.ListYearNavigation))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.month."), c.ListMonthTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("list.page."), c.ListTransactionPage))

	dispatcher.AddHandler(handlers.NewCommand("edit", c.EditTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("edit.page."), c.EditTransactionPage))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("edit.select."), c.EditTransactionSelect))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("edit.field."), c.EditTransactionField))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("edit.done"), c.EditDone))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("edit.search.category."), c.EditSearchCategorySelected))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("edit.search.page."), c.EditSearchResultsPage))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("edit.search.select."), c.EditSearchTransactionSelected))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("edit.search.cancel"), c.EditSearchCancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("edit.search.home"), c.EditSearchHome))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("edit.search.new"), c.EditSearchNew))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("edit.search.noop"), c.EditSearchNoop))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.page."), c.DeleteTransactionPage))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.showconfirm."), c.ShowDeleteConfirmation))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.confirm."), c.DeleteTransactionConfirm))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.search.category."), c.DeleteSearchCategorySelected))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.search.page."), c.DeleteSearchResultsPage))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("delete.search.select."), c.DeleteSearchTransactionSelected))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("delete.search.cancel"), c.DeleteSearchCancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("delete.search.home"), c.DeleteSearchHome))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("delete.search.new"), c.DeleteSearchNew))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("delete.search.noop"), c.DeleteSearchNoop))

	dispatcher.AddHandler(handlers.NewCommand("cancel", c.Cancel))
	dispatcher.AddHandler(handlers.NewCommand("delete", c.DeleteTransactions))
	dispatcher.AddHandler(handlers.NewCommand("start", c.Start))
	dispatcher.AddHandler(handlers.NewCommand("new", c.Start))
	dispatcher.AddHandler(handlers.NewCommand("week", c.WeekRecap))
	dispatcher.AddHandler(handlers.NewCommand("month", c.MonthRecap))
	dispatcher.AddHandler(handlers.NewCommand("year", c.YearRecap))
	dispatcher.AddHandler(handlers.NewCommand("export", c.ExportTransactions))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("monthrecap.cancel"), c.Cancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("monthrecap.year."), c.MonthRecapYearNavigation))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("monthrecap.month."), c.MonthRecapSelected))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("yearrecap.cancel"), c.Cancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("yearrecap.year."), c.YearRecapSelected))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.week"), c.WeekRecap))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.month"), c.MonthRecap))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.year"), c.YearRecap))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.list"), c.ListTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.delete"), c.DeleteTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.edit"), c.EditTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("home.search"), c.SearchTransactions))

	dispatcher.AddHandler(handlers.NewCommand("search", c.SearchTransactions))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("search.category."), c.SearchCategorySelected))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("search.page."), c.SearchResultsPage))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("search.cancel"), c.SearchCancel))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("search.home"), c.SearchHome))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("search.new"), c.SearchNew))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("search.noop"), c.SearchNoop))
}
