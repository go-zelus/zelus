package redix

import (
	"fmt"
	"time"
)

const (
	EXPIRE = "expire" // 过期时间
	NX     = "nx"     // 如果不存在，则SET
	XX     = "xx"     // 如果存在，则SET
)

type empty struct {
}

type Operation struct {
	Name  string
	Value interface{}
}

type Operations []*Operation

func (o Operations) Find(name string) *Result {
	for _, attr := range o {
		if attr.Name == name {
			return NewResult(attr.Value, nil)
		}
	}
	return NewResult(nil, fmt.Errorf("operation found error: %s", name))
}

// WithExpire 超时
func WithExpire(t int) *Operation {
	return &Operation{
		Name:  EXPIRE,
		Value: time.Duration(t) * time.Second,
	}
}

// WithNX 不存在写入
func WithNX() *Operation {
	return &Operation{
		Name:  NX,
		Value: empty{},
	}
}

// WithXX 存在写入
func WithXX() *Operation {
	return &Operation{
		Name:  XX,
		Value: empty{},
	}
}
