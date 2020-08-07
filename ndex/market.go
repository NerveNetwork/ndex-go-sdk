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
 * @Date: 2020/5/12 下午3:43
 */
package ndex

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/niels1286/nuls-go-sdk/crypto/eckey"
	txprotocal "github.com/niels1286/nuls-go-sdk/tx/protocal"
	"github.com/niels1286/nuls-go-sdk/utils/seria"
	"log"
	"time"

	"github.com/NerveNetwork/ndex-go-sdk/utils"
)

type Market struct {
	Host 		string
	WsHost		string
	Address 	string
	PrivateKey	string

	ndexWs		*NdexWs
}

/**
 * Perform basic testing and set the default server access address
 * 进行基本的检测，设置默认的服务器访问地址
 */
func (market *Market) Initialize() {
	if market.Host == "" {
		market.Host = "https://api.nervedex.com"
	}
	if market.WsHost == "" {
		market.WsHost = "wss://api.nervedex.com"
	}
}

/**
 * Get server time
 * 获取服务器时间
 */
func (market *Market) GetServeTime() (*time.Time, error) {
	uri := "/api/time"
	url := market.Host + uri
	responseBytes, err := utils.RequestGet(url)
	if err != nil {
		return nil, err
	}
	// Parsing the return value 解析返回值
	getTime := &GetTime{}
	err = json.Unmarshal(responseBytes, getTime)
	if err != nil {
		return nil, err
	}
	if !getTime.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", getTime.Code, getTime.Msg))
	}
	thisTime := time.Unix(getTime.Data / 1000, getTime.Data % 1000 * 1000)
	return &thisTime, nil
}

/**
 * Get information about all trading pairs in the trading market
 * 获取交易市场的所有交易对信息
 */
func (market *Market) GetSymbols() ([]*Symbol, error) {
	uri := "/api/tradings"
	url := market.Host + uri
	responseBytes, err := utils.RequestGet(url)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(responseBytes))
	// Parsing the return value 解析返回值
	getSymbols := &GetSymbols{}
	err = json.Unmarshal(responseBytes, getSymbols)
	if err != nil {
		return nil, err
	}
	if !getSymbols.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", getSymbols.Code, getSymbols.Msg))
	}
	return getSymbols.Data, nil
}

/**
 * Get ticker information of trading pairs
 * 获取交易对的ticker信息
 */
func (market *Market) GetTicker(symbol string) (*Ticker, error) {
	if symbol == "" {
		return nil, errors.New("symbol can not empty")
	}
	uri := fmt.Sprintf("/api/ticker/%s", symbol)
	url := market.Host + uri
	responseBytes, err := utils.RequestGet(url)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(responseBytes))
	// Parsing the return value 解析返回值
	getTicker := &GetTicker{}
	err = json.Unmarshal(responseBytes, getTicker)
	if err != nil {
		return nil, err
	}
	if !getTicker.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", getTicker.Code, getTicker.Msg))
	}
	return getTicker.Data, nil
}

/**
 * Get information about all trading pairs in the trading market
 * 获取交易市场的所有交易对信息
 */
func (market *Market) Kline(symbol string, inv, size int) ([]*Kline, error) {
	uri := "/api/kline"
	url := market.Host + uri

	params := map[string]interface{} {
		"symbol":symbol,
		"type":inv,
		"limit":size,
	}
	responseBytes, err := utils.RequestHttpGet(url, params)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(responseBytes))
	// Parsing the return value 解析返回值
	getKline := &GetKline{}
	err = json.Unmarshal(responseBytes, getKline)
	if err != nil {
		return nil, err
	}
	if !getKline.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", getKline.Code, getKline.Msg))
	}
	return getKline.Data, nil
}

/**
 * Get the asset balance of the configured address
 * 获取配置地址的资产余额
 */
func (market *Market) GetBalance() ([]*Balance, error) {
	if market.Address == "" {
		return nil, errors.New("No address is configured")
	}
	return market.GetBalanceByAddress(market.Address)
}

/**
 * Get the asset balance of the specified address
 * 获取指定地址的资产余额
 */
func (market *Market) GetBalanceByAddress(address string) ([]*Balance, error) {
	if address == "" {
		return nil, errors.New("address can not empty")
	}
	uri := fmt.Sprintf("/api/ledger/%s", address)
	url := market.Host + uri
	responseBytes, err := utils.RequestGet(url)
	if err != nil {
		return nil, err
	}
	// Parsing the return value 解析返回值
	getBalance := &GetBalance{}
	err = json.Unmarshal(responseBytes, getBalance)
	if err != nil {
		return nil, err
	}
	if !getBalance.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", getBalance.Code, getBalance.Msg))
	}
	return getBalance.Data, nil
}

