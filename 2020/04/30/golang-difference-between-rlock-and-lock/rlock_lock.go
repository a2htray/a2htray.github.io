package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	m := map[string]string{
		"learn": "to the end",
	}
	lock := sync.RWMutex{}

	go func(name string) {
		lock.RLock()
		fmt.Printf("goroutine %s read: %s\n", name, m["learn"])
		time.Sleep(time.Second * 2)
		fmt.Printf("goroutine %s read: %s again\n", name, m["learn"])
		lock.RUnlock()
	}("A")

	go func(name string) {
		lock.Lock()
		fmt.Printf("goroutine %s read: %s\n", name, m["learn"])
		time.Sleep(time.Second * 2)
		lock.Unlock()
	}("B")

	<- time.After(time.Second * 5)
}