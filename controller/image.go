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

var dcyImage *dcy.Image = &dcy.Image{}
var connectionPoolImage *redis.Pool = misc.CreateConnectionPool()

type ImageController struct{}

func (t *ImageController) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t.Get(w, r)
	} else if r.Method == "POST" {
		t.Post(w, r)
	} else if r.Method == "DELETE" {
		t.Delete(w, r)
	}
}

func (t *ImageController) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	result := t.getStruct(id)

	dcyImage.Damage(result)
	result.CallCount += 1
	t.updateStruct(result)

	response, _ := json.Marshal(result)

	w.Write(response)
}

func (t *ImageController) Post(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *ImageController) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *ImageController) getStruct(id string) *model.Dcy {
	conn, _ := connectionPoolImage.Dial()
	defer conn.Close()

	reply, err := conn.Do("HGET", "image:content", id)
	content, _ := redis.String(reply, err)

	reply, err = conn.Do("HGET", "image:call_count", id)
	callCount, _ := redis.Int(reply, err)

	imageId, _ := strconv.ParseInt(id, 10, 32)

	return &model.Dcy{int(imageId), int64(callCount), content}
}

func (t *ImageController) updateStruct(s *model.Dcy) {
	conn, _ := connectionPoolImage.Dial()
	defer conn.Close()

	id := strconv.FormatInt(int64(s.Id), 10)

	conn.Do("HSET", "image:call_count", id, s.CallCount)
}