/**
 * Get the market information of the specified trading pair
 * 获取指定交易对的盘口信息
 */
func (market *Market) GetOrderBook(symbol string, size int) (*OrderBook, error) {
	if symbol == "" {
		return nil, errors.New("symbol can not empty")
	}
	uri := fmt.Sprintf("/api/orderBook/%s/%d", symbol, size)
	url := market.Host + uri
	responseBytes, err := utils.RequestGet(url)
	if err != nil {
		return nil, err
	}
	// Parsing the return value 解析返回值
	getOrderBook := &GetOrderBook{}
	err = json.Unmarshal(responseBytes, getOrderBook)
	if err != nil {
		return nil, err
	}
	if !getOrderBook.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", getOrderBook.Code, getOrderBook.Msg))
	}
	return getOrderBook.Data, nil
}

/**
 * Obtain the pending order of the trading pair corresponding to the configured address
 * 获取配置地址对应交易对下的挂单
 */
func (market *Market) GetOpenOrder(symbol string) ([]*Order, error) {
	if market.Address == "" {
		return nil, errors.New("No address is configured")
	}
	return market.GetOpenOrderByAddress(market.Address, symbol)
}


/**
 * Obtain the pending order of the corresponding transaction pair at the specified address
 * 获取指定地址对应交易对下的挂单
 */
func (market *Market) GetOpenOrderByAddress(address, symbol string) ([]*Order, error) {
	if address == "" {
		return nil, errors.New("address can not empty")
	}
	if symbol == "" {
		return nil, errors.New("symbol can not empty")
	}
	uri := fmt.Sprintf("/api/openOrder/%s/%s", symbol, address)
	url := market.Host + uri
	responseBytes, err := utils.RequestGet(url)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(responseBytes))
	// Parsing the return value 解析返回值
	getOpenOrder := &GetOpenOrder{}
	err = json.Unmarshal(responseBytes, getOpenOrder)
	if err != nil {
		return nil, err
	}
	if !getOpenOrder.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", getOpenOrder.Code, getOpenOrder.Msg))
	}
	return getOpenOrder.Data, nil
}

/**
 * Get the list of orders (including unfilled pending orders) corresponding to the trading address of the configured address
 * 获取配置地址对应交易对的订单列表（包括未成交的挂单）
 */
func (market *Market) GetOrderList(symbol string, pageNumber, pageSize int) (*OrderList, error) {
	if market.Address == "" {
		return nil, errors.New("No address is configured")
	}
	return market.GetOrderListByAddress(market.Address, symbol, pageNumber, pageSize)
}

/**
 * Get a list of orders (including unfilled pending orders) corresponding to the specified address
 * 获取指定地址对应交易对的订单列表（包括未成交的挂单）
 */
func (market *Market) GetOrderListByAddress(address, symbol string, pageNumber, pageSize int) (*OrderList, error) {
	if address == "" {
		return nil, errors.New("address can not empty")
	}
	if symbol == "" {
		return nil, errors.New("symbol can not empty")
	}
	url := market.Host + "/api/order/list"
	params := map[string]interface{} {
		"address":address,
		"symbol":symbol,
		"pageNumber":pageNumber,
		"pageSize":pageSize,
	}
	responseBytes, err := utils.RequestPost(url, params)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(responseBytes))

	// Parsing the return value 解析返回值
	getOrderList := &GetOrderList{}
	err = json.Unmarshal(responseBytes, getOrderList)
	if err != nil {
		return nil, err
	}
	if !getOrderList.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", getOrderList.Code, getOrderList.Msg))
	}
	//fmt.Println(getOrderList.Data)
	return getOrderList.Data, nil
}

/**
 * Get order details based on order ID
 * 根据订单ID获取订单详情信息
 */
func (market *Market) GetOrder(id string) (*Order, error) {
	if id == "" {
		return nil, errors.New("order id can not empty")
	}
	uri := fmt.Sprintf("/api/order/%s", id)
	url := market.Host + uri
	responseBytes, err := utils.RequestGet(url)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(responseBytes))
	// Parsing the return value 解析返回值
	getOrder := &GetOrder{}
	err = json.Unmarshal(responseBytes, getOrder)
	if err != nil {
		return nil, err
	}
	if !getOrder.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", getOrder.Code, getOrder.Msg))
	}
	return getOrder.Data, nil
}

/**
 * new order
 * 下单
 */
