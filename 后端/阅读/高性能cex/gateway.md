

## golang
`defer logger.Sync()`
logger.Sync() 会刷新日志缓冲区，把还没写入文件的日志条目强制写入磁盘
defer 保证在函数返回（程序退出）时执行这行代码 

`port := os.Getenv("GATEWAY_PORT")`
os.Getenv() 是 Go 标准库函数，用于读取操作系统环境变量
"GATEWAY_PORT" 是环境变量名，需要在启动程序前设置好.在该项目中若没配置则用8100作为默认端口

`orderTargetURL, err := url.Parse(orderServiceURL)`url.Parse() 是 Go 标准库 net/url 包提供的函数，用于解析 URL 字符串.解析成功后的 *url.URL 结构体指针

`redisCfg := infraredis.DefaultConfig()`
infraredis 是当前项目内部包的别名
DefaultConfig() 是该包提供的一个函数，返回一个预填充了默认值的配置结构体，这样只修改需要变更的字段，其他保持默认
import时声明了infraredis是哪个内部包的别名

`parsed, err := strconv.Atoi(val)`将字符串转化为整数类型

`time.Second`时间间隔常量，表示1秒

100行，前面newReverseProxy不知道是啥