package repository

import (
	"happypoor/internal/db"
	"happypoor/internal/model"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Users struct {
	DB *db.DB
}

func (r *Users) GetByUsername(username string) (model.User, bool, error) {
	user, err := r.DB.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "record not found" {
			return model.User{}, false, nil
		}
		return model.User{}, false, err
	}
	return *user, true, nil
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
