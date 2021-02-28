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
	for i := 0; i < 10; i++ {
		wg.Add(2)
		id := rnd.Intn(10) + 1
		go func(id int, wg *sync.WaitGroup) {
			if b, ok := queryCache(id); ok {
				fmt.Println("from cache")
				fmt.Println(b)
			}
			wg.Done()
		}(id, wg)
		go func(id int, wg *sync.WaitGroup) {
			if b, ok := queryDatabase(id); ok {
				fmt.Println("from database")
				cache[id] = b
				fmt.Println(b)
			}
			wg.Done()
		}(id, wg)
	}
	wg.Wait()
}

func queryCache(id int) (book.Book, bool) {
	b, ok := cache[id]
	return b, ok
}

func queryDatabase(id int) (book.Book, bool) {
	time.Sleep(10 * time.Millisecond)
	for _, b := range book.AllBooks {
		if b.ID == id {
			cache[id] = *b
			return *b, true
		}
	}

	return book.Book{}, false
}
