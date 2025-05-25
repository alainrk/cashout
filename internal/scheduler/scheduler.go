package scheduler

import (
	"cashout/internal/client"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

// TODO: Move into config
const WEEKLY_REMINDER_PROCESSING_MIN = 1

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
	// Schedule the creation of weekly recaps "reminders", every day, just to be sure
	// s.scheduler.Every(1).Minute().Do(func() { /* TEST */
	s.scheduler.Every(1).Day().At("05:00").Do(func() {
		if err := s.createWeeklyReminders(); err != nil {
			s.logger.Errorf("Failed to create weekly reminders: %v", err)
		}
	})

	// Process weekly reminders every minute
	s.scheduler.Every(WEEKLY_REMINDER_PROCESSING_MIN).Minute().Do(func() {
		if err := s.processWeeklyReminders(); err != nil {
			s.logger.Errorf("Failed to process weekly reminders: %v", err)
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
