package indictor

import (
// "fmt"
// . "github.com/qshuai162/qian/qrobot/trader"
// "strconv"
// "strings"
)

type TrendOp struct {
	Up5    float64
	Up15   float64
	Up30   float64
	Up60   float64
	Down5  float64
	Down15 float64
	Down30 float64
	Down60 float64
	Offset float64
}

func GetTrend(macdmap map[int][]MACD, op TrendOp) float64 {
	factor := 0.0
	periods := []int{5, 15, 30, 60}
	for _, period := range periods {
		macd := macdmap[period]
		length := len(macd)
		if macd[length-1].BAR > 0 {
			//上升趋势
			switch period {
			case 5:
				factor += op.Up5
			case 15:
				factor += op.Up15
			case 30:
				factor += op.Up30
			case 60:
				factor += op.Up60
			default:
				factor += 0
			}
		} else {
			switch period {
			case 5:
				factor += op.Down5
			case 15:
				factor += op.Down15
			case 30:
				factor += op.Down30
			case 60:
				factor += op.Down60
			default:
				factor += 0
			}
		}
	}
	return factor + op.Offset
}
