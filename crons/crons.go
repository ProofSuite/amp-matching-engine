package crons

import (
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/robfig/cron"
)

type CronService struct {
	tradeService *services.TradeService
}

func NewCronService(tradeService *services.TradeService) *CronService {
	return &CronService{tradeService}
}
func (s *CronService) InitCrons() {
	c := cron.New()
	s.tickStreamingCron(c)
	c.Start()
}
