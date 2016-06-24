package indictor

import (
	. "github.com/qshuai162/qian/common/model"
	// "math"
)

type MACD struct {
	EMAShort float64
	EMALong  float64
	DIF      float64
	DEA      float64
	BAR      float64
}

func GetMACD(records []Record, periodShort int, periodLong int, periodDIF int) (macd []MACD) {

	if len(records) == 0 {
		return nil
	}
	var price []float64
	for _, v := range records {
		price = append(price, v.Close)
	}

	shortEma := EMA(price, periodShort)
	longEma := EMA(price, periodLong)
	macd = make([]MACD, len(price))

	macd[0].EMAShort = price[0]
	macd[0].EMALong = price[0]
	macd[0].DIF = 0
	macd[0].DEA = 0
	macd[0].BAR = 0

	for i := 1; i < len(price); i++ {
		macd[i].EMAShort = shortEma[i]
		macd[i].EMALong = longEma[i]
		macd[i].DIF = macd[i].EMAShort - macd[i].EMALong
		//DEA（MACD）= 前一日DEA×8/10＋今日DIF×2/10
		//DIFF平均值9日（DEA）=9日DIFF之和/9
		if i >= periodDIF {
			sum := 0.0
			for j := 0; j < periodDIF; j++ {
				sum += macd[i-j].DIF
			}
			avg := sum / float64(periodDIF)
			//value := math.Trunc(avg*1e4+0.5) * 1e-4
			macd[i].DEA = avg
		} else {
			macd[i].DEA = 0
		}
		macd[i].BAR = 2 * (macd[i].DIF - macd[i].DEA)
	}
	return
}
