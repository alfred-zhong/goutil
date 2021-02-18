package main

import (
	"github.com/alfred-zhong/goutil/prome"
)

func main() {
	client := prome.NewClientWithOption(
		"test", "/foo",
		prome.WithRuntimeEnable(true, 120),
		prome.WithConstLables("env", "test", "foo", "bar"),
	)
	if err := client.ListenAndServe(":9000"); err != nil {
		panic(err)
	}
}
