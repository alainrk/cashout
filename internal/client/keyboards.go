package client

import (
	"fmt"

	"cashout/internal/model"
	"cashout/internal/utils"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
)

var incomeCategories = []model.TransactionCategory{
	model.CategorySalary,
	model.CategoryOtherIncomes,
}

var expenseCategories = []model.TransactionCategory{
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

// BuildCategoryInlineKeyboard renders the category picker used by every flow.
//
//   - txType: if "" show both income+expense; if Income/Expense filter accordingly.
//   - callbackPrefix: e.g. "list.cat" produces buttons "list.cat.<CATEGORY>".
//   - cancelCallback: full callback string for the Cancel row (empty to omit).
//   - includeAll: prepend an "🔍 All Categories" row with callback "<prefix>.all".
func BuildCategoryInlineKeyboard(txType model.TransactionType, callbackPrefix, cancelCallback string, includeAll bool) [][]gotgbot.InlineKeyboardButton {
	var keyboard [][]gotgbot.InlineKeyboardButton

	if includeAll {
		keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         "🔍 All Categories",
				CallbackData: fmt.Sprintf("%s.all", callbackPrefix),
			},
		})
	}

	if txType == "" || txType == model.TypeIncome {
		for _, cat := range incomeCategories {
			keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
				{
					Text:         fmt.Sprintf("%s %s", utils.GetCategoryEmoji(cat), cat),
					CallbackData: fmt.Sprintf("%s.%s", callbackPrefix, cat),
				},
			})
		}
	}

	if txType == "" || txType == model.TypeExpense {
		for i := 0; i < len(expenseCategories); i += 2 {
			row := []gotgbot.InlineKeyboardButton{
				{
					Text:         fmt.Sprintf("%s %s", utils.GetCategoryEmoji(expenseCategories[i]), expenseCategories[i]),
					CallbackData: fmt.Sprintf("%s.%s", callbackPrefix, expenseCategories[i]),
				},
			}
			if i+1 < len(expenseCategories) {
				row = append(row, gotgbot.InlineKeyboardButton{
					Text:         fmt.Sprintf("%s %s", utils.GetCategoryEmoji(expenseCategories[i+1]), expenseCategories[i+1]),
					CallbackData: fmt.Sprintf("%s.%s", callbackPrefix, expenseCategories[i+1]),
				})
			}
			keyboard = append(keyboard, row)
		}
	}

	if cancelCallback != "" {
		keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         "❌ Cancel",
				CallbackData: cancelCallback,
			},
		})
	}

	return keyboard
}
