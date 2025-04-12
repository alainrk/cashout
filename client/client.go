package client

import (
	"happypoor/internal/db"
	"happypoor/internal/model"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Client struct {
	Db *db.DB
}

func (c *Client) getUserData(_ *ext.Context, username string) (model.User, bool, error) {
	user, err := c.Db.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "record not found" {
			return model.User{}, false, nil
		}
		return model.User{}, false, err
	}
	return *user, true, nil
}

func (c *Client) setUserData(ctx *ext.Context, session model.UserSession) error {
	name := ctx.Message.From.FirstName
	name = strings.Trim(name, " ")
	if name == "" {
		name = ctx.Message.From.Username
	}

	return c.Db.SetUser(&model.User{
		TgID:        ctx.Message.From.Id,
		Name:        name,
		Session:     session,
		TgUsername:  ctx.Message.From.Username,
		TgFirstname: ctx.Message.From.FirstName,
		TgLastname:  ctx.Message.From.LastName,
	})
}
