package trader

import (
	"fmt"
	"github.com/qshuai162/qian/api/chbtc"
	"github.com/qshuai162/qian/api/huobi"
	"github.com/qshuai162/qian/api/okcoin"
	"github.com/qshuai162/qian/api/btcc"
	. "github.com/qshuai162/qian/common/model"
	// . "github.com/qshuai162/qian/tool/config"
	"github.com/qshuai162/qian/common/logger"
	"strconv"
	// "time"
)
type TradeAPI interface {
	Buy(price, amount string) string
	Sell(price, amount string) string
	GetOrder(order_id string) (ret bool, order Order)
	CancelOrder(order_id string) bool
	GetAccount() (Account, bool)
	GetOrderBook() (ret bool, orderBook OrderBook)
}

type Trader struct {
}

func getTradeAPI(name string) (tradeAPI TradeAPI) {
	if name == "huobi" {
		tradeAPI = huobi.NewHuobi()
	} else if name == "okcoin" {
		tradeAPI = okcoin.NewOkcoin()
	} else if name == "chbtc" {
		tradeAPI = chbtc.NewCHBTC()
	} else if name == "btcc" {
		tradeAPI = btcc.NewBTCC()
	} else {
		logger.Fatalln("Please config the exchange center...")
		panic(0)
	}
	//  else if name == "bitvc" {
	// 	tradeAPI = bitvc.NewBitvc()
	// } else if name == "peatio" {
	// 	tradeAPI = peatio.NewPeatio()
	// } else if name == "bittrex" {
	// 	tradeAPI = Bittrex.Manager()
	// } else if name == "simulate" {
	// 	tradeAPI = simulate.NewSimulate()
	// }

	return
}

func (tr *Trader) Buy(center string, price, amount float64) (int, int) {
	api := getTradeAPI(center)
	buyID := api.Buy(fmt.Sprintf("%.2f", price), fmt.Sprintf("%.2f", amount))
	bid, _ := strconv.Atoi(buyID)
	if bid != 0 {
		logger.Infof("[委托成功] %s 买入 <--限价单 数量:%.2f  价格:%.2f  单号：%d", center, amount, price, bid)
		return 0, bid
	} else {
		logger.Infoln("[委托失败]")
		return 1, bid
	}
}

func (tr *Trader) Sell(center string, price, amount float64) (int, int) {
	api := getTradeAPI(center)
	sellID := api.Sell(fmt.Sprintf("%.2f", price), fmt.Sprintf("%.2f", amount))
	sid, _ := strconv.Atoi(sellID)
	if sid != 0 {
		logger.Infof("[委托成功] %s 卖出 <--限价单 数量:%.2f  价格:%.2f  单号：%d", center, amount, price, sid)
		return 0, sid
	} else {
		logger.Infoln("[委托失败]")
		return 1, sid
	}
}

func (tr *Trader) Cancel(center string, id int) (code int, ret bool) {
	api := getTradeAPI(center)

	if api.CancelOrder(strconv.Itoa(id)) {
		logger.Infof("[Cancel委托成功] %s <-----撤单------>  %d\n", center, id)
		ret = true
	} else {
		logger.Infoln("[Cancel委托失败]")
		ret = false
	}
	return
}

func (tr *Trader) GetOrder(center string,orderId int) (bool , Order){
	api := getTradeAPI(center)
	ret,order:=api.GetOrder(strconv.Itoa(orderId))
	return ret,order
}

// //根据实时买卖盘添加一组买单
// func (tr *Trader) BuyGroup(amount float64) bool {
// 	ordercount, _ := strconv.Atoi(Option["splitordercount"])
// 	stepprice, _ := strconv.ParseFloat(Option["stepprice"], 64)
// 	nSplitTradeAmount := amount / float64(ordercount)
// 	splitTradeAmount := fmt.Sprintf("%.2f", nSplitTradeAmount)

// 	buyPrice := tr.GetBuyPrice()
// 	if buyPrice != 0 {
// 		for i := 1; i <= ordercount; i++ {
// 			tradePrice := fmt.Sprintf("%.2f", buyPrice+stepprice*float64(i))
// 			tr.Buy(tradePrice, splitTradeAmount)
// 		}
// 		return true
// 	}
// 	return false
// }

// //撤消所有超时订单
// func (tr *Trader) CancelGroup(timeout int32) {
// 	forcancel := tr.OCaches.WaitForCancel(timeout)
// 	for i := 0; i < len(forcancel); i++ {
// 		cancelorder := tr.OCaches.Find(forcancel[i])
// 		tr.Cancel(cancelorder.Id)
// 	}

// }

// //卖出所有可卖单
// func (tr *Trader) SellGroup() bool {
// 	minprofit, _ := strconv.ParseFloat(Option["minprofit"], 64)
// 	stepprice, _ := strconv.ParseFloat(Option["stepprice"], 64)
// 	sellprice := tr.GetSellPrice()
// 	if sellprice != 0 {
// 		forsell := tr.OCaches.WaitForSell(sellprice - minprofit)
// 		for i := 0; i < len(forsell); i++ {
// 			buyorder := tr.OCaches.Find(forsell[i])
// 			tradePrice := fmt.Sprintf("%.2f", sellprice-stepprice*float64(i))
// 			amount := fmt.Sprintf("%.2f", buyorder.Order_Amount)
// 			tr.Sell(buyorder.Id, tradePrice, amount)
// 		}
// 		return true
// 	}
// 	return false
// }

// //止损操作
// func (tr *Trader) StopLoss(currentPrice, stopLossRate float64) bool {
// 	forstop := tr.OCaches.WaitForStopLoss(currentPrice, stopLossRate)
// 	if len(forstop) != 0 {
// 		logger.Infoln("stop loss")
// 		price := fmt.Sprintf("%.2f", currentPrice)
// 		for i := 0; i < len(forstop); i++ {
// 			stoporder := tr.OCaches.Find(forstop[i])
// 			amount := fmt.Sprintf("%.2f", stoporder.Order_Amount)
// 			if stoporder.RefId == 0 {
// 				tr.Sell(stoporder.Id, price, amount)
// 			}
// 		}
// 		return true
// 	}
// 	return false
// }

// func (tr *Trader) GetOrderBook() *OrderBook {
// 	ret, orderbook := tr.API.GetOrderBook()
// 	if !ret {
// 		ret, orderbook = tr.API.GetOrderBook() // try again
// 		if !ret {
// 			logger.Infoln("get orderbook failed 2")
// 			return nil
// 		}
// 	}
// 	logger.Infoln("卖一", (orderbook.Asks[len(orderbook.Asks)-1]))
// 	logger.Infoln("买一", orderbook.Bids[0])
// 	return &orderbook
// }

// //获取买入价
// func (tr *Trader) GetBuyPrice() float64 {
// 	orderbook := tr.GetOrderBook()
// 	if orderbook != nil {
// 		return orderbook.Bids[0].Price
// 	} else {
// 		return 0
// 	}
// }

// //获取卖出价
// func (tr *Trader) GetSellPrice() float64 {
// 	orderbook := tr.GetOrderBook()
// 	if orderbook != nil {
// 		return orderbook.Asks[len(orderbook.Asks)-1].Price
// 	} else {
// 		return 0
// 	}
// }

// func (tr *Trader) GetAccount() (Account, bool) {
// 	return tr.API.GetAccount()
// }
