package hedge

import (
	"github.com/qshuai162/qian/common/logger"
	. "github.com/qshuai162/qian/common/model"
	"math"
	// . "github.com/qshuai162/qian/qrobot/trader"
)

//CData 交易中心数据
type CData struct {
	Center    string
	PreLast   float64
	CurTicker Ticker
	Btc float64
	Cny float64
}

func newCData(center string) *CData {
	c := new(CData)
	c.Center = center
	c.CurTicker = *new(Ticker)
	return c
}

type hedgeCouple struct {
	Center1        *CData
	Center2        *CData
	PreAvg         float64
	PrevBuyTime    int64
	ThawingTime    int64
	MinProfit      float64 //单边最小利润
	PerAmount      float64 //单次交易数量
}

func newCouple(center1 *CData, center2 *CData, minprofit, peramount float64) *hedgeCouple {
	c := new(hedgeCouple)
	c.Center1 = center1
	c.Center2 = center2
	c.MinProfit = minprofit
	c.PerAmount = peramount
	return c
}

//Hedge 对冲结构
type Hedge struct {
	Option       HedgeOp
	Centers      []CData
	Couples      []hedgeCouple
	Trader       *ITrader
	Provider     *IProvider
	FailedOrders chan FailedOrder
	CurTime 	 int64
	Tickers      map[string][]Ticker
	PreCalcAvgTime int64
	FirstAvg bool          
}

//NewHedge 新建
func NewHedge(t int64, Option HedgeOp, Provider IProvider, Trader ITrader) *Hedge {
	hd := new(Hedge)
	hd.Option = Option
	hd.Trader = &Trader
	hd.Provider = &Provider
	hd.FailedOrders = make(chan FailedOrder, 999)
	hd.Tickers=make(map[string][]Ticker)
	hd.CurTime=t
	hd.PreCalcAvgTime=t
	hd.FirstAvg=false
	for i := 0; i < len(Option.Centers); i++ {
		hd.Centers = append(hd.Centers, *newCData(Option.Centers[i]))
	}
	for i := 0; i < len(Option.HdConfig); i++ {
		cf := hd.Option.HdConfig[i]
		hd.Couples = append(hd.Couples, *newCouple(&hd.Centers[int(cf[0])], &hd.Centers[int(cf[1])], cf[2], cf[3]))
	}

	hd.calDifAvg(t)
	go handlefailedOrders(hd)
	if !hd.Option.Simulation{
		go hd.startSyncAccount()
	}
	return hd
}


func (hd *Hedge) log(formate string, args ...interface{}) {
	if hd.Option.LogEnabled {
		logger.Infof(formate, args...)
	}
}

//Tick 周期响应
func (hd *Hedge) Tick(t int64) bool {
	hd.CurTime=t
	if hd.FirstAvg{
		cc := make(chan int)
		for i := 0; i < len(hd.Couples); i++ {
			go hd.hedgeCore(i, &hd.Couples[i], t, cc)
		}
		count := 0
		for {
			count += <-cc
			if count >= len(hd.Option.HdConfig) {
				break
			}
		}
	}
	return true
}

func (hd *Hedge) hedgeCore(flag int, couple *hedgeCouple, t int64, cc chan int) {
	center1 := couple.Center1
	center2 := couple.Center2
	if center1.Center == "" || center2.Center == "" {
		cc <- 1
		return
	}
	oticker := center1.CurTicker
	hticker := center2.CurTicker

	check := checkData(&oticker, &hticker)
	if !check {
		hd.log("%s  %s Ticker 数据异常\n", center1.Center, center2.Center)
		cc <- 1
		return
	}

	offset := oticker.Last - hticker.Last
	hd.log("%s - %s = %.2f\t上值:%.2f\t下值:%.2f\n", center1.Center, center2.Center, offset, couple.PreAvg+couple.MinProfit, couple.PreAvg-couple.MinProfit)
	
	if center1.PreLast>1 && center2.PreLast>1{
		if (math.Abs(oticker.Last-center1.PreLast) >= hd.Option.TooFast ||math.Abs(hticker.Last-center2.PreLast) >= hd.Option.TooFast) {
			hd.log("%s  %s价格变化过快，暂停%d秒\n", center1.Center, center2.Center, hd.Option.ColdDownTime)
			couple.ThawingTime = t+int64(hd.Option.ColdDownTime)
		}
		if (math.Abs(oticker.Last-center1.PreLast) >=hd.Option.Fast ||math.Abs(hticker.Last-center2.PreLast) >= hd.Option.Fast) {
			if couple.ThawingTime-t<int64(hd.Option.ColdTime){
				hd.log("%s  %s价格变化略快，暂停%d秒\n", center1.Center, center2.Center, hd.Option.ColdTime)
				couple.ThawingTime = t+int64(hd.Option.ColdTime)
			}
		}
	}
	center1.PreLast = oticker.Last
	center2.PreLast = hticker.Last
	if couple.ThawingTime-t>0 {
		hd.log("%s  %s市场冷却中,等待%d秒\n", center1.Center, center2.Center,couple.ThawingTime-t)
		cc <- 1
		return
	}
	
	if t-couple.PrevBuyTime > int64(hd.Option.RestTime) {
		a:=0.03
		// at:=couple.MinProfit/2-a
		if offset > couple.PreAvg+couple.MinProfit {
			buyPrice := hticker.Last + a
			sellPrice := oticker.Last - a
			hd.log("%s -> %s   数量:%.2f\n", center1.Center, center2.Center, couple.PerAmount)
			if center1.Btc >= couple.PerAmount && center2.Cny >= couple.PerAmount*buyPrice {
				// hd.asynTrade(flag, center2, center1, buyPrice, sellPrice, couple.PerAmount)
				hd.syncTrade(flag,center1.Center,center2,center1,buyPrice,sellPrice,couple.PerAmount)
				// hd.syncSafeTrade(flag, center1.Center, center2, center1, buyPrice, sellPrice, couple.PerAmount,couple.DoneTime)
				couple.PrevBuyTime = t
			} else {
				hd.log("%s  %s资产不足！！！\n", center1.Center, center2.Center)
			}
		}

		if offset < couple.PreAvg-couple.MinProfit {
			buyPrice := oticker.Last + a 
			sellPrice := hticker.Last - a
			hd.log("%s -> %s   数量:%.2f\n", center2.Center, center1.Center, couple.PerAmount)
			if center2.Btc >= couple.PerAmount && center1.Cny >= couple.PerAmount*buyPrice {
				// hd.syncSafeTrade(flag, center1.Center, center1, center2, buyPrice, sellPrice, couple.PerAmount,couple.DoneTime)
				hd.syncTrade(flag,center1.Center,center1,center2,buyPrice,sellPrice,couple.PerAmount)
				// hd.asynTrade(flag, center1, center2, buyPrice, sellPrice, couple.PerAmount)
				couple.PrevBuyTime = t
			} else {
				hd.log("%s  %s资产不足！！！\n", center1.Center, center2.Center)
			}
		}
	} else {
		hd.log("%s  %s休息%d秒\n", center1.Center, center2.Center, hd.Option.RestTime)
	}   
	cc <- 1
}