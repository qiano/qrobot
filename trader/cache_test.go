package trader

import (
	"fmt"
	// "strconv"
	"testing"
	"time"
)

func Test_Cache(t *testing.T) {
	os := new(OCaches)
	os.SaveHistory=true
tt:=time.Now().Unix()
	fmt.Println("buy 1")
	os.Buy(tt,"test",1, 1456.12, 0.002)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)
	fmt.Println("cancel 1")
	os.Cancel(1)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("buy 2")
	os.Buy(tt,"test",2, 321.12, 22.2)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("done 2")
	os.Done(2)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("sell 2 3")
	os.Sell(tt,"test",2, 3, 12, 32)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("cancel 3")
	os.Cancel(3)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("buy 4")
	os.Buy(tt,"test",4, 11, 1)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("done 4")
	os.Done(4)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("sell 40")
	os.Sell(tt,"test",4, 40, 123, 1)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("cancel 40")
	os.Cancel(40)
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("buy not sell")
	t2 := os.Buy_NotSell()
	for i := 0; i < len(t2); i++ {
		t := os.Find(t2[i])
		_print(t)
	}
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("wait for sell and sell")
	temp := os.WaitForSell(9999)
	for i := 0; i < len(temp); i++ {
		t := os.Find(temp[i])
		_print(t)
		os.Sell(tt,"test",t.Id, t.Id+1000, 9999, os.Find(temp[i]).Order_Amount)
	}
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("wait for cancel and cancel")
	temp1 := os.WaitForCancel(100000000000)
	for i := 0; i < len(temp1); i++ {
		t := os.Find(temp1[i])
		_print(t)
		os.Cancel(t.Id)
	}
	os._PrintAll()
	time.Sleep(500 * time.Millisecond)

}
