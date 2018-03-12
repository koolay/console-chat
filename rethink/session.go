// Package main provides ...
package rethink

import (
	"errors"
	"fmt"
	"time"

	r "gopkg.in/gorethink/gorethink.v4"
)

// Options options
type Options struct {
	Address  string
	Database string
	Username string
	Password string
}

type ctx struct {
	username string
	inRoom   string
}

// Rethink driver
type Rethink struct {
	address  string
	database string
	username string
	password string
	ctx      *ctx
}

// NewRethink instance driver
func NewRethink(options *Options) *Rethink {
	rethink := &Rethink{
		address:  options.Address,
		username: options.Username,
		password: options.Password,
		database: options.Database,
		ctx:      &ctx{inRoom: "general", username: ""},
	}
	return rethink
}

// getPrivateMsgTable get table name of personal
func (p *Rethink) getPrivateMsgTable() string {
	return "private_message"
}

// getPublicMsgTable get public message table
func (p *Rethink) getPublicMsgTable() string {
	return "public_message"
}

func (p *Rethink) connect() (sess *r.Session, err error) {
	sess, err = r.Connect(r.ConnectOpts{
		Address:  p.address,
		Database: p.database,
		Username: p.username,
		Password: p.password,
	})
	return sess, err
}

func (p *Rethink) Login(username string, password string) error {
	sess, err := p.connect()
	result, err := r.Table("users").Filter(r.Row.Field("username").Eq(username)).Run(sess)
	if err != nil {
		return err
	}

	var user Users
	_, err = result.Peek(&user)
	if err != nil {
		return err
	}
	if !checkPasswordHash(password, user.Password) {
		return errors.New("username or password not match")
	}
	p.ctx.username = username
	return nil
}

func (p *Rethink) Join(username string, password string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	result, err := r.Table("users").Filter(r.Row.Field("username").Eq(username)).Run(sess)
	if err != nil {
		return err
	}
	var user Users
	hasUser, err := result.Peek(&user)
	if err != nil {
		return err
	}
	if hasUser {
		return errors.New(fmt.Sprintf("username %s already exist.", username))
	}
	hPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	user = Users{
		Username:  username,
		Password:  hPassword,
		CreatedAt: time.Now().Unix(),
	}
	_, err = r.Table("users").Insert(user).RunWrite(sess)
	return err

}

func (p *Rethink) SwitchRoom(room string) error {
	if err := checkRoomName(room); err != nil {
		return err
	}
	p.ctx.inRoom = room
	return nil
}

func (p *Rethink) CreateRoom(room string) error {

	err := checkRoomName(room)
	if err != nil {
		return err
	}

	sess, err := p.connect()
	if err != nil {
		return err
	}
	result, err := r.DB(p.database).TableCreate(room).RunWrite(sess)
	if err != nil {
		return err
	}
	fmt.Println("*** Create table result: ***")
	fmt.Println(result)
	return nil
}

func (p *Rethink) SendPublicMessage(message string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	msg := PublicMessage{
		Content: message,
		Sender:  p.ctx.username,
	}
	msg.CreatedAt = time.Now().Unix()
	_, err = r.Table(p.getPublicMsgTable()).Insert(msg).RunWrite(sess)
	return err
}

func (p *Rethink) SendRoomMessage(message string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	msg := RoomMessage{
		Content: message,
		Sender:  p.ctx.username,
	}
	msg.CreatedAt = time.Now().Unix()
	_, err = r.Table(p.ctx.inRoom).Insert(msg).RunWrite(sess)
	return err
}

func (p *Rethink) SendPrivateMessage(toUser string, message string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	msg := PrivateMessage{
		Sender:   p.ctx.username,
		Receiver: toUser,
		Content:  message,
	}
	msg.CreatedAt = time.Now().Unix()
	_, err = r.Table(p.getPrivateMsgTable()).Insert(msg).RunWrite(sess)
	return err
}

func (p *Rethink) FeedsPublic() error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	cursor, err := r.Table(p.getPublicMsgTable()).Changes().Run(sess)
	if err != nil {
		return err
	}
	var msg PublicMessage
	for cursor.Next(&msg) {
		fmt.Println(msg)
	}
	return nil
}

func (p *Rethink) FeedsRoom(room string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	cursor, err := r.Table(room).Changes().Run(sess)
	if err != nil {
		return err
	}
	var msg RoomMessage
	for cursor.Next(&msg) {
		fmt.Println(msg)
	}
	return nil
}

func (p *Rethink) FeedPrivate() error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	cursor, err := r.Table(p.getPrivateMsgTable()).
		Filter(r.Row.Field("Receiver").Eq(p.ctx.username)).
		Changes().
		Run(sess)
	if err != nil {
		return err
	}
	var msg PrivateMessage
	for cursor.Next(&msg) {
		fmt.Println(msg)
	}
	return nil
}
