# goutil

![Go](https://github.com/alfred-zhong/goutil/workflows/Go/badge.svg?branch=master) [![GoDoc](https://godoc.org/github.com/alfred-zhong/goutil?status.svg)](https://godoc.org/github.com/alfred-zhong/goutil) [![Go Report Card](https://goreportcard.com/badge/github.com/alfred-zhong/goutil)](https://goreportcard.com/report/github.com/alfred-zhong/goutil)

## 区间

Range。类似 Guava 中的 Range，属于简化版。

## 即时 Ticker

InstantTicker。与 time.Ticker 类似，创建后会立马发送当前时间到 channel 中，接下来的行为类似于 time.Ticker。
