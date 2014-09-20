package controller

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/lafaulx/decay/dcy"
	"github.com/lafaulx/decay/misc"
	"github.com/lafaulx/decay/model"
	"net/http"
	"strconv"
)

var dcyText *dcy.Text = &dcy.Text{}
var connectionPoolText *redis.Pool = misc.CreateConnectionPool()

type TextController struct{}

func (t *TextController) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t.Get(w, r)
	} else if r.Method == "POST" {
		t.Post(w, r)
	} else if r.Method == "DELETE" {
		t.Delete(w, r)
	}
}

func (t *TextController) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	result := t.getStruct(id)

	dcyText.Damage(result)
	result.CallCount += 1
	t.updateStruct(result)

	response, _ := json.Marshal(result)

	w.Write(response)
}

func (t *TextController) Post(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mdl model.Dcy
	decoder.Decode(&mdl)

	mdl.Id = misc.GenerateID()
	mdl.CallCount = 0

	t.updateStruct(&mdl)

	response, _ := json.Marshal(mdl)

	w.Write(response)
}

func (t *TextController) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *TextController) getStruct(id string) *model.Dcy {
	conn, _ := connectionPoolText.Dial()
	defer conn.Close()

	reply, err := conn.Do("HGET", "text:content", id)
	content, _ := redis.String(reply, err)

	reply, err = conn.Do("HGET", "text:call_count", id)
	callCount, _ := redis.Int(reply, err)

	textId, _ := strconv.ParseInt(id, 10, 32)

	return &model.Dcy{int(textId), int64(callCount), content}
}

func (t *TextController) updateStruct(s *model.Dcy) {
	conn, _ := connectionPoolText.Dial()
	defer conn.Close()

	id := strconv.FormatInt(int64(s.Id), 10)

	conn.Do("HSET", "text:content", id, s.Content)
	conn.Do("HSET", "text:call_count", id, s.CallCount)
}
