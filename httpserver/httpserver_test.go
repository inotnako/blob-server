package httpserver

import (
	"blob-server/storage"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

type MockStorage struct {
	filelist      []string
	postedfile    string
	deletedfileid string
}

var mockStorage MockStorage

func (storage MockStorage) Stop() {
}

func (storage *MockStorage) Post(reader io.Reader) (string, error) {
	b, _ := ioutil.ReadAll(reader)
	storage.postedfile = string(b)
	return "idcreated", nil
}

func (storage MockStorage) Get(id string, writer io.Writer) storage.IdRequestError {
	if id == "badid" {
		return IdRequestError{false, true}
	} else if id == "nosuchid" {
		return IdRequestError{true, false}
	} else {
		writer.Write([]byte("file content"))
		return nil
	}
}

func (storage MockStorage) GetList() ([]string, error) {
	return storage.filelist, nil
}

func (storage *MockStorage) Delete(id string) storage.IdRequestError {
	if id == "badid" {
		return IdRequestError{false, true}
	} else if id == "nosuchid" {
		return IdRequestError{true, false}
	} else {
		storage.deletedfileid = id
		return nil
	}
}

type IdRequestError struct {
	notFound  bool
	illFormed bool
}

func (err IdRequestError) NotFound() bool {
	return err.notFound
}

func (err IdRequestError) IllFormed() bool {
	return err.illFormed
}

func (err IdRequestError) Error() string {
	return "xxx error"
}

func TestMain(m *testing.M) {
	go Serve("localhost:5555", &mockStorage)
	os.Exit(m.Run())
}

func TestBadUrls(t *testing.T) {
	cases := []struct {
		method, url string
	}{
		{"GET", "http://localhost:5555/"},
		{"GET", "http://localhost:5555/api/v1/fail"},
		{"POST", "http://localhost:5555/api/v1/file/:0123456789abcdef01234567"},
		{"DELETE", "http://localhost:5555/api/v1/file"},
		{"GET", "http://localhost:5555/api/v1/file/:badid"},
		{"DELETE", "http://localhost:5555/api/v1/file/:"},
	}
	for _, c := range cases {
		req, err := http.NewRequest(c.method, c.url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 400 {
			t.Error(c.method, "for url", c.url, "status is ", resp.StatusCode, "not 400")
		}
	}
}

func TestGetFile(t *testing.T) {
	resp, err := http.Get("http://localhost:5555/api/v1/file/:0123456789abcdef01234578")
	if err != nil || resp.StatusCode != 200 {
		t.Fail()
	} else {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil || string(b) != "file content" {
			t.Fail()
		}
	}

	resp, err = http.Get("http://localhost:5555/api/v1/file/:nosuchid")
	if err != nil || resp.StatusCode != 404 {
		t.Fail()
	}
}

func TestGetFileList(t *testing.T) {
	mockStorage.filelist = []string{"abc", "def"}
	resp, err := http.Get("http://localhost:5555/api/v1/file")
	if err != nil || resp.StatusCode != 200 {
		t.Fail()
	} else {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil || string(b) != "[\"abc\",\"def\"]" {
			t.Error("unexpected get file list response", string(b))
		}
	}
}

func TestPostFile(t *testing.T) {
	resp, err := http.Post("http://localhost:5555/api/v1/file", "", bytes.NewBufferString("postcontent"))
	if err != nil || resp.StatusCode != 200 {
		t.Fail()
	} else {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil || string(b) != "idcreated" {
			t.Error("unexpected post file response", string(b))
		}
		if mockStorage.postedfile != "postcontent" {
			t.Error("was not posted", mockStorage.postedfile)
		}
	}
}

func TestDeleteFile(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:5555/api/v1/file/:todelete", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		t.Fail()
	} else {
		if mockStorage.deletedfileid != "todelete" {
			t.Error("was not deleted", mockStorage.deletedfileid)
		}
	}

	req, err = http.NewRequest("DELETE", "http://localhost:5555/api/v1/file/:nosuchid", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 404 {
		t.Fail()
	}
}
