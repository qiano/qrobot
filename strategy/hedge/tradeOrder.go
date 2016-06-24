package hedge

import (
	"github.com/qshuai162/qian/common/logger"
	"github.com/qshuai162/qian/common/model"
	"strconv"
	"time"
	// "math"
	)

func (hd *Hedge) Trade(flag string,center *CData, method string, price, amount float64) int  {
	valid := false
	hd.log("%s btc:%.4f  cny:%.4f  \n",center.Center,center.Btc,center.Cny)
	if method == "buy" && center.Cny >= amount*price {
		center.Cny -= amount * price
		// center.Btc+=amount
		valid = true
	}
	if method == "sell" && center.Btc >= amount {
		center.Btc -= amount
		// center.Cny+=amount*price
		valid = true
	}

	if valid {
		var code, id int
		if method == "buy" {
			code, id = (*hd.Trader).Buy(center.Center, price, amount)
		} else {
			code, id = (*hd.Trader).Sell(center.Center, price, amount)
		}
		if code == 0 && id != 0 {
			if method == "buy" {
				logger.TradeSingle(flag, ",%s,%d,%d,%.2f,%.2f,%.2f", center.Center, -1, id, amount, price, -1*amount*price)
			} else {
				logger.TradeSingle(flag, ",%s,%d,%d,%.2f,%.2f,%.2f", center.Center, 1, id, amount, price, amount*price)
			}
			return id
		}
	}
	return 0
}

func (hd *Hedge) CancelOrder(flag string,oid int,center *CData,method string,price,amount float64) bool{
	_, ret := (*hd.Trader).Cancel(center.Center, oid)
	if ret {
		if method=="buy"{
			center.Cny+=amount*price
			logger.TradeSingle(flag, ",%s,%d,%d,%.2f,%.2f,%.2f,cancel", center.Center, 1, oid,amount,price, amount*price)
		}else{
			center.Btc+=amount
			logger.TradeSingle(flag, ",%s,%d,%d,%.2f,%.2f,%.2f,cancel", center.Center, -1, oid,amount,price, -1*amount*price)
		}
		return true
	} 
	logger.TradeSingle(flag, ",%s,,%d,%.2f,%.2f,,cancel faild", center.Center, oid, amount,price)
	return false
}
//GetOrder 获取委托信息
func (hd *Hedge) GetOrder(center *CData,oid int) model.Order{
	_,or:=(*hd.Trader).GetOrder(center.Center,oid)
	return or
}


//同步下单
func (hd *Hedge) syncTrade(flag int, first string, buyCenter, sellCenter *CData, buyPrice, sellPrice, amount float64) {
	if buyCenter.Center == first {
		bid := hd.Trade(strconv.Itoa(flag), buyCenter, "buy", buyPrice, amount)
		if bid != 0 {
			sid := hd.Trade(strconv.Itoa(flag), sellCenter, "sell", sellPrice, amount)
			if sid == 0 {
				if !hd.CancelOrder(strconv.Itoa(flag),bid,buyCenter,"buy",buyPrice,amount){
					hd.FailedOrders <- FailedOrder{flag, sellCenter, "sell", sellPrice, amount}
				}
			}
		}
	}
	if sellCenter.Center == first {
		sid := hd.Trade(strconv.Itoa(flag), sellCenter, "sell", sellPrice, amount)
		if sid != 0{
			bid := hd.Trade(strconv.Itoa(flag), buyCenter, "buy", buyPrice, amount)
			if bid == 0 {
				if !hd.CancelOrder(strconv.Itoa(flag),sid,sellCenter,"sell",sellPrice,amount){
					hd.FailedOrders <- FailedOrder{flag, buyCenter, "buy", buyPrice, amount}
				}
			}
		}
	}
}

