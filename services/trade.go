package services

import (
	"fmt"
	"math"
	"time"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
	"labix.org/v2/mgo/bson"
)

type TradeService struct {
	tradeDao *daos.TradeDao
}

func NewTradeService(TradeDao *daos.TradeDao) *TradeService {
	return &TradeService{TradeDao}
}

func (t *TradeService) GetByPairName(pairName string) ([]*types.Trade, error) {
	return t.tradeDao.GetByPairName(pairName)
}

func (t *TradeService) GetByUserAddress(addr string) ([]*types.Trade, error) {
	return t.tradeDao.GetByUserAddress(addr)
}

// func (t *TradeService) GetMarketStats(pairName ...string) (resp map[string]interface{}, err error) {
// 	var match bson.M
// 	if len(pairName) == 0 {
// 		match = bson.M{"$match": bson.M{"createdAt": bson.M{"$gte": time.Now().Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)}}}
// 	} else {
// 		match = bson.M{"$match": bson.M{"createdAt": bson.M{"$gte": time.Now().Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)}, "pairName": bson.M{"$in": pairName}}}
// 	}
// 	fmt.Sprintf("%s", match)
// 	sort := bson.M{
// 		"$sort": bson.M{"createdAt": -1},
// 	}
// 	group := bson.M{
// 		"$group": bson.M{
// 			"_id":              bson.M{"market": "$market"},
// 			"volume":           bson.M{"$sum": "$tradeSize"},
// 			"hourlyHigh":       bson.M{"$max": "$tradePrice"},
// 			"hourlyLow":        bson.M{"$min": "$tradePrice"},
// 			"lastPrice":        bson.M{"$last": "$tradePrice"},
// 			"firstPrice":       bson.M{"$first": "$tradePrice"},
// 			"baseCurrency":     bson.M{"$first": "$baseCurrency"},
// 			"exchangeCurrency": bson.M{"$first": "$exchangeCurrency"},
// 			"market":           bson.M{"$first": "$market"},
// 		},
// 	}
// 	addFields := bson.M{
// 		"$addFields": bson.M{
// 			"change": bson.M{
// 				"$multiply": []interface{}{
// 					bson.M{"$divide": []interface{}{
// 						bson.M{"$subtract": []interface{}{
// 							"$lastPrice",
// 							"$firstPrice",
// 						},
// 						}, "$firstPrice"},
// 					}, 100},
// 			},
// 		},
// 	}

