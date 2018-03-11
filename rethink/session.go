// Package main provides ...
package rethink

import (
	"fmt"

	r "gopkg.in/gorethink/gorethink.v4"
)

type RethinkOptions struct {
	Address  string
	Database string
	Username string
	Password string
}

type Rethink struct {
	address  string
	database string
	username string
	password string
}

func NewRethink(options *RethinkOptions) *Rethink {
	rethink := &Rethink{
		address:  options.Address,
		username: options.Username,
		password: options.Password,
		database: options.Database,
	}
	return rethink
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

func (p *Rethink) Feeds(room string) error {
	sess, err := p.connect()
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {

		func() {
			cursor, err := r.Table(room).Changes().Run(sess)
			if err != nil {
				return err
			}
			var record interface{}
			for cursor.Next(&record) {
				fmt.Println(record)
			}
		}()
	}

	return nil
}
