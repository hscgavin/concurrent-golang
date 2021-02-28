package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hscgavin/concurrent-golang/book"
)

var cache = map[int]book.Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	for i := 0; i < 10; i++ {
		id := rnd.Intn(10) + 1
		go func(id int) {
			if b, ok := queryCache(id); ok {
				fmt.Println("from cache")
				fmt.Println(b)
			}
		}(id)
		go func(id int) {
			if b, ok := queryDatabase(id); ok {
				fmt.Println("from database")
				cache[id] = b
				fmt.Println(b)
			}
		}(id)
		time.Sleep((150 * time.Millisecond))
	}
	time.Sleep(2 * time.Second)
}

func queryCache(id int) (book.Book, bool) {
	b, ok := cache[id]
	return b, ok
}

func queryDatabase(id int) (book.Book, bool) {
	time.Sleep(300 * time.Millisecond)
	for _, b := range book.AllBooks {
		if b.ID == id {
			return *b, true
		}
	}

	return book.Book{}, false
}
