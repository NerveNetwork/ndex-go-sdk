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
 * @Date: 2020/5/20 下午12:46
 */
package ndex

import (
	"log"
	"testing"
)

var (
	ndexWs *NdexWs
	waitChan chan string
)

func init() {
	host := "ws://beta.nervedex.com"
	ndexWs = &NdexWs{Host: host}
	ndexWs.Conn()

	waitChan = make(chan string)
}

func TestNdexWs_Conn(t *testing.T) {
	msg := "{\"ping\":1589271332077}"
	ndexWs.writeChannel <- msg
	t.Log("success")
	<- waitChan
}

func TestNdexWs_Ping(t *testing.T) {
	ndexWs.Ping()
	t.Log("send success")
	<- waitChan
}

func TestNdexWs_SubscribeOrderBook(t *testing.T) {
	orderBookEvent, err := ndexWs.SubscribeOrderBook("NVTNULS", 10)
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

func TestNdexWs_SubscribeOrderChange(t *testing.T) {
	orderChangeEvent, err :=  ndexWs.SubscribeOrderChange("TNVTdN9iCdXS46SuN8UPd2kkdXmGNYgMMxuSJ")
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