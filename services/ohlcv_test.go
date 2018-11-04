package services

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

const timeLayoutString = "Jan 2 2006 15:04:05"

type TickSorter []*types.Tick

func (a TickSorter) Len() int           { return len(a) }
func (a TickSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TickSorter) Less(i, j int) bool { return a[i].Ts < a[j].Ts }

var durations = map[string][]int64{
	"year":  {1},
	"month": {1, 3, 6, 9},
	"week":  {1},
	"day":   {1},
	"hour":  {1, 6, 12},
	"min":   {1, 5, 15, 30},
	"sec":   {15, 30},
}
var testTimes = []string{
	"Dec 17 2017 00:00:00",
	"Dec 23 2017 23:59:59",
	"Dec 24 2017 00:00:00",
	"Dec 30 2017 23:59:59",
	"Dec 31 2017 23:59:58",

	"Jan 1 2018 00:00:00",
	"Jan 31 2018 23:59:59",

	"Apr 1 2018 00:00:00",
	"Apr 30 2018 23:59:59",

	"Jun 1 2018 00:00:00",
	"Jun 30 2018 23:59:59",

	"Aug 21 2018 08:00:00",
	"Aug 21 2018 08:00:14",
	"Aug 21 2018 08:00:15",
	"Aug 21 2018 08:00:29",
	"Aug 21 2018 08:00:30",
	"Aug 21 2018 08:00:44",
	"Aug 21 2018 08:00:45",
	"Aug 21 2018 08:00:59",
	"Aug 21 2018 08:01:00",
	"Aug 21 2018 08:04:59",
	"Aug 21 2018 08:05:00",
	"Aug 21 2018 08:14:59",
	"Aug 21 2018 08:15:00",
	"Aug 21 2018 08:29:59",
	"Aug 21 2018 08:30:00",
	"Aug 21 2018 08:59:59",
	"Aug 21 2018 09:00:00",
	"Aug 21 2018 13:59:59",
	"Aug 21 2018 14:00:00",
	"Aug 21 2018 19:59:59",
	"Aug 21 2018 20:00:00",
	"Aug 21 2018 23:59:59",
	"Aug 22 2018 00:00:00",
}

var durationMap = make(map[string]map[int64]*types.Tick)

func TestOHLCV(t *testing.T) {
	pair := types.PairAddresses{
		Name:       "HPC/AUT",
		BaseToken:  common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken: common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
	}

	sampleTrade := types.Trade{
		Taker:          common.HexToAddress("0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63"),
		Maker:          common.HexToAddress("0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b58"),
		BaseToken:      pair.BaseToken,
		QuoteToken:     pair.QuoteToken,
		MakerOrderHash: common.HexToHash("0x4ac68946450e5a6273b92d81aa58f288d7b5515942456b89fb5c7e982efead7c"),
		TakerOrderHash: common.HexToHash("0x4ac68946450e5a6273b92d81aa58f288d7b5515942456b89fb5c7e982efead7c"),
		Hash:           common.HexToHash("0x4ac68946450e5a6273b92d81aa58f288d7b5515942456b89fb5c7e982efeas3f"),
		PairName:       pair.Name,
		Signature:      &types.Signature{},
		Side:           "BUY",
		PricePoint:     big.NewInt(9987),
		Amount:         big.NewInt(125772),
	}
	app.Config.DBName = "proofdex"
	tradeDao := daos.NewTradeDao()
	ohlcvService := NewOHLCVService(tradeDao)

	for _, t := range testTimes {
		tTime, err := time.Parse(timeLayoutString, t)
		if err != nil {
			panic("invalid date: " + err.Error())
		}
		amt := new(big.Int)
		prc := new(big.Int)
		sampleTrade.CreatedAt = tTime
		sampleTrade.Amount = amt.Add(sampleTrade.Amount, big.NewInt(10))
		sampleTrade.PricePoint = prc.Add(sampleTrade.PricePoint, big.NewInt(5))
		sampleTrade.ID = bson.NewObjectId()
		sampleTrade.Hash = sampleTrade.ComputeHash()

		if err := db.DB(app.Config.DBName).C("trades").Insert(&sampleTrade); err != nil {
			panic(err)
		}

		if err := updateExpectedResponse(&sampleTrade); err != nil {
			panic(err)
		}

	}

	for unit, durationSlice := range durations {
		for _, duration := range durationSlice {
			response, err := ohlcvService.GetOHLCV([]types.PairAddresses{pair}, duration, unit, 0, time.Now().Unix())
			if err != nil {
				t.Errorf("%s", err)
				return
			}
			expectedResponse := getExpectedResponse(unit, duration)
			rab, _ := json.Marshal(response)
			erab, _ := json.Marshal(expectedResponse)

			assert.JSONEq(t, string(erab), string(rab))
		}
	}
}

