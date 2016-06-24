package provider

import (
	// "fmt"
	. "github.com/qshuai162/qian/common/model"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"math"
	"time"
)

type KCache struct {
	Time    int64
	Records []Record
}

var k1Cache KCache = *new(KCache)

func (p *Provider) GetKLine(center string, period int, end int64) (int, []Record) {
	var k1datas []Record
	if k1Cache.Time == end {
		k1datas = k1Cache.Records
	} else {
		_, k1datas = getMarketAPI(center).GetKLine(1, 60*48)
		k1Cache.Time = end
		k1Cache.Records = k1datas
	}
	records := createKByK1(k1datas, period)
	return 0, records
}

func createKByK1(k1data []Record, period int) (records []Record) {
	if period == 1 {
		return k1data
	}
	start := int64(math.Ceil(float64(k1data[0].Time)/float64(60*period)) * float64(60*period))
	length := len(k1data)
	for i := 0; i < length; i++ {
		if k1data[i].Time < start {
			continue
		} else {
			record := new(Record)
			record.TimeStr = time.Unix(k1data[i].Time, 0).Format("2006-01-02 15:04:05")
			record.Time = k1data[i].Time
			record.Open = k1data[i].Open
			record.Volumn = k1data[i].Volumn
			record.Close = k1data[i].Close
			record.High = k1data[i].High
			record.Low = k1data[i].Low
			for j := 1; j < period && i+j < length; j++ {
				k1 := k1data[i+j]
				if k1.Open != 0 {
					record.Volumn += k1.Volumn
					record.Close = k1.Close
					if k1.High > record.High {
						record.High = k1.High
					}
					if k1.Low < record.Low {
						record.Low = k1.Low
					}
				}
			}
			i += period - 1
			records = append(records, *record)
		}
	}
	return
}
