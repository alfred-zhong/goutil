package main

import (
	"github.com/alfred-zhong/goutil/prome"
)

func main() {
	client := prome.NewClient("test", "/foo")
	client.UpdateMemStatsInterval = 5
	if err := client.ListenAndServe(":9000"); err != nil {
		panic(err)
	}
}
