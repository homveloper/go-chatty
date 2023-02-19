package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/julienschmidt/httprouter"
)

const (
	MESSAGE_FETCH_SIZE = 10
)

type Message struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	RoomID   bson.ObjectId `bson:"room_id" json:"room_id"`
	Content  string        `bson:"content" json:"content"`
	CreateAt time.Time     `bson:"create_at" json:"create_at"`
	User     *User         `bson:"user" json:"user"`
}

func (m *Message) Create() error {
	// MongoDB 세션 생성
	session := mongoSession.Clone()
	defer session.Close()

	// MongoDB 아이디 생성
	m.ID = bson.NewObjectId()
	// 메시지 생성 시간
	m.CreateAt = time.Now()

	// room 정보 저장을 위한 몽고DB 컬렉션 생성
	collection := session.DB("simple_chat").C("messages")

	// Room 컬렉션 생성 및 저장
	if err := collection.Insert(m); err != nil {
		return err
	}

	return nil
}

func retrievMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = MESSAGE_FETCH_SIZE
	}

	var messages []Message

	err = session.DB("simple_chat").C("messages").
		Find(bson.M{"room_id": bson.ObjectIdHex(ps.ByName("room_id"))}).
		Sort("-_id").
		Limit(limit).
		All(&messages)

	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	renderer.JSON(w, http.StatusOK, messages)
}
