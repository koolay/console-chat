// Package rethink provides subscribe
package rethink

import (
	"sync"
)

var (
	ChatRecordPoll = &sync.Pool{
		New: func() interface{} {
			return &ChatRecord{}
		},
	}

	ChatRecordChan = make(chan *ChatRecord)
)
