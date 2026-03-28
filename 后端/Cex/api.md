### http用法
```golang
func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")//设置响应头，告诉客户端返回内容为JSON格式
	var order model.Order

	json.NewDecoder(r.Body).Decode(&order)//从 HTTP 请求的 Body 中读取 JSON 数据，并将其解析到 order 变量中。(放入用户输入的订单数据)

	h.Engine.PlaceOrder(&order)//调用PlaceOrder()，这里缺少报错处理

	json.NewEncoder(w).Encode(Response{
		Code: 0,
		Msg:  "ok",
	})//向客户端返回一个http响应，包含状态码0和消息ok
}

id := r.URL.Query().Get("id")//从URL参数中获取id值

```
- `w http.ResponseWriter`：用于**发送 HTTP 响应**的接口（写回数据、设置状态码等）。
- `r *http.Request`：代表客户端发送的 **HTTP 请求**对象（包含请求方法、URL、Body 等）。

### main函数里http
- `http.HandleFunc("/order", handler.PlaceOrder)`：注册路由。当有请求访问 `/order` 路径时，调用 `handler.PlaceOrder` 函数来处理。
- `http.ListenAndServe(":8080", nil)`：启动 HTTP 服务，监听本机 `8080` 端口。`nil` 表示使用默认的路由器（即上面注册的路由）。

### Response结构体
```golang
Code int         `json:"code"`  //状态码，0表示成功，非0表示失败，根据具体值做不同处理
Msg  string      `json:"msg"`   //错误信息
Data interface{} `json:"data,omitempty"`
```