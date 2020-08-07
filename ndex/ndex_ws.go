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
 * @Date: 2020/5/20 下午12:08
 */
package ndex

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type NdexWs struct {
	Host 				string
	readChannel 		chan string
	writeChannel 		chan string
	done 				chan struct{}
	conn 				*websocket.Conn

	subscribeMap 		map[string]*WsSubInfo
}

func (ws *NdexWs) Ping() {
	msg := fmt.Sprintf("{\"ping\":%d}", time.Now().UnixNano() / 1e6)
	ws.writeChannel <- msg
}

func (ws *NdexWs) SubscribeOrderBook(symbol string, top int) (chan *OrderBook, error) {
	msg := fmt.Sprintf("{\"action\":\"Subscribe\",\"channel\":\"apiOrderBook:{\\\"symbol\\\":\\\"%s\\\",\\\"top\\\":%d}\"}", symbol, top)
	ws.writeChannel <- msg

	channel := fmt.Sprintf("apiOrderBook:%s", symbol)
	subInfo := ws.subscribeMap[channel]
	if subInfo == nil {
		orderBookEvent := make(chan *OrderBook, 30)
		subInfo = &WsSubInfo{
			Channel: 	channel,
			SubMessage: msg,
			Event: 		orderBookEvent,
		}
		ws.subscribeMap[channel] = subInfo
	}
	return subInfo.Event.(chan *OrderBook), nil
}

func (ws *NdexWs) SubscribeOrderChange(address string) (chan *WsOrderChange, error) {
	channel := fmt.Sprintf("order:%s", address)
	msg := fmt.Sprintf("{\"action\":\"Subscribe\",\"channel\":\"%s\"}", channel)
	ws.writeChannel <- msg

	subInfo := ws.subscribeMap[channel]
	if subInfo == nil {
		orderChangeEvent := make(chan *WsOrderChange, 10)
		subInfo = &WsSubInfo{
			Channel: 	channel,
			SubMessage: msg,
			Event: 		orderChangeEvent,
		}
		ws.subscribeMap[channel] = subInfo
	}
	return subInfo.Event.(chan *WsOrderChange), nil
}

func (ws *NdexWs) SubscribeBalanceChange(address string) (chan *WsBalanceChange, error) {
	channel := fmt.Sprintf("account:%s", address)
	msg := fmt.Sprintf("{\"action\":\"Subscribe\",\"channel\":\"%s\"}", channel)
	ws.writeChannel <- msg

	subInfo := ws.subscribeMap[channel]
	if subInfo == nil {
		balanceChangeEvent := make(chan *WsBalanceChange, 10)
		subInfo = &WsSubInfo{
			Channel: 	channel,
			SubMessage: msg,
			Event: 		balanceChangeEvent,
		}
		ws.subscribeMap[channel] = subInfo
	}
	return subInfo.Event.(chan *WsBalanceChange), nil
}

func (ws *NdexWs) UnSubscribeOrderBook(symbol string) error {
	msg := fmt.Sprintf("{\"action\":\"Unsubscribe\",\"channel\":\"apiOrderBook:{\\\"symbol\\\":\\\"%s\\\"}\"}", symbol)
	ws.writeChannel <- msg
	channel := fmt.Sprintf("apiOrderBook:%s", symbol)
	wsSubInfo := ws.subscribeMap[channel]
	close(wsSubInfo.Event.(chan *OrderBook))
	delete(ws.subscribeMap, channel)
	return nil
}

func (ws *NdexWs) UnSubscribeOrderChange(address string) error {
	msg := fmt.Sprintf("{\"action\":\"Unsubscribe\",\"channel\":\"order:%s\"}", address)
	ws.writeChannel <- msg
	channel := fmt.Sprintf("order:%s", address)
	wsSubInfo := ws.subscribeMap[channel]
	close(wsSubInfo.Event.(chan *WsOrderChange))
	delete(ws.subscribeMap, channel)
	return nil
}

func (ws *NdexWs) UnSubscribeBalanceChange(address string) error {
	msg := fmt.Sprintf("{\"action\":\"Unsubscribe\",\"channel\":\"account:%s\"}", address)
	ws.writeChannel <- msg
	channel := fmt.Sprintf("account:%s", address)
	wsSubInfo := ws.subscribeMap[channel]
	close(wsSubInfo.Event.(chan *WsBalanceChange))
	delete(ws.subscribeMap, channel)
	return nil
}

func (ws *NdexWs) reSubscribe() error {
	for _, wsSubInfo := range ws.subscribeMap {
		ws.writeChannel <- wsSubInfo.SubMessage
	}
	return nil
}

func (ws *NdexWs) ReConn() error {
	log.Println("ndex reConning...")
	err := ws.Conn()
	if err != nil {
		return err
	}
	err = ws.reSubscribe()
	if err == nil {
		log.Println("ndex reConn success.")
	} else {
		log.Println("ndex reConn error : ", err)
	}
	return err
}

