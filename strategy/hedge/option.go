package hedge

import (
// . "github.com/qshuai162/qian/common/config"
// "strconv"
)

type HedgeOp struct {
	Centers      []string
	HdConfig     [][]float64
	RestTime     int     //休息时间，单位秒
	Fast    	 float64
	ColdTime     int
	TooFast      float64 //价格变化过快判定值
	ColdDownTime int     //价格变化过快等待时间 ，单位秒
	Refreshtime  int     //中值刷新时间，单位秒
	LogEnabled   bool    //是否开启日志
	Simulation   bool    //是否是模拟
}

func GetOpion() HedgeOp {
	op := new(HedgeOp)
	op.RestTime = 3
	op.Fast=0.7
	op.ColdTime=20
	op.TooFast = 1.8
	op.ColdDownTime = 360
	op.Refreshtime = 300
	op.LogEnabled = true
	op.Simulation = false
	return *op
}
