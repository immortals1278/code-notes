# channel 进阶
## 数据结构
先进先出，环形数组
## 并发安全
多个goroutine同时访问同一个数据时不会出现数据竞争
### 锁机制
加锁来确保同一时刻只有一个人在修改：

goroutine：lock -> 修改数据逻辑 -> unlock
```golang
func TestConcurrencyWithLock(t *testing.T) {
    var shared []int
    var mu sync.Mutex
    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            shared = append(shared, 1)
            mu.Unlock()
        }()
    }
    //如果不lock的话，10个goroutine之间就会冲突
    wg.Wait()
    fmt.Println(len(shared))
}
```

