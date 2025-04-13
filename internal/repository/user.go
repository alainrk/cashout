package repository

import (
	"happypoor/internal/db"
	"happypoor/internal/model"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Users struct {
	DB *db.DB
}

func (u *Users) GetByUsername(username string) (model.User, bool, error) {
	user, err := u.DB.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "record not found" {
			return model.User{}, false, nil
		}
		return model.User{}, false, err
	}
	return *user, true, nil
}

func (u *Users) UpsertWithContext(ctx *ext.Context, session model.UserSession) error {
	name := ctx.Message.From.FirstName
	name = strings.Trim(name, " ")
	if name == "" {
		name = ctx.Message.From.Username
	}

	return u.DB.SetUser(&model.User{
		TgID:        ctx.Message.From.Id,
		Name:        name,
		Session:     session,
		TgUsername:  ctx.Message.From.Username,
		TgFirstname: ctx.Message.From.FirstName,
		TgLastname:  ctx.Message.From.LastName,
	})
}

func (u *Users) Update(user *model.User) error {
	return u.DB.SetUser(user)
}
