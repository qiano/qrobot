package hedge

import(
    "time"
    "strconv"
)

func (hd *Hedge) SyncAccount(){
    	for i := 0; i < len(hd.Centers); i++ {
			if hd.Centers[i].Center != "" {
				code2, ok2 := (*hd.Provider).GetAccount(hd.Centers[i].Center)
				if code2 == 0 {
                    hd.Centers[i].Btc,_=strconv.ParseFloat(ok2.Available_btc,64)
                    hd.Centers[i].Cny,_=strconv.ParseFloat(ok2.Available_cny,64)
				}
			}
		}
}

func (hd *Hedge) startSyncAccount() {
	for {
		time.Sleep(5 * time.Second)
	    hd.SyncAccount()
	}
}