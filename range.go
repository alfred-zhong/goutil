package goutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Range 表示区间。
type Range interface {
	LowerEndpoint() EndPoint
	UpperEndpoint() EndPoint

	IsValid() bool

	IsConnected(Range) bool
}

// BoundType 表示边界类型。
type BoundType int

const (
	// OPEN 表示边界开区间。
	OPEN BoundType = 0
	// CLOSED 表示边界闭区间。
	CLOSED BoundType = 1
)

// IsOpen 返回该边界类型是否为开区间。
func (b BoundType) IsOpen() bool {
	return b == OPEN
}

// IsClosed 返回该边界类型是否为闭区间。
func (b BoundType) IsClosed() bool {
	return b == CLOSED
}

// EndPoint 表示边界点。
type EndPoint struct {
	BoundType
	Value float64
}

// CompareTo 返回该边界点与参数边界点 e2 的比较值。
func (e1 EndPoint) CompareTo(e2 EndPoint) int {
	if e1.Value < e2.Value {
		return -1
	} else if e1.Value > e2.Value {
		return 1
	} else {
		return 0
	}
}

// newOpenEndPoint 新建一个边界类型为开的边界点。
func newOpenEndPoint(value float64) EndPoint {
	return EndPoint{
		BoundType: OPEN,
		Value:     value,
	}
}

// newClosedEndPoint 新建一个边界类型为闭的边界点。
func newClosedEndPoint(value float64) EndPoint {
	return EndPoint{
		BoundType: CLOSED,
		Value:     value,
	}
}

// NumberRange 表示数字化的区间.
type NumberRange struct {
	lower, upper EndPoint
}

func (n NumberRange) String() string {
	var sb strings.Builder
	if n.lower.IsOpen() {
		sb.WriteByte('(')
	} else if n.lower.IsClosed() {
		sb.WriteByte('[')
	} else {
		sb.WriteByte('?')
	}

	sb.WriteString(fmt.Sprintf(
		"%s,%s",
		strconv.FormatFloat(n.lower.Value, 'f', -1, 64),
		strconv.FormatFloat(n.upper.Value, 'f', -1, 64),
	))

	if n.upper.IsOpen() {
		sb.WriteByte(')')
	} else if n.upper.IsClosed() {
		sb.WriteByte(']')
	} else {
		sb.WriteByte('?')
	}

	return sb.String()
}

// LowerEndpoint 返回下限边界点。
func (n NumberRange) LowerEndpoint() EndPoint {
	return n.lower
}

// UpperEndpoint 返回上限边界点。
func (n NumberRange) UpperEndpoint() EndPoint {
	return n.upper
}

// IsValid 返回该 Range 是否有效合法。
func (n NumberRange) IsValid() bool {
	return n.lower.CompareTo(n.upper) <= 0
}

// IsConnected 返回该 Range 与其他 Range 是否有交集。
func (n NumberRange) IsConnected(r Range) bool {
	return n.LowerEndpoint().CompareTo(r.UpperEndpoint()) <= 0 &&
		r.LowerEndpoint().CompareTo(n.UpperEndpoint()) <= 0
}

// NewClosedRange 新建一个闭区间。
func NewClosedRange(lower, upper float64) Range {
	return NumberRange{
		lower: newClosedEndPoint(lower),
		upper: newClosedEndPoint(upper),
	}
}

// NewOpenRange 新建一个开区间。
func NewOpenRange(lower, upper float64) Range {
	return NumberRange{
		lower: newOpenEndPoint(lower),
		upper: newOpenEndPoint(upper),
	}
}

// NewClosedOpenRange 新建一个左闭右开区间。
func NewClosedOpenRange(lower, upper float64) Range {
	return NumberRange{
		lower: newClosedEndPoint(lower),
		upper: newOpenEndPoint(upper),
	}
}

// NewOpenClosedRange 新建一个左开右闭区间。
func NewOpenClosedRange(lower, upper float64) Range {
	return NumberRange{
		lower: newOpenEndPoint(lower),
		upper: newClosedEndPoint(upper),
	}
}
