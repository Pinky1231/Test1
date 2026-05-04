package main

import (
	"fmt"
	"sync"
	"time"
)

type Cash struct {
	data map[string]interface{}
	mx   sync.RWMutex
	wg   sync.WaitGroup
}

func main() {
	Csh := Cash{
		data: make(map[string]interface{}),
	}

	// Csh.Set("Rok", "Huynya", time.Second*6)
	// Csh.Set("Bob", 21, time.Second*3)
	// Csh.Set("Ric_astley", 67, time.Minute*15)
	// Csh.Set("John", 14, time.Minute*1)

	// fmt.Println(Csh.data)
	// fmt.Println(Csh.Get("Rick_astley"))
	// <-time.After(time.Second * 7)
	// fmt.Println(Csh.data)

	Csh.wg.Wait()

}

func (cash *Cash) Set(key string, value interface{}, ttl time.Duration) {
	if _, ok := cash.Get(key); !ok {
		cash.mx.Lock()
		defer cash.mx.Unlock()
		cash.data[key] = value
		cash.wg.Add(1)
		go func(ttl time.Duration) {
			defer cash.wg.Done()

			<-time.After(ttl)

			cash.mx.Lock()
			defer cash.mx.Unlock()
			fmt.Println(key, "был удален спустя", ttl)
			delete(cash.data, key)

		}(ttl)
	}

}

func (cash *Cash) Get(key string) (interface{}, bool) {
	cash.mx.RLock()
	defer cash.mx.RUnlock()

	result, ok := cash.data[key]
	if !ok {
		return nil, false
	} else {
		return result, true
	}

}
