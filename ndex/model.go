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
 * @Date: 2020/5/12 下午3:55
 */
package ndex

type BaseResponse struct {
	Code	int		`"json:code"`
	Success bool	`"json:success"`
	Msg		string	`"json:msg"`
}

type GetTime struct {
	BaseResponse
	Data	int64	`"json:data"`
}

type GetSymbols struct {
	BaseResponse
	Data	[]*Symbol	`"json:data"`
}

type Symbol struct {
	Symbol 				string		`"json:symbol"`		//交易对名称
	BaseAssetName 		string		`"json:baseAssetName"`		//交易资产名称
	BaseDecimal			int			`"json:baseDecimal"`		//交易资产小数位数
	QuoteAssetName 		string		`"json:quoteAssetName"`		//货币资产名称
	QuoteDecimal		int			`"json:quoteDecimal"`		//货币资产小数位数
	BaseMinTradingAmount	float64	`"json:baseMinTradingAmount"`	//最小委托的数量（交易资产）
}

type GetBalance struct {
	BaseResponse
	Data	[]*Balance	`"json:data"`
}

type Balance struct {
	Available		float64		`"json:available"`	//可用金额
	Freeze			float64		`"json:freeze"`		//冻结金额
	AssetName		string		`"json:assetName"`	//资产名称
	Nonce			string 		`"json:nonce"`		//地址的Nonce值
}

type GetOrderBook struct {
	BaseResponse
	Data	*OrderBook	`"json:data"`
}

type OrderBook struct {
	Symbol 			string		`"json:symbol"`
	UpdateTime		int64		`"json:updateTime"`
	SellList		[][]float64	`"json:sellList"`
	BuyList			[][]float64	`"json:buyList"`
}

type GetOpenOrder struct {
	BaseResponse
	Data	[]*Order	`"json:data"`
}

type Order struct {
	Id 				string		`"json:id"`					//订单ID
	Symbol 			string		`"json:symbol"`				//交易对名称
	Address 		string		`"json:address"`			//订单对应的地址
	Type 			int			`"json:type"`				//订单类型，1买，2卖
	BaseAmount		float64		`"json:baseAmount"`			//委托数量
	BaseDealAmount	float64		`"json:baseDealAmount"`		//已成交数量
	Price			float64		`"json:price"`				//委托价格
	AvgPrice		float64		`"json:avgPrice"`			//平均成交价格
	QuoteDealAmount	float64		`"json:quoteDealAmount"`	//已成交金额
	LeftAmount		float64		`"json:leftAmount"`			//未成交数量
	Status			int			`"json:status"`				//委托单状态 1：挂单中，2:部分成交、3:已成交 、4已撤销、5，部分成交已撤单。
	CreateTime		int64		`"json:createTime"`			//创建时间
}

type GetOrderList struct {
	BaseResponse
	Data	*OrderList	`"json:data"`
}

type OrderList struct {
	Data	[]*Order	`"json:data"`
	Total	int			`"json:total"`
}

type GetOrder struct {
	BaseResponse
	Data	*Order	`"json:data"`
}

type WsResponse struct {
	Channel 				string		`"json:channel"`
	Action 					string		`"json:action"`
	Status					int			`"json:status"`
}

type WsPong struct {
	Pong		int64		`"json:pong"`
}

type WsOrderBookResponse struct {
	Channel 				string		`"json:channel"`
	Action 					string		`"json:action"`
	Status					int			`"json:status"`
	Data 					*OrderBook	`"json:data"`
}

type WsOrderChangeResponse struct {
	Channel 				string		`"json:channel"`
	Action 					string		`"json:action"`
	Status					int			`"json:status"`
	Data 					*WsOrderChange	`"json:data"`
}

type WsOrderChange struct {
	T			string		`"json:t"`
	D			[]*Order	`"json:d"`
}

type WsSubInfo struct {
	Channel 	string
	SubMessage	string
	Event		interface{}
}

type NewOrderResponse struct {
	BaseResponse
	Data	string	`"json:data"`
}

type BroadcastResponse struct {
	BaseResponse
	Data	string	`"json:data"`
}