// Package rethink provides login
package rethink

import (
	"time"
)

const (
	ChannelPrivate = iota
	ChannelPublic
	ChannelRoom
)

// MessageBase message base
type MessageBase struct {
	Id        string `gorethink:"id,omitempty"`
	Sender    string `gorethink:"sender"`
	Content   string `gorethink:"content"`
	CreatedAt int64  `gorethink:"created_at"`
}

// PublicMessage public message
type PublicMessage struct {
	Id        string `gorethink:"id,omitempty"`
	Sender    string `gorethink:"sender"`
	Content   string `gorethink:"content"`
	CreatedAt int64  `gorethink:"created_at"`
	//MessageBase
}

// RoomMessage message
type RoomMessage struct {
	MessageBase
}

// PrivateMessage private message
type PrivateMessage struct {
	Receiver string `gorethink:"receiver"`
	MessageBase
}

// ChatRecord of list
type ChatRecord struct {
	Id        string
	Sender    string
	CreatedAt int64
	Channel   int
	Content   string
}

func (p *ChatRecord) TimeSince() string {
	return time.Unix(p.CreatedAt, 0).Format("01-02 15:04:05")
}