// 	_, err := t.tradeDao.Aggregate([]bson.M{match, sort, group, addFields}) //dao.db.DB(dao.dbName).C(dao.collectionName).Pipe([]bson.M{match, sort, group, addFields}).All(&resp)
// 	return
// }
func (t *TradeService) GetTicks(pairName string, duration int64, unit string, timeInterval ...int64) (resp []interface{}, nextTime int64, err error) {
	var match bson.M
	currentTs := time.Now().UnixNano() / int64(time.Second)
	var lt time.Time
	var gt time.Time
	sort := bson.M{"$sort": bson.M{"createdAt": 1}}
	group := bson.M{
		"count": bson.M{"$sum": 1},
		"h":     bson.M{"$max": "$price"},
		"l":     bson.M{"$min": "$price"},
		"o":     bson.M{"$first": "$price"},
		"c":     bson.M{"$last": "$price"},
		"v":     bson.M{"$sum": "$amount"},
	}
	var addFields bson.M
	var intervalSeconds int64
	var modTime int64
	switch unit {

	case "sec":
		group["_id"], addFields = getGroupTsBson("$createdAt", "sec", duration)
		intervalSeconds = duration
		modTime = currentTs - int64(math.Mod(float64(currentTs), float64(intervalSeconds)))

	case "hour":
		group["_id"], addFields = getGroupTsBson("$createdAt", "hour", duration)
		intervalSeconds = duration * 60 * 60
		modTime = currentTs - int64(math.Mod(float64(currentTs), float64(intervalSeconds)))

	case "day":
		group["_id"], addFields = getGroupTsBson("$createdAt", "day", duration)
		intervalSeconds = duration * 24 * 60 * 60
		modTime = currentTs - int64(math.Mod(float64(currentTs), float64(intervalSeconds)))

	case "week":
		group["_id"], addFields = getGroupTsBson("$createdAt", "week", duration)
		intervalSeconds = duration * 7 * 24 * 60 * 60
		modTime = currentTs - int64(math.Mod(float64(currentTs), float64(intervalSeconds)))

	case "month":
		group["_id"], addFields = getGroupTsBson("$createdAt", "month", duration)
		d := time.Date(time.Now().Year(), time.Now().Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		intervalSeconds = duration * int64(d) * 7 * 24 * 60 * 60
		modTime = currentTs - int64(math.Mod(float64(currentTs), float64(intervalSeconds)))

	case "yr":
		group["_id"], addFields = getGroupTsBson("$createdAt", "yr", duration)
		// Number of days in current year
		d := time.Date(time.Now().Year()+1, 0, 0, 0, 0, 0, 0, time.UTC).Sub(time.Date(time.Now().Year(), 0, 0, 0, 0, 0, 0, time.UTC)).Hours() / 24

		intervalSeconds = duration * int64(d) * 7 * 24 * 60 * 60
		modTime = currentTs - int64(math.Mod(float64(currentTs), float64(intervalSeconds)))

	case "":
	case "min":
		group["_id"], addFields = getGroupTsBson("$createdAt", "min", duration)
		intervalSeconds = duration * 60
		modTime = currentTs - int64(math.Mod(float64(currentTs), float64(intervalSeconds)))

	default:
		err = fmt.Errorf("Invalid unit please try again")
		return
	}
	lt = time.Unix(modTime, 0)
	gt = time.Unix(modTime-intervalSeconds, 0)

	if err != nil {
		err = fmt.Errorf("Invalid units value. " + unit + " Please try again")
		return
	}
	if len(timeInterval) == 0 {
		match = bson.M{"$match": bson.M{"pairName": pairName, "createdAt": bson.M{"$lt": lt}}}
	} else if len(timeInterval) >= 1 {
		lt = time.Unix(timeInterval[1], 0)
		gt = time.Unix(timeInterval[0], 0)
		match = bson.M{"$match": bson.M{"pairName": pairName, "createdAt": bson.M{"$gte": gt, "$lt": lt}}}
	}
	group = bson.M{"$group": group}
	query := []bson.M{match, sort, group, addFields, bson.M{"$sort": bson.M{"ts": 1}}}
	resp, err = t.tradeDao.Aggregate(query) // dao.db.DB(dao.dbName).C(dao.collectionName).Pipe(query).All(&resp)
	if err != nil {
		return
	}
	return
}

func getGroupTsBson(key, units string, duration int64) (resp bson.M, addFields bson.M) {
	t := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	var d interface{}
	if key == "now" {
		d = time.Now()
	} else {
		d = key
	}
	if units == "sec" {

		resp = bson.M{
			"year":   bson.M{"$year": d},
			"day":    bson.M{"$dayOfMonth": d},
			"month":  bson.M{"$month": d},
			"hour":   bson.M{"$hour": d},
			"minute": bson.M{"$minute": d},
			"second": bson.M{
				"$subtract": []interface{}{
					bson.M{"$second": d},
					bson.M{"$mod": []interface{}{bson.M{"$second": d}, duration}},
				},
			},
		}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":   "$_id.year",
			"month":  "$_id.month",
			"day":    "$_id.day",
			"hour":   "$_id.hour",
			"minute": "$_id.minute",
			"second": "$_id.second",
		},
		}, t,
		},
		},
		},
		}
	} else if units == "min" {

		resp = bson.M{
			"year":  bson.M{"$year": d},
			"day":   bson.M{"$dayOfMonth": d},
			"month": bson.M{"$month": d},
			"hour":  bson.M{"$hour": d},
			"minute": bson.M{
				"$subtract": []interface{}{
					bson.M{"$minute": d},
					bson.M{"$mod": []interface{}{bson.M{"$minute": d}, duration}},
				},
			},
		}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":   "$_id.year",
			"month":  "$_id.month",
			"day":    "$_id.day",
			"hour":   "$_id.hour",
			"minute": "$_id.minute",
		}}, t}}}}

	} else if units == "hour" {

		resp = bson.M{
			"year":  bson.M{"$year": d},
			"day":   bson.M{"$dayOfMonth": d},
			"month": bson.M{"$month": d},
			"hour": bson.M{
				"$subtract": []interface{}{
					bson.M{"$hour": d},
					bson.M{"$mod": []interface{}{bson.M{"$hour": d}, duration}},
				},
			},
		}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":  "$_id.year",
			"month": "$_id.month",
			"day":   "$_id.day",
			"hour":  "$_id.hour",
		}}, t}}}}

	} else if units == "day" {

		resp = bson.M{
			"year":  bson.M{"$year": d},
			"month": bson.M{"$month": d},
			"day": bson.M{
				"$subtract": []interface{}{
					bson.M{"$dayOfMonth": d},
					bson.M{"$mod": []interface{}{bson.M{"$dayOfMonth": d}, duration}},
				},
			},
		}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":  "$_id.year",
			"month": "$_id.month",
			"day":   "$_id.day",
		}}, t}}}}

	} else if units == "week" {

		resp = bson.M{
			"year": bson.M{"$year": d},
			"isoWeek": bson.M{
				"$subtract": []interface{}{
					bson.M{"$isoWeek": d},
					bson.M{"$mod": []interface{}{bson.M{"$isoWeek": d}, duration}},
				},
			},
		}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":    "$_id.year",
			"isoWeek": "$_id.isoWeek",
		}}, t}}}}

	} else if units == "month" {

		resp = bson.M{
			"year": bson.M{"$year": d},
			"month": bson.M{
				"$subtract": []interface{}{
					bson.M{"$month": d},
					bson.M{"$mod": []interface{}{bson.M{"$month": d}, duration}},
				},
			},
		}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":  "$_id.year",
			"month": "$_id.month",
		}}, t}}}}
	} else if units == "yr" {

		resp = bson.M{
			"year": bson.M{
				"$subtract": []interface{}{
					bson.M{"$year": d},
					bson.M{"$mod": []interface{}{bson.M{"$year": d}, duration}},
				},
			},
		}

		addFields = bson.M{"$addFields": bson.M{"ts": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year": "$_id.year"}}, t}}}}
	}
	return
}
