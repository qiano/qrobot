package hedge

import (
	"github.com/qshuai162/qian/common/logger"
	"time"
)

func statistical(hd *Hedge) {
	for {
		if hd.Centers[0].CurTicker.Last != 0 {
			t := time.Now().Unix()
			syncAccountHandle(hd)
			logger.Profitf(",%s", time.Unix(t, 0).Format("2006-01-02 15:04:05"))
			total := newCData("total")
			realas := 0.0
			calas := 0.0
			for i := 0; i < len(hd.Centers); i++ {
				c := hd.Centers[i]
				realAsset := (c.BTC+c.FBTC)*c.CurTicker.Last + c.CNY + c.FCNY
				calAsset := (c.BTC+c.FBTC)*3000 + c.CNY + c.FCNY

				total.BTC += c.BTC
				total.FBTC += c.FBTC
				total.CNY += c.CNY
				total.FCNY += c.FCNY
				realas += realAsset
				calas += calAsset
				logger.Profitf(",%s,%.4f,%.4f,%.4f,%.4f,,%.4f,%.4f,,%.2f,,%.4f,%.4f", c.Center, c.BTC, c.CNY, c.FBTC, c.FCNY, c.BTC+c.FBTC, c.CNY+c.FCNY, c.CurTicker.Last, realAsset, calAsset)
			}
			c := total
			logger.Profitf(",%s,%.4f,%.4f,%.4f,%.4f,,%.4f,%.4f,,,,%.4f,%.4f", c.Center, c.BTC, c.CNY, c.FBTC, c.FCNY, c.BTC+c.FBTC, c.CNY+c.FCNY, realas, calas)
			logger.Profitf("")
			time.Sleep(time.Hour)
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}