func updateExpectedResponse(trade *types.Trade) error {

	tradeTs := trade.CreatedAt.Unix()
	for unit, durationSlice := range durations {
		for _, duration := range durationSlice {
			var ts int64
			switch unit {
			case "sec":
				intervalSeconds := duration
				ts = tradeTs - int64(math.Mod(float64(tradeTs), float64(intervalSeconds)))

			case "hour":
				intervalSeconds := duration * 60 * 60
				ts = tradeTs - int64(math.Mod(float64(tradeTs), float64(intervalSeconds)))

			case "day":
				intervalSeconds := duration * 24 * 60 * 60
				ts = tradeTs - int64(math.Mod(float64(tradeTs), float64(intervalSeconds)))

			case "week":
				yr, wk := trade.CreatedAt.ISOWeek()
				ts = firstDayOfISOWeek(yr, wk, time.UTC).Unix()

			case "month":
				mnth, _ := strconv.ParseInt(fmt.Sprintf("%d", trade.CreatedAt.Month()), 10, 64)
				ts = time.Date(trade.CreatedAt.Year(), time.Month((math.Ceil(float64(mnth)/float64(duration))*float64(duration))-float64(duration-1)), 1, 0, 0, 0, 0, time.UTC).Unix()

			case "year":
				ts = time.Date(trade.CreatedAt.Year(), 1, 1, 0, 0, 0, 0, time.UTC).Unix()

			case "min":
				intervalSeconds := duration * 60
				ts = tradeTs - int64(math.Mod(float64(tradeTs), float64(intervalSeconds)))

			default:
				panic("Invalid unit: " + unit)
			}

			addTick(unit, duration, ts, trade)

		}
	}
	return nil
}

func addTick(unit string, duration int64, ts int64, trade *types.Trade) {
	key := getKey(unit, duration)
	if durationMap[key] == nil {
		durationMap[key] = make(map[int64]*types.Tick)
	}
	durationMap[key][ts] = tradeToTick(trade, durationMap[key][ts], ts)
}

func tradeToTick(trade *types.Trade, tick *types.Tick, ts int64) *types.Tick {
	if tick == nil {
		tick = &types.Tick{
			ID: types.TickID{
				Pair:       trade.PairName,
				BaseToken:  trade.BaseToken,
				QuoteToken: trade.QuoteToken,
			},
			O:     trade.PricePoint,
			H:     trade.PricePoint,
			L:     trade.PricePoint,
			C:     trade.PricePoint,
			V:     trade.Amount,
			Count: big.NewInt(1),
			Ts:    ts * 1000,
		}
	} else {
		tick.C = trade.PricePoint
		tv := new(big.Int)
		tv.Add(tick.V, trade.Amount)
		tick.V = tv

		tick.Count.Add(tick.Count, big.NewInt(1))
		if trade.PricePoint.Cmp(tick.H) == 1 {
			tick.H = trade.PricePoint
		}
		if trade.PricePoint.Cmp(tick.L) == -1 {
			tick.L = trade.PricePoint
		}
	}
	return tick
}
func getKey(unit string, duration int64) string {
	return fmt.Sprintf("%d%s", duration, unit)
}
func getExpectedResponse(unit string, duration int64) (ticks []*types.Tick) {
	k := getKey(unit, duration)
	ticks = make([]*types.Tick, 0)
	if durationMap[k] != nil {
		for _, tick := range durationMap[k] {
			ticks = append(ticks, tick)
		}
	}
	sort.Sort(TickSorter(ticks))
	return
}

func firstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}
