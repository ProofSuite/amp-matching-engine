package services

import (
	"math"
	"time"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"gopkg.in/mgo.v2/bson"

	"github.com/ethereum/go-ethereum/common"
)

type OHLCVService struct {
	tradeDao interfaces.TradeDao
}

func NewOHLCVService(TradeDao interfaces.TradeDao) *OHLCVService {
	return &OHLCVService{TradeDao}
}

// Unsubscribe handles all the unsubscription messages for ticks corresponding to a pair
func (s *OHLCVService) Unsubscribe(conn *ws.Conn, bt, qt common.Address, params *types.Params) {
	id := utils.GetOHLCVChannelID(bt, qt, params.Units, params.Duration)
	ws.GetTradeSocket().Unsubscribe(id, conn)
}

// Subscribe handles all the subscription messages for ticks corresponding to a pair
// It calls the corresponding channel's subscription method and sends trade history back on the connection
func (s *OHLCVService) Subscribe(conn *ws.Conn, bt, qt common.Address, params *types.Params) {

	socket := ws.GetOHLCVSocket()

	ohlcv, err := s.GetOHLCV([]types.PairSubDoc{types.PairSubDoc{BaseToken: bt, QuoteToken: qt}},
		params.Duration,
		params.Units,
		params.From,
		params.To,
	)

	if err != nil {
		socket.SendErrorMessage(conn, err.Error())
	}

	id := utils.GetOHLCVChannelID(bt, qt, params.Units, params.Duration)
	err = socket.Subscribe(id, conn)
	if err != nil {
		message := map[string]string{
			"Code":    "UNABLE_TO_SUBSCRIBE",
			"Message": "UNABLE_TO_SUBSCRIBE: " + err.Error(),
		}

		socket.SendErrorMessage(conn, message)
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeHandler(id))
	socket.SendInitMessage(conn, ohlcv)
}

// GetOHLCV fetches OHLCV data using
// pairName: can be "" for fetching data for all pairs
// duration: in integer
// unit: sec,min,hour,day,week,month,yr
// timeInterval: 0-2 entries (0 argument: latest data,1st argument: from timestamp, 2nd argument: to timestamp)
func (s *OHLCVService) GetOHLCV(pairs []types.PairSubDoc, duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error) {
	match := make(bson.M)
	addFields := make(bson.M)
	resp := make([]*types.Tick, 0)

	currentTs := time.Now().UnixNano() / int64(time.Second)
	sort := bson.M{"$sort": bson.M{"createdAt": 1}}
	toDecimal := bson.M{"$addFields": bson.M{
		"pd": bson.M{"$toDecimal": "$price"},
		"ad": bson.M{"$toDecimal": "$amount"},
	}}

	modTime, intervalSeconds := getModTime(currentTs, duration, unit)
	group, addFields := getGroupAddFieldBson("$createdAt", unit, duration)

	lt := time.Unix(currentTs, 0)
	gt := time.Unix(modTime-intervalSeconds, 0)

	if len(timeInterval) >= 1 {
		lt = time.Unix(timeInterval[1], 0)
		gt = time.Unix(timeInterval[0], 0)
	}
	match = getMatchQuery(lt, gt, pairs...)

	match = bson.M{"$match": match}
	group = bson.M{"$group": group}
	query := []bson.M{match, sort, toDecimal, group, addFields, {"$sort": bson.M{"ts": 1}}}
	resp, err := s.tradeDao.Aggregate(query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getMatchQuery(lt, gt time.Time, pairs ...types.PairSubDoc) bson.M {

	match := bson.M{"createdAt": bson.M{"$gte": gt, "$lt": lt}}

	if len(pairs) >= 1 {
		or := make([]bson.M, 0)

		for _, pair := range pairs {
			or = append(or, bson.M{
				"$and": []bson.M{
					{
						"baseToken":  pair.BaseToken.Hex(),
						"quoteToken": pair.QuoteToken.Hex(),
					},
				},
			},
			)
		}

		match["$or"] = or
	}
	return match
}

func getModTime(ts, duration int64, unit string) (int64, int64) {
	var modTime, interval int64
	switch unit {
	case "sec":
		interval = duration
		modTime = ts - int64(math.Mod(float64(ts), float64(interval)))

	case "hour":
		interval = duration * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(interval)))

	case "day":
		interval = duration * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(interval)))

	case "week":
		interval = duration * 7 * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(interval)))

	case "month":
		d := time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.UTC).Day()
		interval = duration * int64(d) * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(interval)))

	case "year":
		// Number of days in current year
		d := time.Date(time.Now().Year()+1, 1, 1, 0, 0, 0, 0, time.UTC).Sub(time.Date(time.Now().Year(), 0, 0, 0, 0, 0, 0, time.UTC)).Hours() / 24
		interval = duration * int64(d) * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(interval)))

	case "min":
		interval = duration * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(interval)))
	}
	return modTime, interval
}