//同步安全下单
func (hd *Hedge) syncSafeTrade(flag int, first string, buyCenter, sellCenter *CData, buyPrice, sellPrice, amount float64,donetime int) {
	if buyCenter.Center == first {
		bid := hd.Trade(strconv.Itoa(flag), buyCenter, "buy", buyPrice, amount)
		if bid != 0 {
			time.Sleep(time.Duration(donetime)*time.Millisecond)
			or:=hd.GetOrder(buyCenter,bid)
			sellamount:=or.Amount
			if or.Status!=2 {
				logamount:=or.Amount-or.Deal_amount
				if or.Deal_amount<0.01{
					logamount=or.Amount
				}
				result:=hd.CancelOrder(strconv.Itoa(flag),bid,buyCenter,"buy",or.Price,logamount)
				if result {
					sellamount=or.Deal_amount
				}
			}
			if sellamount>=0.01{
				sid := hd.Trade(strconv.Itoa(flag), sellCenter, "sell", sellPrice, sellamount)
				if sid == 0 {
					hd.FailedOrders <- FailedOrder{flag, sellCenter, "sell", sellPrice, sellamount}
				}	
			}
		}
	}
	if sellCenter.Center == first {
		sid := hd.Trade(strconv.Itoa(flag), sellCenter, "sell", sellPrice, amount)
		if sid != 0{
			time.Sleep(time.Duration(donetime)*time.Millisecond)
			or:=hd.GetOrder(sellCenter,sid)
			buyamount:=or.Amount
			if or.Status!=2 {
				logamount:=or.Amount-or.Deal_amount
				if or.Deal_amount<0.01{
					logamount=or.Amount
				}
				result:=hd.CancelOrder(strconv.Itoa(flag),sid,sellCenter,"sell",or.Price,logamount)
				if result {
					buyamount=or.Deal_amount
				}
			}
			if buyamount>=0.01{
				bid := hd.Trade(strconv.Itoa(flag), buyCenter, "buy", buyPrice, buyamount)
				if bid == 0 {
					hd.FailedOrders <- FailedOrder{flag, buyCenter, "buy", buyPrice, buyamount}
				}	
			}
		}
	}
}
//异步下单
func (hd *Hedge) asynTrade(flag int, buyCenter, sellCenter *CData, buyPrice, sellPrice, amount float64) {
	go func() {
		id := hd.Trade(strconv.Itoa(flag), buyCenter, "buy", buyPrice, amount)
		if id == 0 {
			hd.FailedOrders <- FailedOrder{flag, buyCenter, "buy", buyPrice, amount}
		}
	}()
	go func() {
		id := hd.Trade(strconv.Itoa(flag), sellCenter, "sell", sellPrice, amount)
		if id == 0 {
			hd.FailedOrders <- FailedOrder{flag, sellCenter, "sell", sellPrice, amount}
		}
	}()
}

//FailedOrder 失败的委托
type FailedOrder struct {
	Flag   int
	Center *CData
	Type   string
	Price  float64
	Amount float64
}

func handlefailedOrders(hd *Hedge) {
	for {
		bo := <-hd.FailedOrders
		time.Sleep(3 * time.Second)
		logger.Infof("开始重试失败的委托: %s %s %.2f %.2f\n", bo.Center.Center, bo.Type, bo.Price, bo.Amount)
		price := bo.Price
		amount := bo.Amount
		if (bo.Type == "buy" && bo.Center.Cny >= price*amount) ||
			(bo.Type == "sell" && bo.Center.Btc >= amount) {
			go func(){
				curPrice:=price
				if curPrice!=0 &&
				((bo.Type=="buy" && bo.Center.CurTicker.Last<curPrice-0.1) ||
				 (bo.Type=="sell" && bo.Center.CurTicker.Last>curPrice+0.1)) {
					curPrice=bo.Center.CurTicker.Last
				}
				id := hd.Trade(strconv.Itoa(bo.Flag), bo.Center, bo.Type, curPrice, amount)
				if id == 0 {
					hd.FailedOrders <- FailedOrder{bo.Flag, bo.Center, bo.Type, price, amount}
				}
			}()
		} else {
			hd.log("%s 资产不足！！！\n", bo.Center.Center)
			hd.FailedOrders <- FailedOrder{bo.Flag, bo.Center, bo.Type, price, amount}
		}
	}
}

