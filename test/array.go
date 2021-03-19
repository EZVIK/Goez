package main

import (
	"fmt"
	"sync"
)

func main() {

	//randints := []int{1,2,3,4,5,6,7,8,9,10}
	//
	//rew := map[int]int{}
	//res := [5]int{}
	//
	//s1 := rand.NewSource(time.Now().UnixNano())
	//r1 := rand.New(s1)
	//for i:= 0; i< 5; i++ {
	//
	//	for {
	//		num := r1.Intn(len(randints))
	//		if rew[num] != 1 {
	//			rew[num] = 1
	//			res[i] = num
	//			break
	//		}
	//	}
	//}
	//fmt.Println(res)

	a := 1
	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			a += 1
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println(a)

}
