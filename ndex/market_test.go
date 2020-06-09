/**
* MIT License
* <p>
Copyright (c) 2019-2020 nerve.network
* <p>
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
* <p>
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
* <p>
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
* SOFTWARE.
*/

/**
 * @Author: nerve.network core team
 * @Date: 2020/5/12 下午4:01
 */
package ndex

import (
	"fmt"
	"log"
	"testing"
)

var (
	market *Market
)

func init() {
	host := ""
	market = &Market{Host: host}
	market.Initialize()
}

func TestMarket_GetServeTime(t *testing.T) {
	tm, err := market.GetServeTime()
	if err != nil {
		t.Error(err)
		return
	}
	if tm == nil {
		t.Fail()
		return
	}
	fmt.Println(tm)
}

func TestMarket_GetSymbols(t *testing.T) {
	symbols, err := market.GetSymbols()
	if err != nil {
		t.Error(err)
		return
	}
	if symbols == nil {
		t.Fail()
		return
	}
	for _, symbol := range symbols {
		fmt.Println(symbol)
	}
}

func TestMarket_GetBalance(t *testing.T) {
	address := "TNVTdN9iCdXS46SuN8UPd2kkdXmGNYgMMxuSJ"
	balances, err := market.GetBalanceByAddress(address)
	if err != nil {
		t.Error(err)
		return
	}
	if balances == nil {
		t.Fail()
		return
	}
	for _, balance := range balances {
		fmt.Println(balance)
	}
}

func TestMarket_GetOrderBook(t *testing.T) {
	symbol := "BTCUSDT"
	size := 10
	orderBook, err := market.GetOrderBook(symbol, size)
	if err != nil {
		t.Error(err)
		return
	}
	if orderBook == nil {
		t.Fail()
		return
	}
	fmt.Println(orderBook)
}

func TestMarket_GetOpenOrder(t *testing.T) {
	address := "TNVTdN9iCdXS46SuN8UPd2kkdXmGNYgMMxuSJ"
	symbol := "BTCUSDT"
	openOrders, err := market.GetOpenOrderByAddress(address, symbol)
	if err != nil {
		t.Error(err)
		return
	}
	if openOrders == nil {
		t.Fail()
		return
	}
	fmt.Println(openOrders)
	fmt.Println("total order count is : ", len(openOrders))
	for _, order := range openOrders {
		fmt.Printf("%#v\n", order)
	}
}

func TestMarket_GetOrderList(t *testing.T) {
	address := "TNVTdN9iCdXS46SuN8UPd2kkdXmGNYgMMxuSJ"
	symbol := "BTCUSDT"
	orderList, err := market.GetOrderListByAddress(address, symbol, 5000700, 10)
	if err != nil {
		t.Error(err)
		return
	}
	if orderList == nil {
		t.Fail()
		return
	}
	fmt.Println(orderList)
	fmt.Println("total order count is : ", len(orderList.Data))
	for _, order := range orderList.Data {
		fmt.Printf("%#v\n", order)
	}
}

func TestMarket_GetOrder(t *testing.T) {
	id := "b0113ba5efb8b3c0a01c7829e6b4e4e775437af7334a3bb02eb33b6db9beaee7"
	order, err := market.GetOrder(id)
	if err != nil {
		t.Error(err)
		return
	}
	if order == nil {
		t.Fail()
		return
	}
	fmt.Printf("%#v\n", order)
}

func TestMarket_NewOrder(t *testing.T) {
	order, err := market.NewOrderByAddress("TNVTdN9iCdXS46SuN8UPd2kkdXmGNYgMMxuSJ", "1b0470a2a8c8a02c5dee364fd6d0f56dc7e30a4a817c18b4b00c419c349dc7df", "BTCUSDT", 1, 8000, 1)
	if err != nil {
		t.Error(err)
		return
	}
	if order == nil {
		t.Fail()
		return
	}
	fmt.Println(order)

}

func TestMarket_CancelOrder(t *testing.T) {
	txId, err := market.CancelOrderByAddress("82ad1756b22f480195c3f7eee5dfe2cf43e661b404f828255903cedb70172e40", "1b0470a2a8c8a02c5dee364fd6d0f56dc7e30a4a817c18b4b00c419c349dc7df")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(txId)

}

func TestMarket_SubscribeOrderBook(t *testing.T) {
	orderBookEvent, err := market.SubscribeOrderBook("BTCUSDT", 10)
	if err == nil {
		t.Log("subscribe success")
	}
	for {
		select {
		case orderBook := <- orderBookEvent:
			log.Println(orderBook)
		}
	}
}

func TestMarket_SubscribeOrderChange(t *testing.T) {
	orderChangeEvent, err :=  market.SubscribeOrderChangeByAddress("TNVTdN9iCdXS46SuN8UPd2kkdXmGNYgMMxuSJ")
	if err == nil {
		t.Log("subscribe success")
	}
	for {
		select {
		case order := <- orderChangeEvent:
			log.Printf("%#v\n", order)
		}
	}
}