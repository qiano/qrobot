package quantify

import (
	. "github.com/qshuai162/qian/common/config"
	idc "github.com/qshuai162/qian/qrobot/indictor"
	"strconv"
)

type QuantifyOp struct {
	Center     string
	KPeriod    int
	PerAmount  float64
	MinProfit  float64 //最低利润
	MaxCount   int
	StopLoss   float64 //止损百分比
	TimeOut    int     //订单超时时间，单位秒
	LogEnabled bool    //是否开启日志
	Simulation bool    //是否是模拟
}

func GetOpion() QuantifyOp {
	op := new(QuantifyOp)
	op.Center = "okcoin"
	op.LogEnabled = true
	op.Simulation = false
	op.KPeriod = 1
	op.TimeOut = 600
	op.MaxCount = 10
	op.MinProfit = 0.2
	op.PerAmount = 0.01
	op.StopLoss = 0.005
	return *op

}
