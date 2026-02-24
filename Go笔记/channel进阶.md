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
## 阻塞
```golang
ch := make(chan int,3)

go func(){
    ch <- 1
    ch <- 2
    ch <- 3
    fmt.Println("send done")
    ch <- 4//这个协程会阻塞在这里（因为容量满了还往里加）
    fmt.Println("this log cant be seen")
}()

```
### 协程和线程
线程运行在操作系统（OS）中，而协程运行在线程中。

协程的数量远大于线程的数量

`GMP`:协程，线程，线程执行的环境。一个M绑定一个P，P维护一个全局列表，列表里面存储了很多runnable的协程。G运行在M里

1.谁来阻塞协程：运行时调度器将调`gopark()`，被阻塞的G放在channel里（hchan的等待队列recvq里）

2.过程：在P里挑一个runnable的G，调度器将它与M关系绑定开始执行

3.阻塞后唤醒：调度器执行`goready()`，把唤醒的G放回P里