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

### 其他
`http.ListenAndServe(":8080", nil)`启动 HTTP 服务器，监听 **8080 端口**，使用默认路由（需预先注册）。

`"root:123456@tcp(127.0.0.1:3306)/exchangeengine"`指定 **TCP 连接**，数据库地址为本地 `127.0.0.1`，端口 `3306`（MySQL 默认）。

```golang
type Response struct {
	Code int         `json:"code"`//指定 JSON 序列化/反序列化时的字段映射：
	Data interface{} `json:"data,omitempty"` //如果 Data 为空值，JSON 中省略该字段
}
```

前端代码在浏览器运行，浏览器转发给后端
```golang
w.Header().Set("Content-Type", "application/json")               // 设置响应格式为JSON
w.Header().Set("Access-Control-Allow-Origin", "*")              // *表示允许所有域名跨域访问,所有网站都能调用该api
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // 允许的HTTP方法
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // 允许的请求头，前端发送 JSON 时必须带这个头，不设置会被浏览器拦截。
```

浏览器发送跨域请求前，会先发送一个 **OPTIONS 预检请求**（不带业务数据），询问后端是否允许实际请求。
后端收到 OPTIONS 后：
- `w.WriteHeader(http.StatusOK)` → 返回 200 状态码表示"允许"
- `return` → 立即结束，不处理后续业务逻辑

这样浏览器收到 200 后，才会继续发送真实的 GET/POST 请求。
```golang
if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
```
