package crons

import (
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/robfig/cron"
)

// CronService contains the services required to initialize crons
type CronService struct {
	ohlcvService *services.OHLCVService
}

// NewCronService returns a new instance of CronService
func NewCronService(ohlcvService *services.OHLCVService) *CronService {
	return &CronService{ohlcvService}
}

// InitCrons is responsible for initializing all the crons in the system
func (s *CronService) InitCrons() {
	c := cron.New()
	s.tickStreamingCron(c)
	c.Start()
}