func (market *Market) NewOrder(symbol string, slide int, price, quantity float64) (*Order, error) {
	if market.Address == "" {
		return nil, errors.New("No address is configured")
	}
	if market.PrivateKey == "" {
		return nil, errors.New("No privateKey is configured")
	}
	return market.NewOrderByAddress(market.Address, market.PrivateKey, symbol, slide, price, quantity)
}

/**
 * new order
 * 下单
 */
func (market *Market) NewOrderByAddress(address, privateKey, symbol string, slide int, price, quantity float64) (*Order, error) {
	if symbol == "" {
		return nil, errors.New("symbol can not empty")
	}
	if address == "" {
		return nil, errors.New("address can not empty")
	}
	if privateKey == "" {
		return nil, errors.New("privateKey can not empty")
	}
	url := market.Host + "/api/order"
	params := map[string]interface{} {
		"address":address,
		"symbol":symbol,
		"quantity":quantity,
		"price":price,
		"type":slide,
	}
	responseBytes, err := utils.RequestPost(url, params)
	if err != nil {
		return nil, err
	}

	newOrderResponse := &NewOrderResponse{}
	err = json.Unmarshal(responseBytes, newOrderResponse)
	if err != nil {
		return nil, err
	}
	if !newOrderResponse.Success {
		return nil, errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", newOrderResponse.Code, newOrderResponse.Msg))
	}
	txBytes, err := hex.DecodeString(newOrderResponse.Data)
	if err != nil {
		return nil, err
	}
	tx := txprotocal.ParseTransactionByReader(seria.NewByteBufReader(txBytes, 0))
	//sign
	hash, err := tx.GetHash().Serialize()
	if err != nil {
		return nil, err
	}
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	ecKey, err := eckey.FromPriKeyBytes(privateKeyBytes)
	if err != nil {
		log.Println("New Order ERROR : private key error, ", err)
		return nil, err
	}
	signData, err := ecKey.Sign(hash)
	if err != nil {
		return nil, err
	}
	sign := txprotocal.P2PHKSignature{
		SignValue: signData,
		PublicKey: ecKey.GetPubKeyBytes(true),
	}
	writer := seria.NewByteBufWriter()
	writer.WriteBytesWithLen(sign.PublicKey)
	writer.WriteBytesWithLen(sign.SignValue)
	tx.SignData = writer.Serialize()
	// broadcast tx
	txBytes, err = tx.Serialize()
	if err != nil {
		return nil, err
	}
	txHash, err := market.broadcast(hex.EncodeToString(txBytes))
	if err != nil {
		return nil, err
	}
	order := &Order{
		Id: txHash,
		Address: address,
		Symbol: symbol,
		Price: price,
		Type: slide,
		BaseAmount: quantity,
		Status: 1,
	}
	return order, nil
}

/**
 * Cancel the order, note that the configured private key must match the address corresponding to the order
 * 取消订单, 注意，配置的私钥必须和订单对应的地址匹配
 */
func (market *Market) CancelOrder(orderId string) (string, error) {
	if market.Address == "" {
		return "", errors.New("No address is configured")
	}
	if market.PrivateKey == "" {
		return "", errors.New("No privateKey is configured")
	}
	return market.CancelOrderByAddress(orderId, market.PrivateKey)
}

/**
 * Cancel the order, note that the incoming private key must match the address corresponding to the order
 * 取消订单, 注意，传入的私钥必须和订单对应的地址匹配
 */
