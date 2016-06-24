package quantify

import (
	"github.com/qshuai162/qian/common/logger"
	. "github.com/qshuai162/qian/common/model"
	idc "github.com/qshuai162/qian/qrobot/indictor"
	. "github.com/qshuai162/qian/qrobot/trader"
	"github.com/qshuai162/qian/qrobot/strategy/hedge"
	// "strconv"
	// "fmt"
	"time"
)

type QuantifyStrategy struct {
	Option           QuantifyOp
	PrevBuyTime      int64
	PrevStoplossTime int64
	Orders           OCaches
	HD   *hedge.Hedge
}	


func NewQuantifyStrategy(t int64, Option QuantifyOp,hd *hedge.Hedge) *QuantifyStrategy {
	qu := new(QuantifyStrategy)
	qu.Option = Option
	qu.PrevBuyTime = t - 3600*24*365
	qu.PrevStoplossTime = t - 3600*24*365
	qu.HD=hd
	qu.Orders.SaveHistory=true
	return hd
}

func (kdj *QuantifyStrategy) Tick(t int64,records []Record,ticker Ticker) bool {
	
	if kdj.Option.Simulation {
		kdj.Orders.SimDone(kdj.Option.Center,ticker.Last)
	} else {
		_, orders := (*kdj.Provider).GetOrders(kdj.Option.Center)
		kdj.Orders.Sync(orders)
	}

	length:=len(records)
	macd := idc.GetMACD(records, 12, 26, 9)
	kdj.log("Prev:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-3].Close, macd[length-3].DIF, macd[length-3].DEA, macd[length-3].BAR)
	kdj.log("LAST:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-2].Close, macd[length-2].DIF, macd[length-2].DEA, macd[length-2].BAR)
	kdj.log("CURR:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-1].Close, macd[length-1].DIF, macd[length-1].DEA, macd[length-1].BAR)

	//买
	if macd[length-4].BAR>macd[length-3].BAR   && macd[length-2].BAR > macd[length-3].BAR &&  macd[length-1].BAR > macd[length-2].BAR  {
		kdj.log("--->MACD up cross,准备买入\n")
		if t-kdj.PrevStoplossTime >= 60*5 {
			if len(kdj.Orders.Buy_NotSell()) < kdj.Option.MaxCount && t-kdj.PrevBuyTime >= 15 {
				amount := kdj.Option.PerAmount 
				price:=ticker.Sell+0.03
				if amount > 0 {
					bid:=kdj.HD.Trade("MACD1",&kdj.HD.Centers[0],"buy",price,amount)
					if bid!=0{
						kdj.Orders.Buy(t,kdj.Option.Center,bid,price,amount)
					}
					kdj.PrevBuyTime = t
				}
			}
		} else {
			kdj.log("止损保护中，不允许交易\n")
		}
	}

	//卖
	if macd[length-1].BAR<macd[length-2].BAR{ // && macd[length-3].BAR>macd[length-2].BAR {//&& j[length-2]>j[length-1]{
		kdj.log("<---MACD down cross,尝试卖出\n")
		sellPrice := ticker.Buy - 0.03
		forsell := kdj.Orders.WaitForSell(sellPrice - kdj.Option.MinProfit)
		for i := 0; i < len(forsell); i++ {
			buyorder := kdj.Orders.Find(forsell[i])
			sid:=kdj.HD.Trade("MACD1",&kdj.HD.Centers[0],"sell",sellPrice,buyorder.Order_Amount)
			if sid!=0{
				kdj.Orders.Sell(t,kdj.Option.Center,buyorder.Id,sid,sellPrice,buyorder.Order_Amount)
			}
		}
		return true
	}

	forcancel := kdj.Orders.WaitForCancel(t - int64(kdj.Option.TimeOut))
	for i := 0; i < len(forcancel); i++ {
		cancelorder := kdj.Orders.Find(forcancel[i])
		kdj.cancel(t, cancelorder.Id)
		if cancelorder.Type == "sell" {
			kdj.sell(t, cancelorder.Id, ticker.Buy-0.02)
		}
	}

	p := ticker.Buy - 0.1
	forstop := kdj.Orders.WaitForStopLoss(p, kdj.Option.StopLoss)
	if len(forstop) != 0 {
		kdj.log("stop loss")
		for i := 0; i < len(forstop); i++ {
			stoporder := kdj.Orders.Find(forstop[i])
			if stoporder.RefId != 0 {
				sod := kdj.Orders.Find(stoporder.RefId)
				kdj.cancel(t, sod.Id)
			}
			kdj.sell(t, stoporder.Id, p)
			kdj.StoplossCount++
		}
		kdj.PrevStoplossTime = t
	}
	// kdj.Orders.PrintAll()
	return true
}

//下跌百分比
func DownRate(records []Record) (rate float64) {
	length := len(records)
	for i := length - 2; i >= 0; i-- {
		if records[i-1].Open < records[i].Close {
			open := records[i].Open
			rate = (records[length-1].Close - open) / open
			return
		}
	}
	return
}

func (hd *QuantifyStrategy) buy(t int64, price float64, amount float64) {
	_, id := (*hd.Trader).Buy(hd.Option.Center, price, amount)
	if id != 0 {
		// hd.Orders.Buy(t, id, price, hd.Option.PerAmount)
		logger.Tradef(",%s,%s,%d,%d,%.2f,%.2f,%.2f", time.Unix(t, 0).Format("2006-01-02 15:04:05"), hd.Option.Center, -1, id, amount, price, -1*amount*price)
	}
}

func (hd *QuantifyStrategy) sell(t int64, buyid int, price float64) {
	_, id := (*hd.Trader).Sell(hd.Option.Center, price, hd.Option.PerAmount)
	if id != 0 {
		bo := hd.Orders.Find(buyid)
		// hd.Orders.Sell(t, buyid, id, price, hd.Option.PerAmount)
		logger.Tradef(",%s,%s,%d,%d,%.2f,%.2f,%.2f,%.2f", time.Unix(t, 0).Format("2006-01-02 15:04:05"), hd.Option.Center, 1, id, hd.Option.PerAmount, price, hd.Option.PerAmount*price, hd.Option.PerAmount*(price-bo.Price))
	}
}

func (hd *QuantifyStrategy) cancel(t int64, oid int) {
	_, suc := (*hd.Trader).Cancel(hd.Option.Center, oid)
	if suc {
		hd.Orders.Cancel(oid)
		logger.Tradef(",%s,%s,%d,%d,,,", time.Unix(t, 0).Format("2006-01-02 15:04:05"), hd.Option.Center, 0, oid)
	}
}

func (hd *QuantifyStrategy) log(formate string, args ...interface{}) {
	if hd.Option.LogEnabled {
		logger.Infof(formate, args...)
	}
}
