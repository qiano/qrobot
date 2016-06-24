package provider

import (
	. "github.com/qshuai162/qian/common/model"
	
)

func (p *Provider) GetAccount(center string) (code int, result Account) {
	result, success := getMarketAPI(center).GetAccount()
	if success {
		code = 0
	} else {
		code = 1
	}
	return
}



func (p *Provider) GetOrders(center string) (code int, orders []Order) {
	_, orders = getMarketAPI(center).GetOrders()
	return
}