func (ws *NdexWs) Conn() error {
	url := fmt.Sprintf(ws.Host + "/ws")
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	if ws.subscribeMap== nil {
		ws.subscribeMap = make(map[string]*WsSubInfo)
	}
	ws.readChannel = make(chan string, 100)
	ws.writeChannel = make(chan string, 10)
	ws.done = make(chan struct{})
	ws.conn = c

	go ws.readHandler()
	go ws.writeHandler()
	go ws.pingHandler()
	go ws.exitHandler()
	go ws.messageHandler()

	return nil
}

func (ws *NdexWs) messageHandler() {
	for {
		select {
		case <- ws.done:
			return
		case message := <- ws.readChannel:
			//log.Println("received message :  " + message)
			wsResponse := &WsResponse{}
			messageBytes := []byte(message)
			err := json.Unmarshal(messageBytes, wsResponse)
			if err != nil || wsResponse.Channel == "" {
				wsPong := &WsPong{}
				err = json.Unmarshal(messageBytes, wsPong)
				if err == nil {
					// processPong
					ws.processPong(wsPong)
				}
				break
			}
			if wsResponse.Status == 200 {
				if wsResponse.Action != "Data" {
					log.Println("[WARING] receive unknown message : ", message)
					break
				}
				switch wsResponse.Channel {
				case "apiOrderBook":
					orderBookResponse := &WsOrderBookResponse{}
					err = json.Unmarshal(messageBytes, orderBookResponse)
					if err != nil {
						log.Println(err, ", [apiOrderBook] message : ", message)
						break
					}
					symbol := orderBookResponse.Data.Symbol
					wsSubInfo := ws.subscribeMap["apiOrderBook:" + symbol]
					if wsSubInfo != nil {
						wsSubInfo.Event.(chan *OrderBook) <- orderBookResponse.Data
					}
				case "order":
					orderChangeResponse := &WsOrderChangeResponse{}
					err = json.Unmarshal(messageBytes, orderChangeResponse)
					if err != nil {
						log.Println(err, ", [order] message : ", message)
						break
					}
					if orderChangeResponse.Data.T == "update" {
						for _, order := range orderChangeResponse.Data.D {
							channel := fmt.Sprintf("order:%s", order.Address)
							wsSubInfo := ws.subscribeMap[channel]
							if wsSubInfo == nil {
								continue
							}
							orderList := []*Order{ order  }
							orderChange := &WsOrderChange{
								T:	orderChangeResponse.Data.T,
								D:	orderList,
							}
							wsSubInfo.Event.(chan *WsOrderChange) <- orderChange
						}
					} else if orderChangeResponse.Data.T == "init" {
						if len(orderChangeResponse.Data.D) == 0 {
							break
						}
						channel := fmt.Sprintf("order:%s", orderChangeResponse.Data.D[0].Address)
						wsSubInfo := ws.subscribeMap[channel]
						if wsSubInfo == nil {
							continue
						}
						wsSubInfo.Event.(chan *WsOrderChange) <- orderChangeResponse.Data
					}
				case "account":
					balanceChangeResponse := &WsBalanceChangeResponse{}
					err = json.Unmarshal(messageBytes, balanceChangeResponse)
					if err != nil {
						log.Println(err, ", [balance] message : ", message)
						break
					}
					channel := fmt.Sprintf("account:%s", balanceChangeResponse.Data.A)
					wsSubInfo := ws.subscribeMap[channel]
					if wsSubInfo == nil {
						continue
					}
					wsSubInfo.Event.(chan *WsBalanceChange) <- balanceChangeResponse.Data
				default:
					log.Println("[NOTICE] Not yet parsed message : ", message)
				}
			} else {
				log.Println("[ERROR] receive exception data : ", message)
			}
		}
	}
}

func (ws *NdexWs) readHandler() {
	for {
		select {
		case <- ws.done:
			log.Println("ndex websocket closing reader.")
			return
		default:
			_, message, err := ws.conn.ReadMessage()
			if err != nil {
				log.Println("ndex websocket close, " + err.Error())
				close(ws.done)
				for {
					err = ws.ReConn()
					if err == nil {
						break
					}
					time.Sleep(5*time.Second)
				}
				return
			}
			//log.Println("received message:  " + string(message))
			ws.readChannel <- string(message)
		}
	}
}

func (ws *NdexWs) writeHandler() {
	for {
		select {
		case <- ws.done:
			return
		case msg := <- ws.writeChannel:
			if msg != "" {
				err := ws.conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					//log.Println("write message error, now close the connection.", err)
					//close(ws.done)
					return
				}
			}
		}
	}
}

func (ws *NdexWs) pingHandler() {
	ws.Ping()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ws.Ping()
		case <- ws.done:
			return
		}
	}
}

func (ws *NdexWs) exitHandler() {
	defer ws.conn.Close()
	defer close(ws.readChannel)
	defer close(ws.writeChannel)
	for {
		select {
		case <- ws.done:
			log.Println("ndex websocket closed.")
			return
		}
	}
}

func (ws *NdexWs) processPong(pong *WsPong) {
	// Do nothing
}
