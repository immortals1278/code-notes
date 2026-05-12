## orderService
初始化db，redis，kafka，outbox

创建一个订单服务实例

启动kafka consumer(每个微服务都要有)

设置gin，为业务的核心功能（internal/order）注册路由，启动http server

优雅关闭

db ，repository ，domain ，outbox


## 数据库
数据库连接url写在环境变量中

初始化数据库池
```golang
dbCfg := db.DefaultDBConfig(dbURL) // 创建一个默认配置的 数据库配置 。

dbCfg.MaxOpenConns = 150 // 设置连接池最多同时打开 150 个数据库连接。

dbCfg.MaxIdleTime = 5 * time.Minute // 一个空闲连接最多保持 5 分钟，超时后会被关闭并回收。

pool, err := db.NewPostgresPool(context.Background(), dbCfg) // 创建连接池
if err != nil {
	logger.Log.Fatal("order-service: 無法連接資料庫", zap.Error(err))
}

defer pool.Close() // 在 main() 返回前关闭池
```

`context.Background()` 创建一个空的根 context，用于控制连接池的生命周期

**连接池**：预先创建的一批连接，反复使用，用完不销毁，放回池中

## repository
`repo := repository.NewPostgresRepository(pool)` 把 pool 实例放进去创建一个 repo 实例，不知道有什么用

## redis
初始化redis
```golang
redisCfg := redis.DefaultConfig() // 创建 redis 默认配置

if redisAddr := os.Getenv("REDIS_URL"); redisAddr != "" { // 从系统环境变量读取 redis_url ，放到 redis 默认配置中
	redisCfg.Addr = redisAddr
}

redisClient, err := redis.NewClient(redisCfg) // 根据 redis 默认配置创建 redis 客户端对象
if err != nil {
	logger.Log.Warn("Redis 不可用，市價單預估功能將受限", zap.Error(err)) // 用 Warn 而不是 Fatal，非核心功能即使连不上 redis 订单功能依然可用
}

var cacheRepo domain.CacheRepository // ?

if redisClient != nil {
	cacheRepo = repository.NewRedisCacheRepository(redisClient) // 创建 redis 的 repo
	logger.Log.Info("Redis 已連線")
}
```

## kafka
持久化的消息队列系统，用于微服务之间传递消息

订单事件必须通过他发布给其他服务

kafka的consumer用来接受消息，producer用来发消息（放在outbox里）

`topic` 消息主题
```golang
kafkaCfg := kafka.DefaultConfig() // 创建默认 kafka 配置

if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" { //从环境变量中读取信息给 kafka 配置
	kafkaCfg.Brokers = strings.Split(brokers, ",")
}
if resetOffset := os.Getenv("KAFKA_RESET_OFFSET"); resetOffset != "" {
	kafkaCfg.ResetOffset = strings.ToLower(resetOffset)
}

if os.Getenv("GIN_MODE") == "release" { //判断 gin 是否是生产模式，若是则禁止kafka自动创建topic
	kafkaCfg.AllowAutoTopicCreation = false
}

kafkaProducer, err := kafka.NewProducer(kafkaCfg)
if err != nil {
	logger.Log.Fatal("order-service: Kafka 連線失敗，純微服務模式無法啟動", zap.Error(err))
}
eventBus := domain.EventPublisher(kafkaProducer)
```‘
启动kafka消费者
```golang
eventSubscriber := order.NewEventSubscriber(repo, repo, repo, repo, kafkaProducer)

consumerCtx, cancelConsumers := context.WithCancel(context.Background())
defer cancelConsumers()

settleConsumer, err := kafka.NewConsumer(kafkaCfg, "order-service", []string{domain.TopicSettlements})
if err != nil {
	logger.Log.Fatal("order-service: 建立結算 consumer 失敗", zap.Error(err))
}
settleConsumer.Start(consumerCtx, eventSubscriber.HandleEvents)
logger.Log.Info("Kafka settlement consumer 已啟動", zap.String("topic", domain.TopicSettlements))```
```

## Outbox Pattern
满足数据库操作与发送消息需要原子性的需求：避免发送了信息没更新数据库，导致数据不一致

原理：发送事件也变成一次数据库操作，操作“发送表”。一个goroutine不断轮询“发送表”并执行对应操作。

`outboxCtx, cancelOutbox := context.WithCancel(context.Background())`
- `context.WithCancel()`：基于父 Context 创建一个可以手动取消的子 Context
- `outboxCtx`：返回的可取消 Context，用于传递给需要受控的 goroutine
- `cancelOutbox`：取消函数，调用后会通知所有使用 outboxCtx 的 goroutine 停止工作

## api和业务逻辑
为业务的核心功能注册路由
```golang
handler := api.NewHandler(svc, nil)
v1 := r.Group("/api/v1")
handler.RegisterRoutes(v1) // 设置各种功能的路由
```
所有发到当前微服务 /api/v1 路径的请求，都会通过 handler 注册的路由，最终调用 svc（order包的实例） 中的业务逻辑代码。

handler（api包）调用order包中的核心业务逻辑


**路由**
```golang
r.POST("/orders", h.CreateOrder)
//  └─路由规则─┘   └─处理器─┘
```

## 优雅关闭
优雅关闭kafka consumer
```golang
cancelConsumers()
cancelOutbox()
shutdownDone := make(chan struct{})
go func() {
	if settleConsumer != nil { //如果是nil就不用关闭
		settleConsumer.Wait() // 等待直到消费者完成所有消息处理
	}
	close(shutdownDone) // 关闭chan，接收方立即受到该chan的0值
}()
select {
case <-shutdownDone:
	logger.Info("Kafka consumers 已完整關閉")
case <-time.After(10 * time.Second): // 10s没关闭强制关闭
	logger.Warn("Kafka consumer 等待超時，強制繼續關機")
}
```
直接关闭kafka producer
优雅关闭http server
略


## golang
`defer logger.Sync()` 确保在程序退出前将所有日志写入磁盘

`strings.ToLower(resetOffset)`：转换为小写，避免大小写不一致

`url := ginSwagger.URL("/docs/architecture/swagger.yaml")`告诉gin-swagger去哪里加载api配置文件

`<-quit`阻塞等待，直到收到信号

