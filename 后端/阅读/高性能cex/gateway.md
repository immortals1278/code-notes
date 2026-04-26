## gateway
将环境变量中不同功能的本地url，解析成对应的url.URL结构体指针，生成对应的反向代理。
设置好gin引擎，与反向代理连接
看到srv

## gin
Gin 框架中用于将 HTTP 请求映射到处理函数的组件

`r := gin.New()`创建一个不带任何中间件的gin引擎

`r.Use()`手动添加中间件，**中间件**是一个可插拔的函数，在请求到达最终处理函数之前或之后执行。
- `r.Use(gin.Recovery())`  Gin 内置的异常恢复中间件，当代码有错，中间件会捕获它，返回 500 Internal Server Error
- `r.Use(metrics.Middleware("gateway"))`  监控指标中间件，记录网关运行数据
- `r.Use(cors.New(cors.Config{ 这里 }))`cors中间件:规定允许的http方法，请求头，域名...

`r.GET()`Gin 框架中注册 HTTP GET 路由的方法。当客户端发送 GET 请求到指定路径时，执行对应的处理函数。
- `r.GET("/metrics", gin.WrapH(metrics.Handler()))`访问/metrics时，输出指标数据(网关工作中的统计数据) ，函数wrapH()成gin认得的形式
- `r.GET("/health", func(c *gin.Context){c.JSON(http.StatusOK, gin.H{  }) })`用来证明网关健康
  - c：*gin.Context，当前请求的**上下文对象**，包含了请求/响应的所有信息
  - .JSON：Gin 的方法，向客户端返回 JSON 格式(键值对)的响应
  - http.StatusOK：常量，表示成功响应
  - gin.H{}：一个对象，里面可以放很多键值对，用来放很多信息

`r.Any()`注册一个路由响应所有http方法
- `r.Any("/ws", gin.WrapH(marketProxy))`所有发往 /ws 的请求转发给 marketProxy 反向代理。将标准 http.Handler wrapH() 成 Gin 的处理函数

`apiGroup := r.Group("/api/v1")`创建一个路由分组，所有以 /api/v1 开头的路由都属于这个分组
`public.Use(middleware.RateLimitMiddleware(publicLimiter))`该分组下所有路由都会被限流(比如每秒100次)
`orders.Use(middleware.IdempotencyMiddleware(idempStore, 24*time.Hour))` 幂等性中间件：确保同一个请求多次发送，只会被执行一次(同一个 Key 24 小时内只能请求 1 次,不区分用户)

`r.NoRoute(gin.WrapH(orderProxy))`当请求的路径没有任何匹配的路由时，转发给 orderProxy

## 反向代理
把客户端请求转发给后端服务器处理，然后把结果返回给客户端。

反向代理会生成一个客户端调用的url(域名)，调用后反向代理调用 http://localhost:8103 执行本地机器的后端代码（这部分后端代码运行在这个本地url上）

生成反向代理：
```golang
func newReverseProxy(targetURL *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
    //调用 Go 标准库函数创建基础反向代理
	originalDirector := proxy.Director
    //保存proxy.Director最初的逻辑：改写请求目标地址，发给本地url
	proxy.Director = func(req *http.Request) {
    //将匿名函数赋给 proxy.Director函数，覆盖原有逻辑。
		originalDirector(req)
    //在原有逻辑的基础上增加新逻辑
		req.Header.Set("X-Forwarded-Host", req.Host)
    //记录客户端原本访问的主机名:将字段X-Forwarded-Host的值设置为req.Host
    //主机名：url中域名的部分
		if req.TLS != nil {
			req.Header.Set("X-Forwarded-Proto", "https")
			return
		}
		req.Header.Set("X-Forwarded-Proto", "http")
    //记录客户端使用的协议
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Log.Error("gateway: 反向代理失敗", zap.String("path", r.URL.Path), zap.Error(err))
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error":   "upstream unavailable",
			"service": "gateway",
		})
	}
    //自定义错误处理逻辑
	return proxy
}
```
## 微服务架构的标准目录结构
每个文件夹下都是main.go。每个 main.go 都是一个独立微服务的可执行程序入口，用于启动该服务。

## golang
`defer logger.Sync()`
logger.Sync() 会刷新日志缓冲区，把还没写入文件的日志条目强制写入磁盘
defer 保证在函数返回（程序退出）时执行这行代码 

`port := os.Getenv("GATEWAY_PORT")`
os.Getenv() 是 Go 标准库函数，用于读取操作系统环境变量
"GATEWAY_PORT" 是环境变量名，需要在启动程序前设置好.在该项目中若没配置则用8100作为默认端口

`orderTargetURL, err := url.Parse(orderServiceURL)`url.Parse() 是 Go 标准库 net/url 包提供的函数，用于解析 URL 字符串.解析成功后返回 *url.URL 结构体指针

`redisCfg := infraredis.DefaultConfig()`
infraredis 是当前项目内部包的别名
DefaultConfig() 是该包提供的一个函数，返回一个预填充了默认值的配置结构体，这样只修改需要变更的字段，其他保持默认
import时声明了infraredis是哪个内部包的别名

`parsed, err := strconv.Atoi(val)`将字符串转化为整数类型

`time.Second`时间间隔常量，表示1秒

` { } `定义一个代码块， { } 里的变量外面访问不到