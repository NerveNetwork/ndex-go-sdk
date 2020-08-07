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
package utils

import (
	json2 "encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func RequestGet(url string) ([]byte, error) {
	h := NewHttpSend(url)
	content, err := h.Get()
	return content, err
}

func RequestHttpGet(url string, params map[string]interface{}) ([]byte, error) {
	client := http.Client{}
	json,error := json2.Marshal(params)
	if error != nil {
		return nil, error
	}
	request, err := http.NewRequest("GET", url, strings.NewReader(string(json))) //请求
	if err != nil {
		return nil, err // handle error
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8") //设置Content-Type
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36") //设置User-Agent
	response, err := client.Do(request)                //返回
	if err != nil {
		return nil, err
	}
	if response == nil || response.Body == nil {
		return nil, errors.New("response is nil")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func RequestPost(url string, params map[string]interface{}) ([]byte, error) {
	client := http.Client{}
	json,error := json2.Marshal(params)
	if error != nil {
		return nil, error
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(string(json))) //请求
	if err != nil {
		return nil, err // handle error
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8") //设置Content-Type
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36") //设置User-Agent
	response, err := client.Do(request)                //返回
	if err != nil {
		return nil, err
	}
	if response == nil || response.Body == nil {
		return nil, errors.New("response is nil")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
