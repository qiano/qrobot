package provider

import (
	// "fmt"
	// idc "github.com/qshuai162/qian/qrobot/indictor"
	"testing"
	"time"
)

// func Test_Kline(t *testing.T) {
// 	p := new(Provider)
// 	// start := time.Now().Unix() - 10*60
// 	// end := time.Now().Unix()
// 	// fmt.Println(sart, end)
// 	// for i := 0; i < 3; i++ {
// 	cur := time.Now().Unix()
// 	t1 := time.Now().UnixNano()
// 	_, records := p.GetKLine("okcoin", 1, cur)
// 	fmt.Println(records[len(records)-1])
// 	t2 := time.Now().UnixNano()
// 	_, records = p.GetKLine("okcoin", 5, cur)
// 	macd5 := idc.GetMACD(records, 12, 26, 9)
// 	fmt.Println(records[len(records)-1])
// 	t3 := time.Now().UnixNano()
// 	_, records = p.GetKLine("okcoin", 15, cur)
// 	macd15 := idc.GetMACD(records, 12, 26, 9)
// 	fmt.Println(records[len(records)-1])
// 	t4 := time.Now().UnixNano()
// 	_, records = p.GetKLine("okcoin", 30, cur)
// 	macd30 := idc.GetMACD(records, 12, 26, 9)
// 	fmt.Println(records[len(records)-1])
// 	t5 := time.Now().UnixNano()
// 	_, records = p.GetKLine("okcoin", 60, cur)
// 	macd60 := idc.GetMACD(records, 12, 26, 9)
// 	fmt.Println(records[len(records)-1])
// 	t6 := time.Now().UnixNano()
// 	fmt.Println(t2-t1, t3-t2, t4-t3, t5-t4, t6-t5, len(records))
// 	macdmap := make(map[int][]idc.MACD)
// 	macdmap[5] = macd5
// 	macdmap[15] = macd15
// 	macdmap[30] = macd30
// 	macdmap[60] = macd60
// 	fmt.Println(macd5[len(macd5)-1])
// 	fmt.Println(macd15[len(macd15)-1])
// 	fmt.Println(macd30[len(macd30)-1])
// 	fmt.Println(macd60[len(macd60)-1])
// 	// trend := idc.GetTrend(macdmap)
// 	// fmt.Println(trend)
// 	// time.Sleep((30 * time.Second))
// 	// }
// 	// fmt.Println(k[length-2], j[length-2], d[length-2])
// 	// fmt.Println(k[length-1], j[length-1], d[length-1])
// }

func Test_Ticker(t *testing.T) {
	p := new(Provider)
	p.GetDifAvg("okcoin", "chbtc", time.Now().Unix(), 4800)

}
