package kdj

import (
	"github.com/qshuai162/qian/common/logger"
	idc "github.com/qshuai162/qian/qrobot/indictor"
	. "github.com/qshuai162/qian/qrobot/trader"
	"github.com/qshuai162/qian/qrobot/strategy/hedge"
	// "strconv"	
	. "github.com/qshuai162/qian/common/model"

	// "fmt"
	// "time"
)

type KdjStrategy struct {
	Option           KdjOp
	PrevBuyTime      map[string]int64
	PrevStoplossTime int64
	Orders           map[string]*OCaches
	HD   *hedge.Hedge
	
}

func (kdj *KdjStrategy) getOrders(str string) *OCaches{
	temp:=kdj.Orders[str]
	return temp
}

func NewKdjStrategy(t int64, Option KdjOp,hd *hedge.Hedge) *KdjStrategy {
	kdj := new(KdjStrategy)
	kdj.Option = Option
	kdj.PrevStoplossTime = t - 3600*24*365
	kdj.HD=hd

	kdj.Orders=make(map[string]*OCaches)
	kdj.Orders["KDJ1"]=new(OCaches)
	kdj.Orders["MACD1"]=new(OCaches)
	kdj.getOrders("KDJ1").SaveHistory=true
	kdj.getOrders("MACD1").SaveHistory=true
	
	kdj.PrevBuyTime=make(map[string]int64)
	kdj.PrevBuyTime["KDJ1"]=t-3600*24*365
	kdj.PrevBuyTime["MACD1"]=t-3600*24*365
	return kdj
}

func (kdj *KdjStrategy) Tick(t int64,records []Record,ticker Ticker) bool {
	length:=len(records)
	if kdj.Option.Simulation {
		//kdj.Orders.SimDone(kdj.Option.Center,ticker.Last)
	} else {
		_, orders := (*kdj.HD.Provider).GetOrders(kdj.Option.Center)
		kdj.getOrders("KDJ1").Sync(orders)
		kdj.getOrders("MACD1").Sync(orders)
	}
	
	go kdj.kdjCore(t,records,length,ticker)
	go kdj.macdCore(t,records,length,ticker)
	return true 
}

func (kdj *KdjStrategy) kdjCore(t int64,records []Record,length int,ticker Ticker){
	// K线为白，D线为黄，J线为红，K in middle
	k, d, j := idc.GetKDJ(records, 9)
	kdj.log("LAST: d(黄线）%0.2f\tk(白线）%0.2f\tj(红线）%0.2f\n", d[length-2], k[length-2], j[length-2])
	kdj.log("CURR: d(黄线）%0.2f\tk(白线）%0.2f\tj(红线）%0.2f\n", d[length-1], k[length-1], j[length-1])
	//买
	if j[length-4]>j[length-3] && j[length-3]<j[length-2] && j[length-2]<j[length-1] && j[length-1]<kdj.Option.JHigh {
		kdj.log("--->KDJ buy\n")
		kdj.Buy("KDJ1",t,ticker)
	}
	//卖
	if j[length-1]<j[length-2]{
		kdj.log("<---KDJ sell\n")
		kdj.Sell("KDJ1",t,ticker)
	}
	kdj.TryCancel("KDJ1",t,ticker) 
	kdj.TryStoploss("KDJ1",t,ticker)
}

func (kdj *KdjStrategy) macdCore(t int64,records []Record,length int,ticker Ticker){
	macd := idc.GetMACD(records, 12, 26, 9)
	kdj.log("Prev:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-3].Close, macd[length-3].DIF, macd[length-3].DEA, macd[length-3].BAR)
	kdj.log("LAST:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-2].Close, macd[length-2].DIF, macd[length-2].DEA, macd[length-2].BAR)
	kdj.log("CURR:price=%.2f\tdif=%0.2f\tdea=%0.2f\tmacd=%0.2f\n", records[length-1].Close, macd[length-1].DIF, macd[length-1].DEA, macd[length-1].BAR)
	//买
	if macd[length-4].BAR>macd[length-3].BAR   && macd[length-2].BAR > macd[length-3].BAR &&  macd[length-1].BAR > macd[length-2].BAR  {
		kdj.log("--->MACD buy\n")
		kdj.Buy("MACD1",t,ticker)
	}
	//卖
	if macd[length-1].BAR<macd[length-2].BAR{ 
		kdj.log("<---MACD sell\n")
	    kdj.Sell("MACD1",t,ticker)
	}
	kdj.TryCancel("MACD1",t,ticker) 
	kdj.TryStoploss("MACD1",t,ticker)
}




func (hd *KdjStrategy) log(formate string, args ...interface{}) {
	if hd.Option.LogEnabled {
		logger.Infof(formate, args...)
	}
}
