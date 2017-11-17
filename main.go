package main

import (
	"fmt"
	"reflect"
	"spider/douban"
	"spider/spider"
	"spider/v2ex"
	"spider/ziroom"
	"strings"
)

func main() {
	done := make(chan string, 2)
	go run(&douban.Douban{}, done)
	go run(&ziroom.Ziroom{}, done)
	go run(&v2ex.V2ex{}, done)

	count := 0
	for {
		fmt.Println(<-done)
		count += 1
		if count == 3 {
			fmt.Println("spider over")
			return
		}
	}
}

func run(spider spider.Spider, done chan string) {
	spider.Run()
	done <- "over--->" + strings.Split(reflect.TypeOf(spider).String(), ".")[1]
}
