package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
The final value is not 1000 because counter++ is non-atomic operation
executed concurrently by several goroutines, which lead to data race and unstable output.
*/
//func main() {
//	var counter int
//	var wg sync.WaitGroup
//	var mu sync.Mutex
//
//	for i := 0; i < 1000; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			mu.Lock()
//			counter++
//			mu.Unlock()
//		}()
//	}
//
//	wg.Wait()
//	fmt.Println(counter)
//}

func main() {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Println(counter)
}
