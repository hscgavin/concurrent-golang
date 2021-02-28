package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/hscgavin/concurrent-golang/book"
)

var cache = map[int]book.Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	wg := &sync.WaitGroup{}
	m := &sync.RWMutex{}
	for i := 0; i < 10; i++ {
		wg.Add(2)
		id := rnd.Intn(10) + 1
		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex) {
			if b, ok := queryCache(id, m); ok {
				fmt.Println("from cache")
				fmt.Println(b)
			}
			wg.Done()
		}(id, wg, m)
		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex) {
			if b, ok := queryDatabase(id, m); ok {
				fmt.Println("from database")
				fmt.Println(b)
			}
			wg.Done()
		}(id, wg, m)
	}
	wg.Wait()
}

func queryCache(id int, m *sync.RWMutex) (book.Book, bool) {
	m.RLock()
	b, ok := cache[id]
	m.RUnlock()
	return b, ok
}

func queryDatabase(id int, m *sync.RWMutex) (book.Book, bool) {
	time.Sleep(10 * time.Millisecond)
	for _, b := range book.AllBooks {
		if b.ID == id {
			m.Lock()
			cache[id] = *b
			m.Unlock()
			return *b, true
		}
	}

	return book.Book{}, false
}
