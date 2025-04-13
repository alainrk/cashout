package client

import (
	"fmt"
	"happypoor/internal/ai"
	"happypoor/internal/db"
	"happypoor/internal/model"
	"happypoor/internal/repository"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Client struct {
	Repositories Repositories
	LLM          ai.LLM
}

type Repositories struct {
	Users repository.Users
}

func NewClient(db *db.DB, llm ai.LLM) *Client {
	return &Client{
		Repositories: Repositories{
			Users: repository.Users{DB: db},
		},
		LLM: llm,
	}
}

// authAndGetUser authenticates the user and returns the user data.
func (c *Client) authAndGetUser(ctx *ext.Context) (model.User, error) {
	user, exists, err := c.Repositories.Users.GetByUsername(ctx.Message.From.Username)
	if err != nil {
		return user, fmt.Errorf("failed to get user data: %w", err)
	}

	if exists {
		user.Session.Iterations++
		c.Repositories.Users.Update(&user)
		return user, nil
	}

	// User to be created
	session := model.UserSession{
		Iterations:  0,
		State:       model.StateNormal,
		LastCommand: model.CommandStart,
		LastMessage: ctx.Message.Text,
	}

	if err := c.Repositories.Users.UpsertWithContext(ctx, session); err != nil {
		return user, fmt.Errorf("failed to set user data: %w", err)
	}

	user, _, err = c.Repositories.Users.GetByUsername(ctx.Message.From.Username)
	if err != nil {
		return user, fmt.Errorf("failed to get user data: %w", err)
	}

	return user, nil
}

// Start introduces the bot.
func (c *Client) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Welcome to HappyPoor, %s!", user.Name)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}

// Message handles incoming messages not in a specific flow.
func (c *Client) Message(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	if user.Session.Iterations == 0 {
		return c.Start(b, ctx)
	}

	user.Session.State = model.StateWaiting
	// TODO: This must be understood through the analysis of the text the user have given
	user.Session.LastCommand = model.CommandExpenseAdd
	user.Session.LastMessage = ctx.Message.Text

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	expense, err := c.LLM.ExtractExpense(ctx.Message.Text)
	if err != nil {
		msg := fmt.Sprintf("I'm sorry, I can't understand your expense '%s', %s!", user.Session.LastMessage, user.Name)
		ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		return err
	}

	msg := fmt.Sprintf("%s (€ %.2f), %s", expense.Category, expense.Amount, expense.Description)
	fmt.Println(msg)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}

// Cancel returns to normal state.
func (c *Client) Cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := c.authAndGetUser(ctx)
	if err != nil {
		return err
	}

	user.Session.State = model.StateNormal
	user.Session.LastCommand = model.CommandCancel
	user.Session.LastMessage = ctx.Message.Text

	err = c.Repositories.Users.Update(&user)
	if err != nil {
		return fmt.Errorf("failed to set user data: %w", err)
	}

	msg := fmt.Sprintf("Sure, you can always restart with /start, %s!", user.Name)
	ctx.EffectiveMessage.Reply(b, msg, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})

	return nil
}
