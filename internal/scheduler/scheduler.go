package scheduler

import (
	"cashout/internal/client"
	"time"

	gotgbot "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

const (
	WEEKLY_REMINDER_PROCESSING_MIN  = 60
	MONTHLY_REMINDER_PROCESSING_MIN = 60
)

type Scheduler struct {
	scheduler    *gocron.Scheduler
	bot          *gotgbot.Bot
	repositories client.Repositories
	logger       *logrus.Logger
}

func NewScheduler(bot *gotgbot.Bot, repos client.Repositories, logger *logrus.Logger) *Scheduler {
	// Create scheduler with UTC timezone
	s := gocron.NewScheduler(time.UTC)

	return &Scheduler{
		scheduler:    s,
		bot:          bot,
		repositories: repos,
		logger:       logger,
	}
}

func (s *Scheduler) Start() {
	// Schedule the creation of weekly recaps
	// s.scheduler.Every(1).Minute().Do(func() { /* TEST */
	s.scheduler.Every(1).Day().At("15:00").Do(func() {
		if err := s.createWeeklyReminders(); err != nil {
			s.logger.Errorf("Failed to create weekly reminders: %v", err)
		}
	})

	// Schedule the creation of monthly recaps
	// s.scheduler.Every(1).Minute().Do(func() { /* TEST */
	s.scheduler.Every(1).Day().At("10:00").Do(func() {
		if err := s.createMonthlyReminders(); err != nil {
			s.logger.Errorf("Failed to create monthly reminders: %v", err)
		}
	})

	// Process weekly reminders
	s.scheduler.Every(WEEKLY_REMINDER_PROCESSING_MIN).Minute().Do(func() {
		if err := s.processWeeklyReminders(); err != nil {
			s.logger.Errorf("Failed to process weekly reminders: %v", err)
		}
	})

	// Process monthly reminders
	s.scheduler.Every(MONTHLY_REMINDER_PROCESSING_MIN).Minute().Do(func() {
		if err := s.processMonthlyReminders(); err != nil {
			s.logger.Errorf("Failed to process monthly reminders: %v", err)
		}
	})

	// Start the scheduler
	s.scheduler.StartAsync()
	s.logger.Info("Scheduler started successfully")
}

func (s *Scheduler) Stop() {
	s.scheduler.Stop()
	s.logger.Info("Scheduler stopped")
}
