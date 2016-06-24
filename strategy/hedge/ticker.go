package hedge

import(
    . "github.com/qshuai162/qian/common/model"
    "math"
    // "github.com/qshuai162/qian/common/logger"
)

func (hd *Hedge) calDifAvg(t int64) {
    // logger.Infoln(t,hd.PreCalcAvgTime,hd.Option.Refreshtime,len(hd.Tickers[hd.Centers[0].Center]))
    if t-hd.PreCalcAvgTime > int64(hd.Option.Refreshtime) {
        datas:=hd.Tickers
        for i := 0; i < len(hd.Couples); i++ {
            couple := &hd.Couples[i]
            if couple.Center1.Center != "" && couple.Center2.Center != "" {
                avg := hd.cal(datas[couple.Center1.Center], datas[couple.Center2.Center])
                couple.PreAvg = avg
                hd.log("%s %s Offset AVG: %.2f\n", couple.Center1.Center,couple.Center2.Center,avg)
            }
        }
        hd.PreCalcAvgTime = t
        length:=len(hd.Tickers[hd.Centers[0].Center])
        start:=int(math.Ceil(float64(length)*0.3))
        for i:=0;i<len(hd.Centers);i++{
            hd.Tickers[hd.Centers[i].Center]=datas[hd.Centers[i].Center][start:length-1]
        }
        hd.FirstAvg=true
    }
}
func (hd *Hedge) cal(center1,center2 []Ticker) float64{
    length:=len(center1)
    if len(center2)<length{
        length=len(center2)
    }
	sum := 0.0
	count := 0
	for i := 0; i < length; i++ {
		if center1[i].Last != 0 && center2[i].Last != 0 {
			offset := center1[i].Last - center2[i].Last
			sum += offset
			count++
		}
	}
	avg := sum / (float64(count))
	return avg
}

func (hd *Hedge) HandleTicker(t int64) {
	tickt := int64(0)
	if hd.Option.Simulation {
		tickt = t
        for i := 0; i < len(hd.Centers); i++ {
    		 hd.getTicker(&hd.Centers[i], tickt, nil)
	    }
	}else{
	    cc := make(chan int)
        for i := 0; i < len(hd.Centers); i++ {
            go hd.getTicker(&hd.Centers[i], tickt, cc)
        }
        count := 0
        for {
            count += <-cc
            if count >= len(hd.Centers) {
                break
            }
        }
    }
    hd.calDifAvg(t)   
}

func (hd *Hedge) getTicker(center *CData, t int64, cc chan int) {
	if center.Center != "" {
		_, (*center).CurTicker = (*hd.Provider).GetTicker((*center).Center, t)
        hd.Tickers[(*center).Center]=append(hd.Tickers[(*center).Center],(*center).CurTicker)
		hd.log("%s  buy:%.2f\tlast:%.2f\tsell:%.2f\t变化:%.2f\n", (*center).Center, (*center).CurTicker.Buy, (*center).CurTicker.Last, (*center).CurTicker.Sell, (*center).CurTicker.Last-(*center).PreLast)
	}
    if cc!=nil{
	    cc <- 1
    }
}

//检查tick数据
func checkData(oticker, hticker *Ticker) bool {
	if oticker.Last < 1 || hticker.Last < 1 {
		return false
	}
	if oticker.Buy-oticker.Last > 0.3 || hticker.Buy-hticker.Last > 0.3 {
		return false
	}
	if oticker.Last-oticker.Sell > 0.3 || hticker.Last-hticker.Sell > 0.3 {
		return false
	}
	return true
}

