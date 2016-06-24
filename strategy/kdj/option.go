package kdj

import (
)

type KdjOp struct {
	Center     string
	KPeriod    int
	PerAmount  float64
	MinProfit  float64 //最低利润
	JHigh      float64
	MaxCount   int
	StopLoss   float64 //止损百分比
	TimeOut    int     //订单超时时间，单位秒
	LogEnabled bool    //是否开启日志
	Simulation bool    //是否是模拟
}

func GetOpion() KdjOp {
	op := new(KdjOp)
	op.Center = "okcoin"
	op.LogEnabled = true
	op.Simulation = false
	op.KPeriod = 1
	op.TimeOut = 300
	op.JHigh = 20
	op.MaxCount = 10
	op.MinProfit=0.2
	op.PerAmount=0.01
	op.StopLoss = 0.005
	return *op

}
