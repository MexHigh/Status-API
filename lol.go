package main

import (
	"fmt"
	"sync"
)

func main() {

	iter := []string{"lol", "lol2"}

	c := make(chan map[string]string, len(iter))

	var wg sync.WaitGroup

	for _, p := range iter {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			c <- map[string]string{
				"status": "up",
				"i":      p,
			}
		}(p)
	}
	wg.Wait()
	close(c)

	for i := range c {
		fmt.Println(i)
	}

}
