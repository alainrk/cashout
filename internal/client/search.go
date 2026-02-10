package client

import (
	"cashout/internal/model"
	"cashout/internal/utils"
	"fmt"
	"strconv"
	"strings"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// SearchTransactions initiates the search flow
func (c *Client) SearchTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Reset user state
	user.Session.State = model.StateSelectingSearchCategory
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Show category selection
	return c.showSearchCategorySelection(b, ctx)
}

// showSearchCategorySelection displays the category selection keyboard
func (c *Client) showSearchCategorySelection(b *gotgbot.Bot, ctx *ext.Context) error {
	var keyboard [][]gotgbot.InlineKeyboardButton

	// Add "All" option first
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "🔍 All Categories",
			CallbackData: "search.category.all",
		},
	})

	// Add separator
	// keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
	// 	{
	// 		Text:         "── Income Categories ──",
	// 		CallbackData: "search.noop",
	// 	},
	// })

	// Add income categories
	incomeCategories := []model.TransactionCategory{
		model.CategorySalary,
		model.CategoryOtherIncomes,
	}

	for _, cat := range incomeCategories {
		emoji := utils.GetCategoryEmoji(cat)
		keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%s %s", emoji, cat),
				CallbackData: fmt.Sprintf("search.category.%s", cat),
			},
		})
	}

	// Add separator
	// keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
	// 	{
	// 		Text:         "── Expense Categories ──",
	// 		CallbackData: "search.noop",
	// 	},
	// })

	// Add expense categories in rows of 2
	expenseCategories := []model.TransactionCategory{
		model.CategoryCar,
		model.CategoryClothes,
		model.CategoryGrocery,
		model.CategoryHouse,
		model.CategoryBills,
		model.CategoryEntertainment,
		model.CategorySport,
		model.CategoryEatingOut,
		model.CategoryTransport,
		model.CategoryLearning,
		model.CategoryToiletry,
		model.CategoryHealth,
		model.CategoryTech,
		model.CategoryGifts,
		model.CategoryTravel,
		model.CategoryPets,
		model.CategoryOtherExpenses,
	}

	for i := 0; i < len(expenseCategories); i += 2 {
		row := []gotgbot.InlineKeyboardButton{}

		emoji := utils.GetCategoryEmoji(expenseCategories[i])
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s %s", emoji, expenseCategories[i]),
			CallbackData: fmt.Sprintf("search.category.%s", expenseCategories[i]),
		})

		// Add second button if exists
		if i+1 < len(expenseCategories) {
			emoji2 := utils.GetCategoryEmoji(expenseCategories[i+1])
			row = append(row, gotgbot.InlineKeyboardButton{
				Text:         fmt.Sprintf("%s %s", emoji2, expenseCategories[i+1]),
				CallbackData: fmt.Sprintf("search.category.%s", expenseCategories[i+1]),
			})
		}

		keyboard = append(keyboard, row)
	}

	// Add cancel button
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "❌ Cancel",
			CallbackData: "search.cancel",
		},
	})

	message := "🔍 <b>Search Transactions</b>\n\nSelect a category to search in:"

	// Send or update message
	if ctx.CallbackQuery != nil {
		_, _, err := ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	} else {
		_, err := b.SendMessage(ctx.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}
}

// SearchCategorySelected handles category selection
func (c *Client) SearchCategorySelected(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: search.category.CATEGORY)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid callback data format")
	}

	category := parts[2]

	// Store selected category in session
	user.Session.State = model.StateEnteringSearchQuery
	user.Session.Body = category
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Ask for search query
	categoryText := "all categories"
	if category != "all" {
		emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
		categoryText = fmt.Sprintf("%s %s", emoji, category)
	}

	_, _, err = ctx.CallbackQuery.Message.EditText(
		b,
		fmt.Sprintf("🔍 Searching in <b>%s</b>\n\nEnter your search text or tap Show All:", categoryText),
		&gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "📋 Show All",
							CallbackData: "search.showall",
						},
					},
					{
						{
							Text:         "❌ Cancel",
							CallbackData: "search.cancel",
						},
					},
				},
			},
		},
	)

	return err
}

// SearchQueryEntered handles the search query input
func (c *Client) SearchQueryEntered(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Get search query
	searchQuery := strings.TrimSpace(ctx.Message.Text)
	if searchQuery == "" {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "Search query cannot be empty. Please try again.", nil)
		return err
	}

	// Get category from session
	category := user.Session.Body

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	// Perform search and show results
	return c.showSearchResults(b, ctx, user, category, searchQuery, 0)
}

