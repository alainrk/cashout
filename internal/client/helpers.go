package client

import (
	"happypoor/internal/model"
)

// getCategoryEmoji returns the appropriate emoji for a transaction category
func getCategoryEmoji(category model.TransactionCategory) string {
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
	}

	if emoji, ok := emojiMap[category]; ok {
		return emoji
	}
	return "📌" // Default emoji
}
