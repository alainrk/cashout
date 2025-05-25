package client

import (
	"fmt"
	"time"

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
			{Text: "📄 Full List", CallbackData: "home.list"},
		},
		{
			{Text: "Current Week Recap", CallbackData: "home.week"},
		},
		{
			{Text: "Year Recap", CallbackData: "home.year"},
			{Text: "Month Recap", CallbackData: "home.month"},
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

// sendRecapWithNavigation sends a recap message with navigation buttons for previous/next period
func (c *Client) sendRecapWithNavigation(b *gotgbot.Bot, ctx *ext.Context, text string, recapType string, year int, month int) error {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Create navigation row with Previous/Next buttons
	var navRow []gotgbot.InlineKeyboardButton

	switch recapType {

	case "month":
		// Calculate previous month
		prevYear, prevMonth := year, month-1
		if prevMonth < 1 {
			prevMonth = 12
			prevYear--
		}

		// Calculate next month
		nextYear, nextMonth := year, month+1
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}

		// Add Previous button if not too far in the past
		if prevYear >= MIN_YEAR_ALLOWED {
			navRow = append(navRow, gotgbot.InlineKeyboardButton{
				Text:         "⬅️ Previous Month",
				CallbackData: fmt.Sprintf("monthrecap.month.%d.%02d", prevYear, prevMonth),
			})
		}

		// Add Next button if not in the future
		currentTime := time.Now()
		if nextYear < currentTime.Year() || (nextYear == currentTime.Year() && nextMonth <= int(currentTime.Month())) {
			navRow = append(navRow, gotgbot.InlineKeyboardButton{
				Text:         "Next Month ➡️",
				CallbackData: fmt.Sprintf("monthrecap.month.%d.%02d", nextYear, nextMonth),
			})
		}
	case "year":
		// Add Previous button if not too far in the past
		if year > MIN_YEAR_ALLOWED {
			navRow = append(navRow, gotgbot.InlineKeyboardButton{
				Text:         "⬅️ Previous Year",
				CallbackData: fmt.Sprintf("yearrecap.year.%d", year-1),
			})
		}

		// Add Next button if not in the future
		currentYear := time.Now().Year()
		if year < currentYear {
			navRow = append(navRow, gotgbot.InlineKeyboardButton{
				Text:         "Next Year ➡️",
				CallbackData: fmt.Sprintf("yearrecap.year.%d", year+1),
			})
		}
	default:
		c.Logger.Warnf("Unknown recap type: %s", recapType)
	}

	// Add navigation row if it has any buttons
	if len(navRow) > 0 {
		keyboard = append(keyboard, navRow)
	}

	// Add standard home keyboard buttons
	keyboard = append(keyboard, [][]gotgbot.InlineKeyboardButton{
		{
			{Text: "💰 Add Income", CallbackData: "transactions.new.income"},
			{Text: "💸 Add Expense", CallbackData: "transactions.new.expense"},
		},
		{
			{Text: "✏️ Edit", CallbackData: "home.edit"},
			{Text: "🗑 Delete", CallbackData: "home.delete"},
		},
		{
			{Text: "📄 Full List", CallbackData: "home.list"},
		},
		{
			{Text: "Current Week Recap", CallbackData: "home.week"},
		},
		{
			{Text: "Year Recap", CallbackData: "home.year"},
			{Text: "Month Recap", CallbackData: "home.month"},
		},
	}...)

	// Send or update message
	var err error
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
