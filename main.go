package main

import (
	"spider/v2ex"
)

func main() {
	done := make(chan string, 1)

	// go func() {
	// 	douban.Run()
	// 	done <- "douban"
	// }()
	// go func() {
	// 	ziroom.Run()
	// 	done <- "ziroom"
	// }()

	go func() {
		v2ex.Run()
		done <- "v2ex"
	}()
	<-done
	// fmt.Print(<-done, "--->", <-done)
}
