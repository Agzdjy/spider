package main

import (
	"fmt"
	"spider/douban"
	"spider/ziroom"
)

func main() {
	done := make(chan string, 2)

	go func() {
		douban.Run()
		done <- "douban"
	}()
	go func() {
		ziroom.Run()
		done <- "ziroom"
	}()

	fmt.Print(<-done, "--->", <-done)
}