// showSearchResults displays paginated search results
func (c *Client) showSearchResults(b *gotgbot.Bot, ctx *ext.Context, user model.User, category, searchQuery string, offset int) error {
	limit := 10

	// Perform search
	var transactions []model.Transaction
	var total int64
	var err error

	if category == "all" {
		transactions, total, err = c.Repositories.Transactions.SearchUserTransactions(
			user.TgID,
			searchQuery,
			"",
			offset,
			limit,
		)
	} else {
		transactions, total, err = c.Repositories.Transactions.SearchUserTransactions(
			user.TgID,
			searchQuery,
			category,
			offset,
			limit,
		)
	}

	if err != nil {
		return fmt.Errorf("failed to search transactions: %w", err)
	}

	if total == 0 {
		message := fmt.Sprintf("🔍 No transactions found matching \"%s\"", searchQuery)
		if category != "all" {
			emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
			message = fmt.Sprintf("🔍 No transactions found matching \"%s\" in %s %s", searchQuery, emoji, category)
		}

		if ctx.CallbackQuery != nil {
			_, _, err = ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
						{
							{
								Text:         "🔍 New Search",
								CallbackData: "search.new",
							},
							{
								Text:         "🏠 Home",
								CallbackData: "search.home",
							},
						},
					},
				},
			})
		} else {
			_, err = b.SendMessage(ctx.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
						{
							{
								Text:         "🔍 New Search",
								CallbackData: "search.new",
							},
							{
								Text:         "🏠 Home",
								CallbackData: "search.home",
							},
						},
					},
				},
			})
		}
		return err
	}

	// Format search results
	message := formatSearchResults(transactions, searchQuery, category, offset, int(total))

	// Create pagination keyboard
	keyboard := createSearchPaginationKeyboard(category, searchQuery, offset, limit, int(total))

	// Send or update message
	if ctx.CallbackQuery != nil {
		_, _, err = ctx.CallbackQuery.Message.EditText(b, message, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	} else {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, message, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			},
		})
		return err
	}
}

// SearchResultsPage handles pagination for search results
func (c *Client) SearchResultsPage(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	query := ctx.CallbackQuery

	// Parse callback data (format: search.page.CATEGORY.OFFSET.QUERY)
	parts := strings.Split(query.Data, ".")
	if len(parts) != 5 {
		return fmt.Errorf("invalid callback data format")
	}

	category := parts[2]
	offset, err := strconv.Atoi(parts[3])
	if err != nil {
		return fmt.Errorf("invalid offset: %v", err)
	}
	searchQuery := parts[4]

	return c.showSearchResults(b, ctx, user, category, searchQuery, offset)
}

// formatSearchResults formats the search results for display
func formatSearchResults(transactions []model.Transaction, searchQuery, category string, offset, total int) string {
	var msg strings.Builder

	msg.WriteString("🔍 <b>Search Results</b>\n")

	// Show query unless it's a wildcard "Show All"
	if searchQuery != "%" {
		msg.WriteString(fmt.Sprintf("Query: \"%s\"", searchQuery))
	}

	if category != "all" {
		emoji := utils.GetCategoryEmoji(model.TransactionCategory(category))
		msg.WriteString(fmt.Sprintf(" in %s %s", emoji, category))
	}

	msg.WriteString(fmt.Sprintf("\nShowing %d–%d of %d\n\n", offset+1, offset+len(transactions), total))

	for _, t := range transactions {
		emoji := utils.GetCategoryEmoji(t.Category)
		sign := "-"
		if t.Type == model.TypeIncome {
			sign = "+"
		}

		// Highlight the search term in description (skip for wildcard)
		desc := t.Description
		if searchQuery != "%" {
			if idx := strings.Index(strings.ToLower(desc), strings.ToLower(searchQuery)); idx != -1 {
				desc = desc[:idx] + "<b>" + desc[idx:idx+len(searchQuery)] + "</b>" + desc[idx+len(searchQuery):]
			}
		}

		msg.WriteString(fmt.Sprintf("%s %s · %s€%.2f · %s\n",
			emoji, desc, sign, t.Amount, t.Date.Format("02/01")))
	}

	return msg.String()
}

// createSearchPaginationKeyboard creates pagination buttons for search results
func createSearchPaginationKeyboard(category, searchQuery string, offset, limit, total int) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton
	var navigationRow []gotgbot.InlineKeyboardButton

	// Previous page button
	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "⬅️ Previous",
			CallbackData: fmt.Sprintf("search.page.%s.%d.%s", category, prevOffset, searchQuery),
		})
	}

	// Page indicator
	currentPage := (offset / limit) + 1
	totalPages := (total + limit - 1) / limit
	navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
		Text:         fmt.Sprintf("%d/%d", currentPage, totalPages),
		CallbackData: "search.noop",
	})

	// Next page button
	if offset+limit < total {
		nextOffset := offset + limit
		navigationRow = append(navigationRow, gotgbot.InlineKeyboardButton{
			Text:         "Next ➡️",
			CallbackData: fmt.Sprintf("search.page.%s.%d.%s", category, nextOffset, searchQuery),
		})
	}

	if len(navigationRow) > 0 {
		keyboard = append(keyboard, navigationRow)
	}

	// Action buttons
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{
			Text:         "🔍 New Search",
			CallbackData: "search.new",
		},
		{
			Text:         "🏠 Home",
			CallbackData: "search.home",
		},
	})

	return keyboard
}

// SearchCancel handles search cancellation
func (c *Client) SearchCancel(b *gotgbot.Bot, ctx *ext.Context) error {
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
		return fmt.Errorf("failed to update user data: %w", err)
	}

	return c.Cancel(b, ctx)
}

// SearchHome returns to home screen
func (c *Client) SearchHome(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.SendHomeKeyboard(b, ctx, "What can I do for you?")
}

// SearchNew starts a new search
func (c *Client) SearchNew(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.SearchTransactions(b, ctx)
}

// SearchNoop handles no-op callbacks (like separators)
func (c *Client) SearchNoop(b *gotgbot.Bot, ctx *ext.Context) error {
	// Answer callback query to remove loading state
	_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	return err
}

// SearchShowAll handles the "Show All" button — searches with wildcard
func (c *Client) SearchShowAll(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Get category from session
	category := user.Session.Body

	// Reset user state
	user.Session.State = model.StateNormal
	user.Session.Body = ""
	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	return c.showSearchResults(b, ctx, user, category, "%", 0)
}
