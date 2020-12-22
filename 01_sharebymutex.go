package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var names = [9]string{"kim", "cho", "won", "a", "b", "c", "d", "e", "ffff"}
var queue []string
var results []string
var mutexR, mutexW sync.Mutex

func feedNames() {
	for _, name := range names {
		for {
			mutexR.Lock()
			if len(queue) < 2 {
				queue = append(queue, name)
				mutexR.Unlock()
				break
			}
			mutexR.Unlock()
			time.Sleep(1 * time.Second)
		}

		for {
			mutexW.Lock()
			if len(results) != 0 {
				fmt.Println(results[0])
				results = results[1:]
				mutexW.Unlock()
				break
			}
			mutexW.Unlock()
			time.Sleep(1 * time.Second)
		}
	}
}

func convert() {
	var name string
	for {
		for {
			mutexR.Lock()
			if len(queue) != 0 {
				name = strings.ToUpper(queue[0])
				queue = queue[1:]
				mutexR.Unlock()
				break
			}
			mutexR.Unlock()
		}
		for {
			mutexW.Lock()
			if len(results) < 2 {
				results = append(results, strings.ToUpper(name))
				mutexW.Unlock()
				break
			}
			mutexW.Unlock()
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	queue = make([]string, 0)
	results = make([]string, 0)
	go feedNames()
	go convert()
	time.Sleep(10 * time.Second)
}
