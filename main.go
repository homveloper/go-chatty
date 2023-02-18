package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"

	"github.com/globalsign/mgo"
)

var renderer *render.Render
var mongoSession *mgo.Session

const (
	SESSION_PUBLIC_KEY  = "simple_chat_public_key"
	SESSION_PRIVATE_KEY = "simple_chat_private_key"
)

func init() {
	renderer = render.New(render.Options{})
	mongoSession, _ = mgo.Dial("localhost")
}

func main() {

	// 라우터 생성
	router := httprouter.New()

	// 핸들러 정의
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Simple Chat!"})
	})

	router.GET("/login", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "login", nil)
	})

	router.GET("/logout", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// 세션에서 사용자 정보 제거 후 로그인 페이지로 이동
		sessions.GetSession(r).Delete(CURRENT_USER_PUBLIC_KEY)
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	router.GET("/auth/:action/:provider", loginHandler)

	router.GET("/rooms", RetrieveRooms)
	router.POST("/rooms", CreateRoom)

	// negroni 서버 생성
	server := negroni.Classic()
	store := cookiestore.New([]byte(SESSION_PRIVATE_KEY))
	server.Use(sessions.Sessions(SESSION_PUBLIC_KEY, store))
	server.Use(LoginRequired("/login", "/auth"))

	// 미들웨어에 라우터 추가
	server.UseHandler(router)

	// 서버 실행
	server.Run(":3000")

}
