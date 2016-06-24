package trader

import (
	"fmt"
	"github.com/qshuai162/qian/common/logger"
	. "github.com/qshuai162/qian/common/model"
	// "strconv"
	"time"
)

type OCache struct {
	Time         int64
	Center 	     string
	Id           int
	RefId        int    //关联订单ID，买单与卖单的对应关系
	Type         string // sell or buy
	Price        float64
	Order_Amount float64 //订单量
	// Process_Amount float64 //成交量
	Status int64 //0-撤销 1-委托 2-全部成交 3-部分成交
}

type OCaches struct {
	List    []OCache //委托中、未卖出的买单
	History []OCache //历史委托（撤销、已成交的卖单、已卖出的买单）
	SaveHistory bool //是否保存历史
}

//新增买单委托
func (os *OCaches) Buy(t int64, center string,id int, price, amount float64) {
	if os.Find(id) == nil {
		o := new(OCache)
		o.Center=center
		o.Time = t
		o.Type = "buy"
		o.Id = id
		o.Price = price
		o.Order_Amount = amount
		o.Status = 1
		os.List = append(os.List, *o)
	} else {

	}
}

//没有卖出的买单
func (os *OCaches) Buy_NotSell() []int {
	rets := make([]int, 0, 0)
	for i := 0; i < len(os.List); i++ {
		temp := os.List[i]
		if temp.Type == "buy"  {
			rets = append(rets, temp.Id)
		}
	}
	for i:=0;i<len(os.History);i++{
		temp := os.History[i]
		if temp.Type == "buy"  {
			if temp.RefId == 0 {
				rets = append(rets, temp.Id)
			} 
		}
	}
	return rets
}

//正在卖的买单（逻辑有错误）
func (os *OCaches) Buy_Selling() ([]int, []int) {
	//逻辑有错误
	sell := make([]int, 0, 0)
	rets := make([]int, 0, 0)
	for i := 0; i < len(os.List); i++ {
		temp := os.List[i]
		if temp.Type == "buy" && temp.Status == 2 && temp.RefId != 0 {
			rets = append(rets, temp.Id)
			sell = append(sell, temp.RefId)
		}
	}
	return rets, sell
}

//委托成交
func (os *OCaches) Done(id int) *OCache{
	temp := os.FindInList(id)
	if temp.Status!=2{
		temp.Status = 2
		if temp.Type=="buy"{
			os.toHistory(id)
		}else{
			if temp.RefId!=0{
				os.RemoveFromHistory(temp.RefId)
			}
			os.RemoveFromList(temp.Id)
		}
		return temp
	}
	return nil
}

func (os *OCaches) RemoveFromList(id int){
	tmp := os.List
	var newList []OCache
	for i := 0; i < len(tmp); i++ {
		if tmp[i].Id != id {
			newList = append(newList, tmp[i])
		}
	}
	os.List = newList
}

func (os *OCaches) RemoveFromHistory(id int) {
	tmp := os.History
	var newList []OCache
	for i := 0; i < len(tmp); i++ {
		if tmp[i].Id != id {
			newList = append(newList, tmp[i])
		}
	}
	os.History = newList
}

//超时需要撤销的委托
func (os *OCaches) WaitForCancel(mintime int64) []int {
	rets := make([]int, 0, 0)
	for i := 0; i < len(os.List); i++ {
		temp := os.List[i]
		if temp.Status == 1 && temp.Time < mintime {
			rets = append(rets, temp.Id)
		}
	}
	return rets
}

func (os *OCaches) toHistory(id int) {
	temp := os.Find(id)
	if os.SaveHistory{
		os.History = append(os.History, *temp)
	}
	tmp := os.List
	var newList []OCache
	for i := 0; i < len(tmp); i++ {
		if tmp[i].Id != temp.Id {
			newList = append(newList, tmp[i])
		}
	}
	os.List = newList
}

//撤销委托
func (os *OCaches) Cancel(id int) {
	temp := os.Find(id)
	temp.Status = 0
	if temp.Type == "sell" {
		if temp.RefId!=0{
			os.Find(temp.RefId).RefId = 0
		}
	}
	os.RemoveFromList(id)
}

//查找委托
func (os *OCaches) Find(id int) *OCache {
	for i := 0; i < len(os.List); i++ {
		if os.List[i].Id == id {
			return &os.List[i]
		}
	}
	for i := 0; i < len(os.History); i++ {
		if os.History[i].Id == id {
			return &os.History[i]
		}
	}
	return nil
}

