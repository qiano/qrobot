package strategy

import (
	. "github.com/qshuai162/qian/common/config"
	. "github.com/qshuai162/qian/common/model"
	// "github.com/qshuai162/btcrobot/src/huobi"
	"fmt"
	"github.com/qshuai162/qian/common/logger"
	idc "github.com/qshuai162/qian/qrobot/indictor"
	. "github.com/qshuai162/qian/qrobot/trader"
	"strconv"
	"time"
)

type MACDStrategyQK1 struct {
	PrevBuyTime time.Time
	Orders      OCaches
}

func init() {
	macd := new(MACDStrategyQK1)
	macd.PrevBuyTime = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	Register("MACD-K1-Q", macd)
}

// KDJ-EX strategy
func (sty *MACDStrategyQK1) Tick() bool {
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
	// timeout, _ := strconv.Atoi(Option["timeout"])
	// stoploss, _ := strconv.ParseFloat(Option["stoploss"], 64)
	// stoploss_resttime, _ := strconv.Atoi(Option["stoploss_resttime"])
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])
	signalPeriod, _ := strconv.Atoi(Option["signalPeriod"])
	// MACDbuyThreshold, _ := strconv.ParseFloat(Option["MACDbuyThreshold"], 64)
	// MACDsellThreshold, _ := strconv.ParseFloat(Option["MACDsellThreshold"], 64)

	logger.Infoln("策略: MACD-K1-Q  MACD参数:", shortEMA, longEMA, signalPeriod, "K线周期:", kperoid)
	sty.Orders.Sync(orders)

	macd := idc.GetMACD(records, shortEMA, longEMA, signalPeriod)
	logger.Infoln("当前价", records[length-1].Close)
	logger.Infof("Prev:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-3].Close, macd[length-3].DIF, macd[length-3].DEA, macd[length-3].BAR)
	logger.Infof("LAST:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-2].Close, macd[length-2].DIF, macd[length-2].DEA, macd[length-2].BAR)
	logger.Infof("CURR:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-1].Close, macd[length-1].DIF, macd[length-1].DEA, macd[length-1].BAR)

	//卖单超时处理，撤销后再次按实价挂卖单
	buys, sells := sty.Orders.Buy_Selling()
	if len(sells) != 0 {
		sellprice := trader.GetSellPrice()
		for i := 0; i < len(sells); i++ {
			s := sty.Orders.Find(sells[i])
			if time.Now().Sub(s.Time) >= 300*time.Second {
				logger.Infoln("resell")
				if trader.Cancel(s.Id) {
					buyorder := sty.Orders.Find(s.RefId)
					tradePrice := fmt.Sprintf("%.2f", sellprice)
					amount := fmt.Sprintf("%.2f", buyorder.Order_Amount)
					trader.Sell(buyorder.Id, tradePrice, amount)
				}
			}
		}
	}

	//新卖单
	if orders != nil {
		sty.addSell(&market, &trader)
	}

	//买
	amount := numTradeAmount * trend
	if amount > 0 {
		c, r := DownRate(records)
		if macd[length-1].BAR > macd[length-2].BAR || c >= 3 {
			logger.Infoln("--->KDJ up cross,准备买入, macd:", macd[length-1].BAR)
			if len(buys) < 20 && time.Now().Sub(sty.PrevBuyTime) >= 60*time.Second {
				if (r < -0.0008 || c >= 3) && records[length-1].Close > records[length-2].Close {

					if trader.BuyGroup(amount) {
						sty.PrevBuyTime = time.Now()
						time.Sleep(1 * time.Second)
						sty.Orders.Sync(market.GetOrders(true))
						sty.addSell(&market, &trader)
					}
				}
			}
		}
	}

	trader.CancelGroup(300)
	return true
}

//下跌百分比
func DownRate(records []Record) (count int, rate float64) {
	count = 0
	length := len(records)
	logger.Infof("records[%d]\topen=%.2f\tclose=%.2f\t\n", length-2, records[length-2].Open, records[length-2].Close)
	for i := length - 2; i >= 0; i-- {
		if records[i-1].Open < records[i].Close {
			logger.Infof("records[%d]\topen=%.2f\tclose=%.2f\t\n", i-1, records[i-1].Open, records[i-1].Close)
			open := records[i].Open
			rate = (records[length-1].Close - open) / open
			logger.Infoln("跌幅：", fmt.Sprintf("%.5f", rate))
			return
		}
		count++
		logger.Infof("records[%d]\topen=%.2f\tclose=%.2f\t\n", i-1, records[i-1].Open, records[i-1].Close)
	}
	return
}

func (sty *MACDStrategyQK1) addSell(market *Market, trader *Trader) {
	forsell := sty.Orders.Buy_NotSell()
	if len(forsell) > 0 {
		sellprice := trader.GetSellPrice() - 0.1
		for i := 0; i < len(forsell)-1; i++ {
			buyorder := sty.Orders.Find(forsell[i])
			var sprice float64
			if buyorder.Price+0.6 < sellprice {
				sprice = sellprice
			} else {
				sprice = buyorder.Price + 1.5
			}
			tradePrice := fmt.Sprintf("%.2f", sprice)
			amount := fmt.Sprintf("%.2f", buyorder.Order_Amount)
			trader.Sell(buyorder.Id, tradePrice, amount)
		}
	}
}
