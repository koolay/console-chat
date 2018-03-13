// Package rethink provides login
package rethink

// MessageBase message base
type MessageBase struct {
	Id        string `gorethink:"id,omitempty"`
	Sender    string `gorethink:"sender"`
	Content   string `gorethink:"content"`
	CreatedAt int64  `gorethink:"created_at"`
}

// PublicMessage public message
type PublicMessage struct {
	MessageBase
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
