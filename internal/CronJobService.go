package internal

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"go-fiber-starter/utils"
	"go.uber.org/fx"
)

type CronService struct {
	cron   *cron.Cron
	logger zerolog.Logger
}

func NewCronService(lc fx.Lifecycle, logger zerolog.Logger) *CronService {
	// Create a new cron job with second-level precision
	c := cron.New(cron.WithSeconds())

	// Create the service
	service := &CronService{
		cron:   c,
		logger: logger,
	}

	// Use Fx lifecycle to manage cron job
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			service.logger.Info().Msg("Starting cron service")
			c.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			service.logger.Info().Msg("Stopping cron service")

			// Stop the cron job and wait for any running jobs to complete
			stopCtx := c.Stop()
			<-stopCtx.Done()
			return nil
		},
	})

	return service
}

// AddJob is a method to add a new cron job
func (s *CronService) AddJob(spec string, cmd func()) error {
	if utils.IsChildProcess() {
		return nil
	}

	_, err := s.cron.AddFunc(spec, func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error().Msg("Cron job panic recovered")
			}
		}()
		cmd()
	})

	return err
}