// query for grouping of the documents and addition of required fields using aggregate pipeline
func getGroupAddFieldBson(key, units string, duration int64) (bson.M, bson.M) {

	var group, addFields bson.M

	t := time.Unix(0, 0)
	var d interface{}
	if key == "now" {
		d = time.Now()
	} else {
		d = key
	}

	decimal1, _ := bson.ParseDecimal128("1")
	group = bson.M{
		"count": bson.M{"$sum": decimal1},
		"h":     bson.M{"$max": "$pd"},
		"l":     bson.M{"$min": "$pd"},
		"o":     bson.M{"$first": "$pd"},
		"c":     bson.M{"$last": "$pd"},
		"v":     bson.M{"$sum": "$ad"},
	}

	gID := make(bson.M)
	switch units {
	case "sec":
		gID = bson.M{
			"year":   bson.M{"$year": d},
			"day":    bson.M{"$dayOfMonth": d},
			"month":  bson.M{"$month": d},
			"hour":   bson.M{"$hour": d},
			"minute": bson.M{"$minute": d},
			"second": bson.M{
				"$subtract": []interface{}{
					bson.M{"$second": d},
					bson.M{"$mod": []interface{}{bson.M{"$second": d}, duration}}}}}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":   "$_id.year",
			"month":  "$_id.month",
			"day":    "$_id.day",
			"hour":   "$_id.hour",
			"minute": "$_id.minute",
			"second": "$_id.second"}}, t}}}}

	case "min":
		gID = bson.M{
			"year":  bson.M{"$year": d},
			"day":   bson.M{"$dayOfMonth": d},
			"month": bson.M{"$month": d},
			"hour":  bson.M{"$hour": d},
			"minute": bson.M{
				"$subtract": []interface{}{
					bson.M{"$minute": d},
					bson.M{"$mod": []interface{}{bson.M{"$minute": d}, duration}},
				}}}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":   "$_id.year",
			"month":  "$_id.month",
			"day":    "$_id.day",
			"hour":   "$_id.hour",
			"minute": "$_id.minute",
		}}, t}}}}

	case "hour":
		gID = bson.M{
			"year":  bson.M{"$year": d},
			"day":   bson.M{"$dayOfMonth": d},
			"month": bson.M{"$month": d},
			"hour": bson.M{
				"$subtract": []interface{}{
					bson.M{"$hour": d},
					bson.M{"$mod": []interface{}{bson.M{"$hour": d}, duration}}}}}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":  "$_id.year",
			"month": "$_id.month",
			"day":   "$_id.day",
			"hour":  "$_id.hour",
		}}, t}}}}

	case "day":
		gID = bson.M{
			"year":  bson.M{"$year": d},
			"month": bson.M{"$month": d},
			"day": bson.M{
				"$subtract": []interface{}{
					bson.M{"$dayOfMonth": d},
					bson.M{"$mod": []interface{}{bson.M{"$dayOfMonth": d}, duration}}}}}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":  "$_id.year",
			"month": "$_id.month",
			"day":   "$_id.day",
		}}, t}}}}

	case "week":
		gID = bson.M{
			"year": bson.M{"$isoWeekYear": d},
			"isoWeek": bson.M{
				"$subtract": []interface{}{
					bson.M{"$isoWeek": d},
					bson.M{"$mod": []interface{}{bson.M{"$isoWeek": d}, duration}}}}}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"isoWeekYear": "$_id.year",
			"isoWeek":     "$_id.isoWeek",
		}}, t}}}}

	case "month":
		gID = bson.M{
			"year": bson.M{"$year": d},
			"month": bson.M{
				"$subtract": []interface{}{
					bson.M{
						"$multiply": []interface{}{
							bson.M{"$ceil": bson.M{"$divide": []interface{}{
								bson.M{"$month": d},
								duration}},
							},
							duration},
					}, duration - 1}}}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":  "$_id.year",
			"month": "$_id.month",
		}}, t}}}}

	case "year":
		gID = bson.M{
			"year": bson.M{
				"$subtract": []interface{}{
					bson.M{"$year": d},
					bson.M{"$mod": []interface{}{bson.M{"$year": d}, duration}}}}}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year": "$_id.year"}}, t}}}}

	}

	gID["pair"] = "$pairName"
	gID["baseToken"] = "$baseToken"
	gID["quoteToken"] = "$quoteToken"
	group["_id"] = gID

	return group, addFields
}
