# ndex-go-sdk

For docking with decentralized exchanges to realize intelligent trading function.



**Useage:**

Use the following code to initialize.

```
market = &Market{
   Host: "",
   WsHost: "",
   Address: "",
   PrivateKey: "",
}
market.Initialize()
```



The usage of websocket is as follows.

```
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
```



For more rest api and websocket usage, please refer to the code and test cases market_test.go and ndex_ws_test.go.