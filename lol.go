package main

import (
	"fmt"
)

func main() {

	iter := []string{"lol", "lol2"}

	c := make(chan map[string]string, len(iter))

	for _, p := range iter {
		go func(p string) {
			c <- map[string]string{
				"status": "up",
				"i":      p,
			}
		}(p)
	}

	for range iter {
		r := <-c
		fmt.Println(r)
	}

	close(c)

}
