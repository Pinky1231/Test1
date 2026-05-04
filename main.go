package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Item struct {
	value interface{}
	ttl   time.Time
}

type Cash struct {
	data map[string]Item
	mx   sync.RWMutex
	wg   sync.WaitGroup
}

func main() {
	Csh := Cash{
		data: make(map[string]Item),
	}
	deadline := time.Now().Add(time.Second * 16)
	context, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	Csh.wg.Add(1)
	go Csh.cleaner(context)

	Csh.Set("Rok", "Huynya", time.Second*6)
	Csh.Set("Bob", 21, time.Second*3)
	Csh.Set("Ric_astley", 67, time.Second*15)
	Csh.Set("John", 14, time.Minute*1)

	fmt.Println(Csh.data)
	<-time.After(time.Second * 7)
	fmt.Println(Csh.Get("Bob"))
	fmt.Println(Csh.data)

	Csh.wg.Wait()

}

func (cash *Cash) Set(key string, value interface{}, ttl time.Duration) {
	cash.mx.Lock()
	defer cash.mx.Unlock()

	cash.data[key] = Item{
		value,
		time.Now().Add(ttl).Local(),
	}

}

func (cash *Cash) Get(key string) (interface{}, bool) {
	cash.mx.RLock()
	defer cash.mx.RUnlock()
	val, ok := cash.data[key]

	if !ok { //if key does'nt exits
		return nil, false
	} else if time.Now().UnixNano() > val.ttl.UnixNano() { // if time is expired
		return nil, false
	}
	return val.value, ok

}

func (cash *Cash) cleaner(ctx context.Context) {
	defer cash.wg.Done()
	ticker := time.NewTicker(time.Second)

	for {
		select {

		case <-ctx.Done():
			ticker.Stop()
			fmt.Println("Монитроинг завершен")
			return

		case <-ticker.C:
			cash.mx.Lock()

			for key, obj := range cash.data {
				if time.Now().UnixNano() > obj.ttl.UnixNano() {

					fmt.Println(key, "был удален")
					delete(cash.data, key)

				}

			}
			cash.mx.Unlock()
		}
	}

}
