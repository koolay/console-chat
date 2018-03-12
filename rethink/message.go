// Package rethink provides login
package rethink

// MessageBase message base
type MessageBase struct {
	Id        string `gorethink:"id,omitempty"`
	CreatedAt int64  `gorethink:"created_at"`
}

// PublicMessage public message
type PublicMessage struct {
	Sender  string `gorethink:"sender"`
	Content string `gorethink:"content"`

	MessageBase
}

// RoomMessage message
type RoomMessage struct {
	Sender  string `gorethink:"sender"`
	Content string `gorethink:"content"`
	MessageBase
}

// PrivateMessage private message
type PrivateMessage struct {
	Sender   string `gorethink:"sender"`
	Receiver string `gorethink:"receiver"`
	Content  string `gorethink:"content"`
	MessageBase
}

type Users struct {
	Username  string `gorethink:"username"`
	Password  string `gorethink:"password"`
	CreatedAt int64  `gorethink:"created_at"`
}
