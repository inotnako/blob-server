package mongostorage

import (
	"testing"
	"gopkg.in/mgo.v2"
	"os"
	"bytes"
)

var s *MongoStorage

func TestMain(m *testing.M) {
	session, _ := mgo.Dial("localhost")
	session.DB("").C("mongostorage_test").DropCollection()
	session.Close()
	s, _ = Start("localhost", "mongostorage_test")
	run := m.Run()
	defer s.Stop()
	os.Exit(run)
}

func Test(t *testing.T) {
	lst, err := s.GetList()
	if (err != nil) {
		t.Error("GetList failed")
	}
	if (len(lst) != 0) {
		t.Error("non-empty file list at start")
	}

	id, err := s.Post(bytes.NewBufferString("filecontent"))
	if (err != nil) {
		t.Error("Post failed")
	}

	var buf bytes.Buffer
	err = s.Get(id, &buf)
	if (err != nil) {
		t.Error("Get failed")
	}
	if (string(buf.Bytes()) != "filecontent") {
		t.Error("Get wrong data")
	}

	iderr := s.Get("badid", &buf)
	if (iderr == nil || iderr.NotFound() || !iderr.IllFormed()) {
		t.Error("Wrong error returned")
	}

	lst, err = s.GetList()
	if (err != nil) {
		t.Error("GetList failed")
	}
	if (len(lst) != 1 || lst[0] != id) {
		t.Error("wrong file list")
	}

	iderr = s.Delete("badid")
	if (iderr == nil || iderr.NotFound() || !iderr.IllFormed()) {
		t.Error("Wrong error returned")
	}

	iderr = s.Delete(id)
	if (iderr != nil) {
		t.Error("Delete failed")
	}

	iderr = s.Delete(id)
	if (iderr == nil || !iderr.NotFound() || iderr.IllFormed()) {
		t.Error("Wrong error returned")
	}

	iderr = s.Get(id, &buf)
	if (iderr == nil || !iderr.NotFound() || iderr.IllFormed()) {
		t.Error("Wrong error returned")
	}

	lst, err = s.GetList()
	if (err != nil) {
		t.Error("GetList failed")
	}
	if (len(lst) != 0) {
		t.Error("non-empty file list")
	}
}
