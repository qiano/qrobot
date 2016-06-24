package strategy

import (
	// . "github.com/qshuai162/qian/common"
	. "github.com/qshuai162/qian/common/config"
	// "github.com/qshuai162/qian/huobi"
	// "fmt"
	"github.com/qshuai162/qian/common/logger"
	idc "github.com/qshuai162/qian/qrobot/indictor"
	. "github.com/qshuai162/qian/qrobot/trader"
	"strconv"
	"time"
)

type MACDStrategy struct {
	PrevBuyTime      time.Time
	PrevStoplossTime time.Time
	Orders           OCaches
}

func init() {
	macd := new(MACDStrategy)
	macd.PrevBuyTime = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	macd.PrevStoplossTime = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	Register("MACD", macd)
}

// KDJ-EX strategy
func (sty *MACDStrategy) Tick() bool {
	center := Option["tradecenter"]
	market := GetMarket(center)
	trader := NewTrader(center, &sty.Orders)
	trend := idc.CalcByMACD(&market)
	if trend <= 0 {
		return false
	}
	kperoid, _ := strconv.Atoi(Option["tick_interval"])
	records := market.GetKLine(kperoid)
	orders := market.GetOrders(false)
	length := len(records)
	if length < 3 {
		return false
	}

	numTradeAmount, _ := strconv.ParseFloat(Option["tradeAmount"], 64)
	timeout, _ := strconv.Atoi(Option["timeout"])
	stoploss, _ := strconv.ParseFloat(Option["stoploss"], 64)
	stoploss_resttime, _ := strconv.Atoi(Option["stoploss_resttime"])
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])
	signalPeriod, _ := strconv.Atoi(Option["signalPeriod"])
	MACDbuyThreshold, _ := strconv.ParseFloat(Option["MACDbuyThreshold"], 64)
	// MACDsellThreshold, _ := strconv.ParseFloat(Option["MACDsellThreshold"], 64)

	logger.Infoln("策略: MACD  MACD参数:", shortEMA, longEMA, signalPeriod, "K线周期:", kperoid)
	sty.Orders.Sync(orders)

	macd := idc.GetMACD(records, shortEMA, longEMA, signalPeriod)
	logger.Infof("Prev:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-3].Close, macd[length-3].DIF, macd[length-3].DEA, macd[length-3].BAR)
	logger.Infof("LAST:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-2].Close, macd[length-2].DIF, macd[length-2].DEA, macd[length-2].BAR)
	logger.Infof("CURR:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-1].Close, macd[length-1].DIF, macd[length-1].DEA, macd[length-1].BAR)
	//买
	if (macd[length-1].BAR > MACDbuyThreshold && macd[length-2].BAR > 0 && macd[length-1].BAR > macd[length-2].BAR && macd[length-1].DEA > 0 && macd[length-1].DIF > 0) &&
		(records[length-2].Close > records[length-3].Close && records[length-1].Close > records[length-2].Close) {
		maxbar_index := length - 1
		for i := length - 2; i >= length-4; i-- {
			if macd[i].BAR < 0 {
				break
			} else {
				if macd[i].BAR > macd[maxbar_index].BAR {
					maxbar_index = i
				}
			}
		}
		if maxbar_index == length-1 {
			logger.Infoln("--->KDJ up cross,准备买入, macd:", macd[length-1].BAR)
			if time.Now().Sub(sty.PrevStoplossTime) >= time.Duration(stoploss_resttime)*time.Minute {
				if len(sty.Orders.Buy_NotSell()) < 20 && time.Now().Sub(sty.PrevBuyTime) >= time.Duration(kperoid*8)*time.Minute {
					amount := numTradeAmount * trend
					if amount > 0 {
						if trader.BuyGroup(amount) {
							sty.PrevBuyTime = time.Now()
						}
					}
				}
			} else {
				logger.Infoln("止损保护中，不允许交易")
			}
		}
	}

	//卖
	if macd[length-2].BAR > 0 && macd[length-2].BAR > macd[length-1].BAR {
		logger.Infoln("<---KDJ down cross,准备卖出, macd:", macd[length-1].BAR)
		trader.SellGroup()
	}

	trader.CancelGroup(int32(timeout))
	if trend < 8 {
		if trader.StopLoss(records[length-1].Close, float64(stoploss)) {
			sty.PrevStoplossTime = time.Now()
		}
	}
	return true
}
