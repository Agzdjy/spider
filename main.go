package main

import (
	"spider/douban"
	"spider/ziroom"
	"fmt"
)
func main()  {
	done := make(chan string, 2)

	go func() {
		douban.Run()
		done <- "douban"
	}()
	go func() {
		ziroom.Run()
		done <- "ziroom"
	}()

	fmt.Print(<- done, "--->", <- done)
}
