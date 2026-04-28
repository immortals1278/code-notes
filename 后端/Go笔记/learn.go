package main

import (
	"fmt"
	"time"

)

func main(){
	ch1 := make(chan string)
	ch2 := make(chan string)//定义两个协程通信的channel

	go func(){
		time.Sleep(2*time.Second)//模拟真实业务执行耗时
		ch1 <- "message1"
	}()

	go func(){
		time.Sleep(1*time.Second)
		ch2 <- "message2"
	}()

	for{
		select{
		case message1 := <-ch1:
			fmt.Println("receive message",message1)
		case message2 := <-ch2:
			fmt.Println("receive message",message2)
		}//监听谁有值谁执行，谁先有信号谁执行
	}//主协程定义一个死循环
}