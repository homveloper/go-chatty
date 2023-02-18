package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/urfave/negroni"
)

const (
	NEXT_PAGE_KEY     = "next_page"
	AUTH_SECURITY_KEY = "auth_security_key"
)

type config struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUrl  string `json:"redirect_url"`
}

func init() {
	var conf config

	file, _ := ioutil.ReadFile("auth.json")
	_ = json.Unmarshal([]byte(file), &conf)

	gomniauth.SetSecurityKey(AUTH_SECURITY_KEY)
	gomniauth.WithProviders(
		google.New(
			conf.ClientId,
			conf.ClientSecret,
			conf.RedirectUrl,
		),
	)
}

func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	provider := ps.ByName("provider")
	session := sessions.GetSession(r)

	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, loginUrl, http.StatusFound)

	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 콜백 결과로부터 사용자 정보 확인
		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 사용자 정보를 세션에 갱신
		u := &User{
			UID:       user.Data().Get("id").MustStr(),
			Name:      user.Name(),
			Email:     user.Email(),
			AvatarUrl: user.AvatarURL(),
		}

		SetCurrentUser(r, u)
		http.Redirect(w, r, session.Get(NEXT_PAGE_KEY).(string), http.StatusFound)
	default:
		http.Error(w, "Auth "+action+" not supported", http.StatusBadRequest)
	}
}

func LoginRequired(ignore ...string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// 로그인이 필요한 페이지인지 확인
		for _, path := range ignore {
			if strings.HasPrefix(r.URL.Path, path) {
				next(w, r)
				return
			}
		}

		// 현재 유저 정보 가져오기
		user := GetCurrentUser(r)

		// 유저 정보가 유효하면 만료 시간을 갱신하고 다음 헨들러 실행
		if user != nil && user.Valid() {
			SetCurrentUser(r, user)
			next(w, r)
			return
		}

		// 유저 정보가 유효하지 않으면 로그인 페이지로 리다이렉트
		SetCurrentUser(r, nil)
		sessions.GetSession(r).Set(NEXT_PAGE_KEY, r.URL.RequestURI())
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
