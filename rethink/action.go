// Package main provides ...
package rethink

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/koolay/console-chat/auth"
	r "gopkg.in/gorethink/gorethink.v4"
)

var (
	roomPrefix   = "room_"
	database     = "chat"
	address      = "172.105.233.187:28015"
	RethinkActor *Rethink
)

func init() {
	options := &Options{
		Address:  address,
		Database: database,
		Username: "",
		Password: "",
	}
	RethinkActor = NewRethink(options)
}

// Options options
type Options struct {
	Address  string
	Database string
	Username string
	Password string
}

// Rethink driver
type Rethink struct {
	address  string
	database string
	username string
	password string
}

// NewRethink instance driver
func NewRethink(options *Options) *Rethink {
	rethink := &Rethink{
		address:  options.Address,
		username: options.Username,
		password: options.Password,
		database: options.Database,
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

// Test test
func (p *Rethink) Test() error {
	//sess, err := p.connect()
	//go func() {
	//for {
	//m := PublicMessage{
	//Sender:    "debug",
	//Content:   "hello" + string(time.Now().Unix()),
	//CreatedAt: time.Now().Unix(),
	//}
	//_, err = r.Table(p.getPublicMsgTable()).Insert(m).RunWrite(sess)
	//time.Sleep(1 * time.Second)
	//}
	//}()
	//user := auth.Users{
	//Username: "ooo",
	//Password: "123",
	//}
	//fmt.Println(result.Created)
	go func() {
		for {
			p.SendPublicMessage("foo" + time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(1 * time.Second)

		}
	}()
	go func() {
		p.FeedsPublic()
	}()
	for {
		record := <-ChatRecordChan
		fmt.Println("chan->", "id:", record.Id, ",", record.Content, ",", record.CreatedAt, ",", record.Sender)
	}
	//return err
	return nil
}

func (p *Rethink) Login(username string, password string) error {
	sess, err := p.connect()
	result, err := r.Table("users").Filter(r.Row.Field("username").Eq(username)).Run(sess)
	if err != nil {
		return err
	}

	var user auth.Users
	_, err = result.Peek(&user)
	if err != nil {
		return err
	}
	if !checkPasswordHash(password, user.Password) {
		return errors.New("username or password not match")
	}
	auth.SetLogin(username)
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
	var user auth.Users
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
	user = auth.Users{
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
	auth.UserSess.InRoom = room
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
	roomName := fmt.Sprintf("%s%s", roomPrefix, room)
	_, err = r.DB(p.database).TableCreate(roomName).RunWrite(sess)
	if err != nil {
		return err
	}
	fmt.Printf("Create room %s Successfully\n", room)
	return nil
}

func (p *Rethink) GetAllRooms() (rooms []string, err error) {
	sess, err := p.connect()
	if err != nil {
		return
	}
	var tables []string
	result, err := r.TableList().Run(sess)
	err = result.All(&tables)
	if err != nil {
		return
	}
	for _, tableName := range tables {
		if strings.HasPrefix(tableName, roomPrefix) {
			rooms = append(rooms, strings.TrimPrefix(tableName, roomPrefix))
		}
	}
	return
}

type ScoreEntry struct {
	ID         string `gorethink:"id,omitempty"`
	PlayerName string
	Score      int
}

func (p *Rethink) SendPublicMessage(message string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	msg := PublicMessage{}
	msg.Content = message
	msg.Sender = auth.UserSess.GetUser()
	msg.CreatedAt = time.Now().Unix()
	_, err = r.Table(p.getPublicMsgTable()).Insert(msg).RunWrite(sess)
	if err != nil {
		panic(err)
	}
	return err
}

func (p *Rethink) SendRoomMessage(message string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	msg := RoomMessage{}
	msg.Content = message
	msg.Sender = auth.UserSess.GetUser()
	msg.CreatedAt = time.Now().Unix()
	_, err = r.Table(auth.UserSess.InRoom).Insert(msg).RunWrite(sess)
	return err
}

func (p *Rethink) SendPrivateMessage(toUser string, message string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	msg := PrivateMessage{}
	msg.Sender = auth.UserSess.GetUser()
	msg.Receiver = toUser
	msg.Content = message
	msg.CreatedAt = time.Now().Unix()
	_, err = r.Table(p.getPrivateMsgTable()).Insert(msg).RunWrite(sess)
	return err
}

func (p *Rethink) FeedsPublic() error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	cursor, err := r.Table(p.getPublicMsgTable()).Changes().Field("new_val").Run(sess)
	if err != nil {
		return err
	}
	var msg PublicMessage
	for cursor.Next(&msg) {
		record := ChatRecordPoll.Get().(*ChatRecord)
		record.Channel = ChannelPublic
		record.Content = msg.Content
		record.CreatedAt = msg.CreatedAt
		record.Id = msg.Id
		record.Sender = msg.Sender
		ChatRecordChan <- record
	}
	return nil
}

func (p *Rethink) FeedsRoom(room string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	cursor, err := r.Table(room).Changes().Field("new_val").Run(sess)
	if err != nil {
		return err
	}
	var msg RoomMessage
	for cursor.Next(&msg) {
		record := ChatRecordPoll.Get().(*ChatRecord)
		record.Channel = ChannelRoom
		record.Content = msg.Content
		record.CreatedAt = msg.CreatedAt
		record.Id = msg.Id
		record.Sender = msg.Sender
		ChatRecordChan <- record
	}
	return nil
}

func (p *Rethink) FeedPrivate() error {
	sess, err := p.connect()
	if err != nil {
		return err
	}
	cursor, err := r.Table(p.getPrivateMsgTable()).
		Filter(r.Row.Field("Receiver").Eq(auth.UserSess.GetUser())).
		Changes().
		Field("new_val").
		Run(sess)
	if err != nil {
		return err
	}
	var msg PrivateMessage
	for cursor.Next(&msg) {
		record := ChatRecordPoll.Get().(*ChatRecord)
		record.Channel = ChannelPrivate
		record.Content = msg.Content
		record.CreatedAt = msg.CreatedAt
		record.Id = msg.Id
		record.Sender = msg.Sender
		ChatRecordChan <- record
	}
	return nil
}
