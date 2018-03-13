// Package auth provides
package auth

// Session user context
type Session struct {
	UserName string
	InRoom   string
}

type Users struct {
	Username  string `gorethink:"username"`
	Password  string `gorethink:"password"`
	CreatedAt int64  `gorethink:"created_at"`
}

var (
	UserSess *Session
)

func newSession(username string) *Session {
	return &Session{
		UserName: username,
		InRoom:   "general",
	}
}

func (p *Session) GetUser() string {
	return p.UserName
}

func IsLogin() bool {
	return UserSess != nil
}

func SetLogin(username string) bool {
	UserSess = newSession(username)
	return true
}
