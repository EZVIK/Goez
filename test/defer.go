package main

import (
	"fmt"
	"time"
)

//
func main() {

	// 1、作用域 defer 会在当前函数和方法返回之前调用
	//{
	//	defer fmt.Println("3 defer runs")
	//
	//	fmt.Println("1 blocks ends")
	//}
	//i := -10
	//defer fmt.Println("10 main ends 2")
	//
	//
	//for i < 0 {
	//	i += 2
	//	time.Sleep(time.Millisecond * 300)
	//	//fmt.Println(i)
	//	defer fmt.Println("* for ends:", i)		// 在for 作用域中多次执行 defer 会从执行的最后一个defer 开始执行
	//}
	//
	//fmt.Println("2 main ends")


	// 2. 预计算参数
	startedAt := time.Now()
	defer fmt.Println(time.Since(startedAt))   // 150ns

	defer func() {
		fmt.Println(time.Since(startedAt))	   // 1.00..s
	}()

	time.Sleep(time.Second)

}
