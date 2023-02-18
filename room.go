package main

import (
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
)

type Room struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}

func (r *Room) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&r.Name: "name",
	}
}

func CreateRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// ...
	r := &Room{}
	errs := binding.Bind(req, r)
	if errs.Handle(w) {
		return
	}

	// MongoDB 세션 생성
	session := mongoSession.Clone()
	defer session.Close()

	// MongoDB 아이디 생성
	r.ID = bson.NewObjectId()
	// room 정보 저장을 위한 몽고DB 컬렉션 생성
	collection := session.DB("simple_chat").C("rooms")

	// Room 컬렉션 생성 및 저장
	if err := collection.Insert(r); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	// 응답 전송
	renderer.JSON(w, http.StatusCreated, r)
}

func RetrieveRooms(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// MongoDB 세션 생성
	session := mongoSession.Clone()
	defer session.Close()

	var rooms []Room
	// 모든 Room 조회
	err := session.DB("simple_chat").C("rooms").Find(nil).All(&rooms)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
	}

	// 응답 전송
	renderer.JSON(w, http.StatusOK, rooms)
}
