package indictor

import (
// "math"
)

func getEMAdifAt(emaShort, emaLong []float64, idx int) float64 {
	var cel = emaLong[idx]
	var ces = emaShort[idx]
	if cel == 0 {
		return 0
	} else {
		return 100 * (ces - cel) / ((ces + cel) / 2)
	}
}

func getEMAdif(emaShort, emaLong []float64) []float64 {
	// loop through data
	var EMAdifs []float64
	length := len(emaShort)
	for i := 0; i < length; i++ {
		EMAdifAt := getEMAdifAt(emaShort, emaLong, i)
		EMAdifs = append(EMAdifs, EMAdifAt)
	}

	return EMAdifs
}

/* Function based on the idea of an exponential moving average.
 *
 * Formula: EMA = Price(t) * k + EMA(y) * (1 - k)
 * t = today y = yesterday N = number of days in EMA k = 2/(2N+1)
 *
 * @param Price : array of y variables.
 * @param periods : The amount of "days" to average from.
 * @return an array containing the EMA.
**/
func EMA(Price []float64, periods int) []float64 {
	var t float64
	y := 0.0
	n := float64(periods)
	var k float64
	k = 2 / (n + 1)
	var ema float64 // exponential moving average.

	var periodArr []float64
	var startpos int
	length := len(Price)
	var emaLine []float64 = make([]float64, length)

	// loop through data
	for i := 0; i < length; i++ {
		if Price[i] != 0 {
			startpos = i + 1
			break
		} else {
			emaLine[i] = 0
		}
	}

	for i := startpos; i < length; i++ {
		periodArr = append(periodArr, Price[i])

		// 0: runs if the periodArr has enough points.
		// 1: set currentvalue (today).
		// 2: set last value. either by past avg or yesterdays ema.
		// 3: calculate todays ema.
		if periods == len(periodArr) {

			t = Price[i]

			if y == 0 {
				y = arrayAvg(periodArr)
			} else {
				ema = (t * k) + (y * (1 - k))
				y = ema
			}
			//四舍五入保留4位小数
			//value := math.Trunc(y*1e4+0.5) * 1e-4
			emaLine[i] = y

			// remove first value in array.
			periodArr = periodArr[1:]

		} else {

			emaLine[i] = 0
		}

	}

	return emaLine
}

/* Function that returns average of an array's values.
 *
**/
func arrayAvg(arr []float64) float64 {
	sum := 0.0

	for i := 0; i < len(arr); i++ {
		sum = sum + arr[i]
	}

	return (sum / (float64)(len(arr)))
}

func Highest(Price []float64, periods int) []float64 {
	var periodArr []float64
	length := len(Price)
	var HighestLine []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, Price[i])
		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.
		if periods == len(periodArr) {
			HighestLine[i] = arrayHighest(periodArr)

			// remove first value in array.
			periodArr = periodArr[1:]
		} else {
			HighestLine[i] = 0
		}
	}

	return HighestLine
}

func Lowest(Price []float64, periods int) []float64 {
	var periodArr []float64
	length := len(Price)
	var LowestLine []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, Price[i])
		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.
		if periods == len(periodArr) {
			LowestLine[i] = arrayLowest(periodArr)

			// remove first value in array.
			periodArr = periodArr[1:]
		} else {
			LowestLine[i] = 0
		}
	}

	return LowestLine
}

/* Function based on the idea of a simple moving average.
 * @param Price : array of y variables.
 * @param periods : The amount of "days" to average from.
 * @return an array containing the SMA.
**/
func SMA(Price []float64, periods int) []float64 {
	var periodArr []float64
	length := len(Price)
	var smLine []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, Price[i])

		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.
		if periods == len(periodArr) {
			smLine[i] = arrayAvg(periodArr)

			// remove first value in array.
			periodArr = periodArr[1:]
		} else {
			smLine[i] = 0
		}
	}

	return smLine
}

func arrayLowest(Price []float64) float64 {
	length := len(Price)
	var lowest = Price[0]

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		if Price[i] < lowest {
			lowest = Price[i]
		}
	}

	return lowest
}

func arrayHighest(Price []float64) float64 {
	length := len(Price)
	var highest = Price[0]

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		if Price[i] > highest {
			highest = Price[i]
		}
	}

	return highest
}

func getMACDdifAt(emaShort, emaLong []float64, idx int) float64 {
	var ces = emaShort[idx]
	var cel = emaLong[idx]
	if cel == 0 {
		return 0
	} else {
		return (ces - cel)
	}
}

func getMACDdif(emaShort, emaLong []float64) []float64 {
	// loop through data
	var MACDdif []float64
	length := len(emaShort)
	for i := 0; i < length; i++ {
		MACDdifAt := getMACDdifAt(emaShort, emaLong, i)
		MACDdif = append(MACDdif, MACDdifAt)
	}

	return MACDdif
}

func getMACDSignal(MACDdif []float64, signalPeriod int) []float64 {
	signal := EMA(MACDdif, signalPeriod)
	return signal
}

func getMACDHistogramAt(MACDdif, MACDSignal []float64, idx int) float64 {
	var dif = MACDdif[idx]
	var signal = MACDSignal[idx]
	if signal == 0 {
		return 0
	} else {
		return dif - signal
	}
}
