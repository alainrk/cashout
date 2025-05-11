package repository

import (
	"cashout/internal/db"

	"github.com/sirupsen/logrus"
)

type Repository struct {
	DB     *db.DB
	Logger *logrus.Logger
}
