package utils

import "cashout/internal/model"

// GetCategoryEmoji returns the appropriate emoji for a transaction category
func GetCategoryEmoji(category model.TransactionCategory) string {
	emojiMap := map[model.TransactionCategory]string{
		model.CategorySalary:        "💵",
		model.CategoryOtherIncomes:  "💵",
		model.CategoryCar:           "🚗",
		model.CategoryClothes:       "👕",
		model.CategoryGrocery:       "🛒",
		model.CategoryHouse:         "🏠",
		model.CategoryBills:         "📄",
		model.CategoryEntertainment: "🎭",
		model.CategorySport:         "🏋️",
		model.CategoryEatingOut:     "🍽️",
		model.CategoryTransport:     "🚆",
		model.CategoryLearning:      "📚",
		model.CategoryToiletry:      "🚿",
		model.CategoryHealth:        "🏥",
		model.CategoryTech:          "💻",
		model.CategoryGifts:         "🎁",
		model.CategoryTravel:        "✈️",
		model.CategoryPets:          "🐈",
		model.CategoryOtherExpenses: "📌",
	}

	if emoji, ok := emojiMap[category]; ok {
		return emoji
	}
	return "📌" // Default emoji
}
