package provider

import (
	// "fmt"
	"github.com/qshuai162/qian/api/chbtc"
	"github.com/qshuai162/qian/api/huobi"
	"github.com/qshuai162/qian/api/okcoin"
	"github.com/qshuai162/qian/api/btcc"
	"github.com/qshuai162/qian/common/logger"
	. "github.com/qshuai162/qian/common/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MarketAPI interface {
	GetKLine(peroid int, size int) (ret bool, records []Record)
	GetTicker() (ret bool, m Ticker)
	GetOrders() (ret bool, orders []Order)
	GetAccount() (account Account, ret bool)
}

type Provider struct {
}

var dbstr string = "121.41.46.25:27017"

func getMarketAPI(center string) (marketAPI MarketAPI) {
	if center == "huobi" {
		marketAPI = huobi.NewHuobi()
	} else if center == "okcoin" {
		marketAPI = okcoin.NewOkcoin()
	} else if center == "chbtc" {
		marketAPI = chbtc.NewCHBTC()
	} else if center == "btcc" {
		marketAPI = btcc.NewBTCC()
	} else {
		logger.Fatalln("Please config the market center...")
		panic(-1)
	}
	return
}

func (p *Provider) GetTicker(center string, t int64) (code int, result Ticker) {
	_, result = getMarketAPI(center).GetTicker()
	return
}

//差值平均值
func (p *Provider) GetDifAvg(center1, center2 string, start, long int64) (int, float64) {
	session, err := mgo.Dial(dbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	ok := session.DB(center1).C("ticker")
	hb := session.DB(center2).C("ticker")
	var hb_results []Ticker
	hb.Find(bson.M{"time": bson.M{"$gt": start - long}}).All(&hb_results)
	var ok_results []Ticker
	ok.Find(bson.M{"time": bson.M{"$gt": start - long}}).All(&ok_results)
	return 0, calcAvg(&ok_results, &hb_results, start-long+1, long)
}

func calcAvg(ok_results, hb_results *[]Ticker, start int64, long int64) float64 {

	m1 := make(map[int64]float64)
	for i := 0; i < len(*ok_results); i++ {
		m1[(*ok_results)[i].Time] = (*ok_results)[i].Last
	}
	m2 := make(map[int64]float64)
	for i := 0; i < len(*hb_results); i++ {
		m2[(*hb_results)[i].Time] = (*hb_results)[i].Last
	}
	sum := 0.0
	count := 0
	for i := start; i < start+long; i++ {
		v1, ok1 := m1[i]
		v2, ok2 := m2[i]
		if ok1 && ok2 && v1 != 0 && v2 != 0 {
			offset := v1 - v2
			sum += offset
			count++
		}
	}
	avg := sum / (float64(count))
	logger.Infof("Offset AVG: %.2f\n", avg)
	return avg
}
