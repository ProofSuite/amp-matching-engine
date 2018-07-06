package crons

import (
	"fmt"
	"log"

	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/ws"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/robfig/cron"
)

func (s *CronService) tickStreamingCron(c *cron.Cron) {
	for unit, durations := range app.Config.TickDuration {
		for _, duration := range durations {
			schedule := getCronScheduleString(unit, duration)
			c.AddFunc(schedule, s.tickStream(unit, duration))
		}
	}
}
func (s *CronService) tickStream(unit string, duration int64) func() {
	return func() {
		// log.Printf("TickStreaming Ran: unit: %s duration: %d\n", unit, duration)
		ticks, err := s.tradeService.GetTicks("", duration, unit)
		if err != nil {
			log.Printf("%s", err)
			return
		}
		for _, tick := range ticks {
			ws.TickBroadcast(utils.GetTickChannelID(tick.ID.Pair, unit, duration), tick)
		}
	}
}

func getCronScheduleString(unit string, duration int64) string {
	switch unit {

	case "sec":
		return fmt.Sprintf("*/%d * * * * *", duration)

	case "min":
		return fmt.Sprintf("0 */%d * * * *", duration)

	case "hour":
		return fmt.Sprintf("0 0 %d * * *", duration)

	case "day":
		return fmt.Sprintf("@daily")

	case "week":
		return fmt.Sprintf("0 0 0 * * */%d", duration)

	case "month":
		return fmt.Sprintf("0 0 0 */%d * *", duration)

	case "year":
		return fmt.Sprintf("@yearly")

	default:
		panic(fmt.Errorf("Invalid unit please try again"))
	}
}
