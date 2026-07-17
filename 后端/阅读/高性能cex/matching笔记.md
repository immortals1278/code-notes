## matching engine
初始化db池。db操作少，20个连接即可

初始化redis，放到 matching svc 里用

设置kafka config，初始化 kafka producer，放matching svc 里面

初始化 svc ：NewSubscriber（engineManage，producer，redis客户端）

设定leader elector，写成为/失去leader时的逻辑，启动elector

搞路由，优雅关闭服务


## db
在db包中import github.com/jackc/pgx/v5/pgxpool  Go 语言的第三开源免费方库来创建db连接池
## leader election
在多实例部署下，确保只有一个实例执行关键任务。防止脑裂（多个实例同时操作db导致冲突）
为什么要部署多实例？

### 初始化
要一个**instance ID**:记录哪个实例执行关键任务
`instanceID, err := os.Hostname()`获取当前OS（容器）的主机名。当前一个容器只有一个实例，所以实例id就是主机名。

需要一个pool（包在electionRepo里）

### 启动
后台开个goroutine启动elector

**成为/失去leader时的逻辑**：启动时需要
- 失去逻辑：获取互斥锁。`svc.SetFencingToken(0)`将 FencingToken 清零，让任何正在执行的撮合逻辑拒绝实际操作。`consumerCancel()`停止consumer接受新信息。`matchConsumer.Wait()`等待正在处理的消息完全处理完
- 成为逻辑：获取互斥锁。冷启动（具体的懒得看了）
## 冷启动
从完全空白的状态初始化的过程

与"热启动"相对，热启动是从缓存恢复，冷启动需要从持久化存储（DB）完整重建状态

## 注册路由
内部撮合引擎，不像orderService一样需要复杂的对外HTTP api服务，仅提供/metrics 和 /health 给运维调用（前端不会有调用这俩的东西所以用户调用不了）

## 优雅关闭


## go
go包可以有子包
