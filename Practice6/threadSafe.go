package main

import (
	"fmt"
	"sync"
)

//func main() {
//	var safeMap sync.Map
//	var wg sync.WaitGroup
//
//	for i := 0; i < 100; i++ {
//		wg.Add(1)
//		go func(key int) {
//			defer wg.Done()
//			safeMap.Store("key", key)
//		}(i)
//	}
//
//	wg.Wait()
//
//	value, _ := safeMap.Load("key")
//	fmt.Printf("Value: %d\n", value)
//}

func main() {
	safeMap := make(map[string]int)
	var mu sync.RWMutex
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()

			mu.Lock()
			safeMap["key"] = key
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	mu.RLock()
	value := safeMap["key"]
	mu.RUnlock()

	fmt.Printf("Value: %d\n", value)
}
