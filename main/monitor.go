package main

import (
	"fmt"
	// . "github.com/qshuai162/qian/common/model"
	"github.com/qshuai162/qian/qrobot/provider"
	"github.com/qshuai162/qian/qrobot/strategy/hedge"
	// "github.com/qshuai162/qian/qrobot/strategy/kdj"
	"github.com/qshuai162/qian/qrobot/trader"
	// . "github.com/qshuai162/qian/tool/config"
	"github.com/qshuai162/qian/common/logger"
	// "strconv"
	"time"
) 

func RobotWorker() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	fmt.Println("trade robot start working...")

	Centers := []string{"okcoin", "chbtc", "huobi","btcc"}
	p := new(provider.Provider)
	t := new(trader.Trader)

	op := hedge.GetOpion()
	op.Centers = Centers
	op.HdConfig = [][]float64{
		[]float64{1, 0, 0.55, 0.05},
		[]float64{2, 0, 0.55, 0.05},
		[]float64{1, 2, 0.55, 0.05},
		[]float64{3, 0, 0.55, 0.05},
		[]float64{1, 3, 0.55, 0.05},
		[]float64{3, 2, 0.55, 0.05}}
	// op.HdConfig=[][]float64{
	// 	[]float64{0,1,0.4,0.05,500} }

	// kdjop:=kdj.GetOpion()
	
	now:=time.Now().Unix();
	hd := hedge.NewHedge(now, op, p, t)
	// kdjs:=kdj.NewKdjStrategy(now,kdjop,hd)

	go func() {
		for _ = range ticker.C {
			curt := time.Now().Unix()
			logger.Infoln("")
			logger.Infof("%s\n", time.Unix(curt, 0).Format("2006-01-02 15:04:05"))
			//准备数据
			cc := make(chan int)
			// var records []Record
			go func(){
				hd.HandleTicker(curt)
				cc<-1
			}()
			go func(){
				// _, records = p.GetKLine(kdjop.Center, kdjop.KPeriod, curt)
				cc<-1
			}()
			count := 0
			for {
				count += <-cc
				if count >= 2 {
					break
				}
			}
			go hd.Tick(curt)
			// go func(){
			// 	if len(records)>0 && hd.Centers[0].CurTicker.Last>0.1{
			// 		kdjs.Tick(curt,records,hd.Centers[0].CurTicker)
			// 	}
			// }()
		}
	}()

	time.Sleep(24 * 365 * 100 * time.Hour)
}

const worker_number = 1

type message struct {
	normal bool                   // true means exit normal, otherwise
	state  map[string]interface{} // goroutine state
}

func worker(mess chan message) {
	defer func() {
		exit_message := message{state: make(map[string]interface{})}
		i := recover()
		if i != nil {
			exit_message.normal = false
		} else {
			exit_message.normal = true
		}
		mess <- exit_message
	}()

	RobotWorker()
}

func supervisor(mess chan message) {
	for i := 0; i < worker_number; i++ {
		m := <-mess
		switch m.normal {
		case true:
			logger.Infoln("exit normal, nothing serious!")
		case false:
			logger.Infoln("exit abnormal, something went wrong")
		}
	}
}

func RunRobot() {
	mess := make(chan message, 10)
	for i := 0; i < worker_number; i++ {
		go worker(mess)
	}

	supervisor(mess)
}