func(os *OCaches) FindInList(id int) *OCache{
	for i := 0; i < len(os.List); i++ {
		if os.List[i].Id == id {
			return &os.List[i]
		}
	}
	return nil
}

//可以卖出的买单
func (os *OCaches) WaitForSell(maxBuyPrice float64) []int {
	rets := make([]int, 0, 0)
	for i := 0; i < len(os.History); i++ {
		temp := os.History[i]
		if temp.Type == "buy" && temp.Status == 2 && temp.Price <= maxBuyPrice {
			if temp.RefId == 0 {
				rets = append(rets, temp.Id)
			} else {
				sell := os.Find(temp.RefId)
				if sell.Status == 0 {
					rets = append(rets, temp.Id)
				}
			}
		}
	}
	return rets
}

//新增卖单委托
func (os *OCaches) Sell(t int64,center string, buyid, sellid int, price, amount float64) {
	if os.Find(sellid) == nil {
		o := new(OCache)
		o.Center=center
		o.Time = t
		o.Type = "sell"
		o.Id = sellid
		o.RefId = buyid
		o.Price = price
		o.Order_Amount = amount
		o.Status = 1
		os.List = append(os.List, *o)
		if buyid!=0{
			os.Find(buyid).RefId = sellid
		}
	}
}

//与服务器同步委托状态
func (os *OCaches) Sync(orders []Order) {
	if orders != nil {
		for i := 0; i < len(os.List); i++ {
			temp := os.List[i]
			if temp.Status != 2 {
				isdone := true
				for j := 0; j < len(orders); j++ {
					if temp.Id == orders[j].Order_id {
						isdone = false
						break
					}
				}
				if isdone {
					os.Done(temp.Id)
				}
			}
		}
		// os.PrintList()
	}
}

//打印所有委托
func (os *OCaches) PrintAll() {
	os.PrintList()
	os.PrintHistory()
}

func (os *OCaches) PrintList() {
	logger.Infoln("LIST")
	for i := 0; i < len(os.List); i++ {
		print(&os.List[i])
	}
}

func (os *OCaches) PrintHistory() {
	logger.Infoln("HISTORY")
	for i := 0; i < len(os.History); i++ {
		print(&os.History[i])
	}
}

func print(temp *OCache) {
	logger.Infoln(time.Unix(temp.Time, 0).Format("2006-01-02 15:04:05"),temp.Center, temp.Id, temp.RefId, temp.Type, temp.Price, temp.Order_Amount, temp.Status)
}

func (os *OCaches) _PrintAll() {
	fmt.Println("LIST")
	for i := 0; i < len(os.List); i++ {
		_print(&os.List[i])
	}
	fmt.Println("HISTORY")
	for i := 0; i < len(os.History); i++ {
		_print(&os.History[i])
	}
	fmt.Println("")
}

func _print(temp *OCache) {
	fmt.Println(time.Unix(temp.Time, 0).Format("2006-01-02 15:04:05"), temp.Id, temp.RefId, temp.Type, temp.Price, temp.Order_Amount, temp.Status)
}

//需要止损卖出的委托
func (os *OCaches) WaitForStopLoss(nowprice, stoploss float64) []int {

	rets := make([]int, 0, 0)
	for i := 0; i < len(os.History); i++ {
		temp := os.History[i]
		if temp.Type == "buy" && temp.Status == 2 && temp.Price*(1-stoploss) >= nowprice {
			rets = append(rets, temp.Id)
		}
	}
	return rets
}

func (os *OCaches) SimDone(center string,curPrice float64)(dones []*OCache){
	if curPrice<1{
		return
	}
	for i := 0; i < len(os.List); i++ {
		o := os.List[i]
		if o.Center==center{
			if (o.Type == "buy" && o.Price > curPrice) ||
				(o.Type == "sell" && o.Price < curPrice){
					d:=os.Done(o.Id)
					if d!=nil{
						dones=append(dones,d)
					}
			}
		}
	}
	return dones
}

func (os *OCaches) Summary() (int,int) {
	// profit := 0.0
	canceled := 0
	donecount := 0
	for i := 0; i < len(os.History); i++ {
		donecount++
		o := os.History[i]
		if o.Status == 0 {
			canceled++
		}
		// if o.Type == "buy" && o.Status == 2 {
		// 	if o.RefId!=0{
		// 	s := os.Find(o.RefId)
		// 	}
		// 	// profit += s.Order_Amount*s.Price - o.Order_Amount*o.Price
		// }
	}
	return  len(os.List),donecount 
}
