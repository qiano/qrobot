package indictor

import (
	. "github.com/qshuai162/qian/common/model"
)

type kdj struct {
	k float64
	d float64
	j float64
}

var kdjcaches map[int64]kdj = make(map[int64]kdj)

func GetKDJ(records []Record, period int) (arrk []float64, arrd []float64, arrj []float64) {
	length := len(records)
	for i := 0; i < length; i++ {
		k, d, j := float64(0), float64(0), float64(0)
		val, ok := kdjcaches[records[i].Time]
		if ok {
			k = val.k
			d = val.d
			j = val.j
		} else {
			var periodLowArr, periodHighArr []float64
			rsv := 0.0
			lowest := 0.0
			highest := 0.0
			if i >= period-1 {
				for j := 0; j < period; j++ {
					periodLowArr = append(periodLowArr, records[i-j].Low)
					periodHighArr = append(periodHighArr, records[i-j].High)
				}
				lowest = arrayLowest(periodLowArr)
				highest = arrayHighest(periodHighArr)
				rsv = float64(100) * (records[i].Close - lowest) / (highest - lowest)
			}
			if i == 0 {
				k = (2.0/3)*0 + 1.0/3*rsv
				d = (2.0/3)*0 + 1.0/3*k
			} else {
				k = (2.0/3)*arrk[i-1] + 1.0/3*rsv
				d = (2.0/3)*arrd[i-1] + 1.0/3*k
			}
			j = 3*k - 2*d

			if k < 0 {
				k = 0
			}
			if k > 100 {
				k = 100
			}
			if d < 0 {
				d = 0
			}
			if d > 100 {
				d = 100
			}
			if j < 0 {
				j = 0
			}
			if j > 100 {
				j = 100
			}

			if i != length-1 {
				kdj := new(kdj)
				kdj.k = k
				kdj.d = d
				kdj.j = j
				kdjcaches[records[i].Time] = *kdj
			}
		}
		arrk = append(arrk, k)
		arrd = append(arrd, d)
		arrj = append(arrj, j)
	}
	return
}
