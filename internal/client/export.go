package client

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// ExportTransactions handles the /export command
func (c *Client) ExportTransactions(b *gotgbot.Bot, ctx *ext.Context) error {
	_, u := c.getUserFromContext(ctx)
	user, err := c.authAndGetUser(u)
	if err != nil {
		return err
	}

	// Get all user transactions
	transactions, err := c.Repositories.Transactions.GetUserTransactions(user.TgID)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	if len(transactions) == 0 {
		_, err = b.SendMessage(ctx.EffectiveSender.ChatId, "You don't have any transactions to export.", nil)
		return err
	}

	// Create CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"tg_id",
		"date",
		"type",
		"category",
		"amount",
		"currency",
		"description",
		"created_at",
		"updated_at",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write transactions
	for _, t := range transactions {
		record := []string{
			strconv.FormatInt(t.TgID, 10),
			t.Date.Format("2006-01-02"),
			string(t.Type),
			string(t.Category),
			fmt.Sprintf("%.2f", t.Amount),
			string(t.Currency),
			t.Description,
			t.CreatedAt.Format("2006-01-02 15:04"),
			t.UpdatedAt.Format("2006-01-02 15:04"),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("CSV writer error: %w", err)
	}

	// Generate filename with current date
	filename := fmt.Sprintf("cashout_export_%s.csv", time.Now().Format("2006-01-02"))

	// Send the CSV file
	_, err = b.SendDocument(ctx.EffectiveSender.ChatId, gotgbot.InputFileByReader(filename, bytes.NewReader(buf.Bytes())), &gotgbot.SendDocumentOpts{
		Caption:   fmt.Sprintf("📊 Exported %d transactions\n\nFile: %s", len(transactions), filename),
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to send CSV file: %w", err)
	}

	return nil
}
