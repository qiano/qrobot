package kdj

import(
    . "github.com/qshuai162/qian/common/model"
    "github.com/qshuai162/qian/common/logger"
    // "fmt"
)

func (kdj *KdjStrategy) Buy(str string,t int64,ticker Ticker) bool{
    if t-kdj.PrevStoplossTime >= 60*5 {
        if len(kdj.getOrders(str).Buy_NotSell()) < kdj.Option.MaxCount && t-kdj.PrevBuyTime[str] >= 15 {
            amount := kdj.Option.PerAmount 
            price:=ticker.Sell+0.03
            bid:=kdj.HD.Trade(str,&kdj.HD.Centers[0],"buy",price,amount)
            if bid!=0{
                kdj.getOrders(str).Buy(t,kdj.Option.Center,bid,price,amount)
                // kdj.getOrders(str).PrintAll()
                kdj.PrevBuyTime[str] = t
                return true
            }
        }
    } else {
        kdj.log("止损保护中，不允许交易\n")
    }
    return false
}

func (kdj *KdjStrategy) Sell(str string,t int64,ticker Ticker) bool {
    sellPrice := ticker.Buy - 0.03
    forsell := kdj.getOrders(str).WaitForSell(sellPrice - kdj.Option.MinProfit)
    for i := 0; i < len(forsell); i++ {
        buyorder := kdj.getOrders(str).Find(forsell[i])
        sid:=kdj.HD.Trade(str,&kdj.HD.Centers[0],"sell",sellPrice,buyorder.Order_Amount)
        if sid!=0{
            kdj.getOrders(str).Sell(t,kdj.Option.Center,buyorder.Id,sid,sellPrice,buyorder.Order_Amount)
        }
    }
    return true
}

func (kdj *KdjStrategy) TryCancel(str string,t int64,ticker Ticker)  {
    forcancel := kdj.getOrders(str).WaitForCancel(t - int64(kdj.Option.TimeOut))
	for i := 0; i < len(forcancel); i++ {
		cancelorder := kdj.getOrders(str).Find(forcancel[i])
		ret:=kdj.HD.CancelOrder(str,cancelorder.Id,&kdj.HD.Centers[0],cancelorder.Type,cancelorder.Price,cancelorder.Order_Amount)
		if ret{
			kdj.getOrders(str).Cancel(cancelorder.Id)
		}
	}
}

func (kdj *KdjStrategy) TryStoploss(str string,t int64,ticker Ticker){
    p := ticker.Buy - 0.1
	forstop := kdj.getOrders(str).WaitForStopLoss(p, kdj.Option.StopLoss)
	if len(forstop) != 0 {
		kdj.log("stop loss")
        logger.TradeSingle(str, ",%s,,,,,,%s", kdj.HD.Centers[0].Center,"stoploss")
		for i := 0; i < len(forstop); i++ {
			stoporder := kdj.getOrders(str).Find(forstop[i])
			if stoporder.RefId == 0 {
				sid:=kdj.HD.Trade("KDJ1",&kdj.HD.Centers[0],"sell",p,stoporder.Order_Amount)
				if sid!=0{
					kdj.getOrders(str).Sell(t,kdj.Option.Center,stoporder.Id,sid,p,stoporder.Order_Amount)
				}
			}
		}
		kdj.PrevStoplossTime = t
	}
}