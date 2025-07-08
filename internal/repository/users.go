package repository

import (
	"cashout/internal/model"
	"errors"
	"strings"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"gorm.io/gorm"
)

type Users struct {
	Repository
}

func (r *Users) UpsertWithContext(user gotgbot.User, session model.UserSession) error {
	name := user.FirstName
	name = strings.Trim(name, " ")
	if name == "" {
		name = user.Username
	}

	return r.DB.SetUser(&model.User{
		TgID:        user.Id,
		Name:        name,
		Session:     session,
		TgUsername:  user.Username,
		TgFirstname: user.FirstName,
		TgLastname:  user.LastName,
	})
}

func (u *Users) Update(user *model.User) error {
	return u.DB.SetUser(user)
}

// GetByUsername retrieves a user by their Telegram username
func (r *Users) GetByUsername(username string) (model.User, bool, error) {
	user, err := r.DB.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, false, nil
		}
		return model.User{}, false, err
	}
	return *user, true, nil
}

// GetByTgID retrieves a user by their Telegram ID
func (r *Users) GetByTgID(tgID int64) (model.User, error) {
	user, err := r.DB.GetUser(tgID)
	if err != nil {
		return model.User{}, err
	}
	return *user, nil
}
