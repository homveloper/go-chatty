package main

import (
	"encoding/json"
	"net/http"
	"time"

	sessions "github.com/goincremental/negroni-sessions"
)

const (
	CURRENT_USER_PUBLIC_KEY = "oauth2_current_user" // 현재 로그인한 사용자의 고유 키
	SESSION_DURATION        = time.Hour             // 로그인 세션 유지 시간
)

type User struct {
	UID       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"avatar"`
	Expired   time.Time `json:"expired"`
}

func (u *User) Valid() bool {
	return time.Until(u.Expired) > 0
}

func (u *User) Refresh() {
	u.Expired = time.Now().Add(SESSION_DURATION)
}

func GetCurrentUser(r *http.Request) *User {

	session := sessions.GetSession(r)

	if session.Get(CURRENT_USER_PUBLIC_KEY) != nil {
		return nil
	}

	data := session.Get(CURRENT_USER_PUBLIC_KEY).([]byte)
	var user User
	json.Unmarshal(data, &user)

	return &user
}

func SetCurrentUser(r *http.Request, user *User) {
	if user != nil {
		user.Refresh()
	}

	session := sessions.GetSession(r)
	data, _ := json.Marshal(user)
	session.Set(CURRENT_USER_PUBLIC_KEY, data)
}
