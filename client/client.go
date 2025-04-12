package client

import (
	"happypoor/internal/ai"
	"happypoor/internal/db"
	"happypoor/internal/model"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Client struct {
	DB  *db.DB
	LLM ai.LLM
}

func (c *Client) getUserData(_ *ext.Context, username string) (model.User, bool, error) {
	user, err := c.DB.GetUserByUsername(username)
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

	return c.DB.SetUser(&model.User{
		TgID:        ctx.Message.From.Id,
		Name:        name,
		Session:     session,
		TgUsername:  ctx.Message.From.Username,
		TgFirstname: ctx.Message.From.FirstName,
		TgLastname:  ctx.Message.From.LastName,
	})
}