func (market *Market) CancelOrderByAddress(orderId, privateKey string) (string, error) {
	if orderId == "" {
		return "", errors.New("orderId can not empty")
	}
	url := market.Host + "/api/cancelOrder"
	params := map[string]interface{} {
		"orderId":orderId,
	}
	responseBytes, err := utils.RequestPost(url, params)
	if err != nil {
		return "", err
	}
	newOrderResponse := &NewOrderResponse{}
	err = json.Unmarshal(responseBytes, newOrderResponse)
	if err != nil {
		return "", err
	}
	if !newOrderResponse.Success {
		return "", errors.New(fmt.Sprintf("the server return false, code=%d , msg=%s", newOrderResponse.Code, newOrderResponse.Msg))
	}
	// sign tx
	txBytes, err := hex.DecodeString(newOrderResponse.Data)
	if err != nil {
		return "", err
	}
	tx := txprotocal.ParseTransactionByReader(seria.NewByteBufReader(txBytes, 0))
	hash, err := tx.GetHash().Serialize()
	if err != nil {
		return "", err
	}
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	ecKey, err := eckey.FromPriKeyBytes(privateKeyBytes)
	if err != nil {
		log.Println("New Order ERROR : private key error, ", err)
		return "", err
	}
	signData, err := ecKey.Sign(hash)
	if err != nil {
		return "", err
	}
	sign := txprotocal.P2PHKSignature{
		SignValue: signData,
		PublicKey: ecKey.GetPubKeyBytes(true),
	}
	writer := seria.NewByteBufWriter()
	writer.WriteBytesWithLen(sign.PublicKey)
	writer.WriteBytesWithLen(sign.SignValue)
	tx.SignData = writer.Serialize()
	// broadcast tx
	txBytes, err = tx.Serialize()
	if err != nil {
		return "", err
	}
	txHash, err := market.broadcast(hex.EncodeToString(txBytes))
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (market *Market) broadcast(txHex string) (string, error) {
	url := market.Host + "/api/broadcast"
	params := map[string]interface{} {
		"txHex":txHex,
	}
	responseBytes, err := utils.RequestPost(url, params)
	if err != nil {
		return "", err
	}
	broadcastResponse := &BroadcastResponse{}
	err = json.Unmarshal(responseBytes, broadcastResponse)
	if err != nil {
		return "", err
	}
	if !broadcastResponse.Success {
		return "", errors.New(fmt.Sprintf("[broadcast] the server return false, code=%d , msg=%s", broadcastResponse.Code, broadcastResponse.Msg))
	}
	return broadcastResponse.Data, nil
}

func (market *Market) getWebsocket() (*NdexWs, error) {
	if market.ndexWs == nil {
		market.ndexWs = &NdexWs{
			Host: market.WsHost,
		}
		err := market.ndexWs.Conn()
		if err != nil {
			return nil, err
		}
	}
	return market.ndexWs, nil
}

/**
 * Pending orders and changes to the configuration address
 * 订阅配置地址的挂单及变化
 */
func (market *Market) SubscribeOrderChange() (chan *WsOrderChange, error) {
	if market.Address == "" {
		return nil, errors.New("No address is configured")
	}
	return market.SubscribeOrderChangeByAddress(market.Address)
}

/**
 * Subscribe to pending orders and changes at specified addresses
 * 订阅指定地址的挂单及变化
 */
func (market *Market) SubscribeOrderChangeByAddress(address string) (chan *WsOrderChange, error) {
	ndexWs, err := market.getWebsocket()
	if err != nil {
		return nil, err
	}
	return ndexWs.SubscribeOrderChange(address)
}

/**
 * Subscription transaction changes to order book
 * 订阅交易对盘口及变化
 */
func (market *Market) SubscribeOrderBook(symbol string, top int) (chan *OrderBook, error) {
	ndexWs, err := market.getWebsocket()
	if err != nil {
		return nil, err
	}
	return ndexWs.SubscribeOrderBook(symbol, top)
}

/**
 * UnSubscription transaction changes to order book
 * 取消订阅交易对盘口及变化
 */
func (market *Market) UnSubscribeOrderBook(symbol string) (error) {
	ndexWs, err := market.getWebsocket()
	if err != nil {
		return err
	}
	return ndexWs.UnSubscribeOrderBook(symbol)
}

/**
 * UnSubscription Pending orders and changes to the configuration address
 * 取消订阅配置地址的挂单及变化
 */
func (market *Market) UnSubscribeOrderChange() (error) {
	if market.Address == "" {
		return errors.New("No address is configured")
	}
	return market.UnSubscribeOrderChangeByAddress(market.Address)
}

/**
 * UnSubscription to pending orders and changes at specified addresses
 * 取消订阅指定地址的挂单及变化
 */
func (market *Market) UnSubscribeOrderChangeByAddress(address string) (error) {
	ndexWs, err := market.getWebsocket()
	if err != nil {
		return err
	}
	return ndexWs.UnSubscribeOrderChange(address)
}


/**
 * Subscription configuration address balance changes
 * 订阅配置地址的余额变化
 */
func (market *Market) SubscribeBalanceChange() (chan *WsBalanceChange, error) {
	if market.Address == "" {
		return nil, errors.New("No address is configured")
	}
	return market.SubscribeBalanceChangeByAddress(market.Address)
}

/**
 * Subscribe to the balance change of the specified address
 * 订阅指定地址的余额变化
 */
func (market *Market) SubscribeBalanceChangeByAddress(address string) (chan *WsBalanceChange, error) {
	ndexWs, err := market.getWebsocket()
	if err != nil {
		return nil, err
	}
	return ndexWs.SubscribeBalanceChange(address)
}