订单怎么进到book里？
# engineManager结构体

一个锁
一个map string -> engine
## 方法
GetEngine 取得指定交易對的引擎，不存在則自動建立。有锁保护map的并发安全
GetSymbols 返回所有已註冊的交易對
reset 删掉所有engine

# engine结构体

一个锁
一个orderbook

## 方法

增删查。都要上锁

增加新订单时执行撮合逻辑，撮合完剩余的限价单加入book

### 查
返回book的快照，输入depth指定返回的asks，bids里的数量（0表示全部）

### 撮合方法
返回trades切片

死循环
拿到“对面最好的“订单
检查userID不一样（避免左手换右手）。检查价格是否一样（仅限价单）
得到成交币数（谁少选谁）
写trade，放到trades里
更新两边order的币数
检查币数：如果对面数量为0则删除订单，如果自己数量为0则退出死循环

买卖撮合逻辑差异：检查价格时都是检查卖价是否大于买价（多的钱买方省下来了）

# orderbook结构体

symbol（一种币一个book）

放order结构体实例：
bids列表（降序）
asks列表（升序）

## 方法

增删查的基础功能

### 估算买入订单所需资金

输入要买的数量，遍历asks，直到找到的卖单提供的币数满足要买的数量，返回总价

如果遍历完asks币数都不够报错

# 更多结构体
限价单，市价单及其创建方法

trade

快照，快照里的orderBookLevel结构体
字段后的json用于转化为json后字段名是小写（json就是要小写）（不然直接用go字段名是大写）
```golang
type OrderBookSnapshot struct {
	Symbol       string           `json:"symbol"`
	Bids         []OrderBookLevel `json:"bids"` // 買單 (Price DESC)
	Asks         []OrderBookLevel `json:"asks"` // 賣單 (Price ASC)
	FencingToken int64            `json:"fencing_token"` // 防腦裂令牌，確保 Redis LWW (Last Write Wins) 語意
}
```
？？
