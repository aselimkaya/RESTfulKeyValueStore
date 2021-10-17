package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/aselimkaya/RESTfulKeyValueStore/src/repository"
)

//TestPost simple test for adding key value pairs. In the first part, we assume that the JSON file is empty. In the second part, the given key already inserted in part 1.
func TestPost(t *testing.T) {
	t.Run("Post key1 value1", func(t *testing.T) {
		l := log.New(os.Stdout, "test", log.LstdFlags)

		absPath, _ := filepath.Abs("../../")
		l.Println(absPath)

		s := New(l, absPath)

		e := repository.Entry{
			Key:   "key1",
			Value: "value1",
		}

		requestByte, _ := json.Marshal(e)

		request, _ := http.NewRequest(http.MethodPost, "/entry", bytes.NewReader(requestByte))
		response := httptest.NewRecorder()

		s.AddEntry(response, request)

		got := response.Body.String()
		want := `{"message":"Key value pair added successfully","status":200}`

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Post key1 value1 again", func(t *testing.T) {
		l := log.New(os.Stdout, "test", log.LstdFlags)

		absPath, _ := filepath.Abs("../../")
		l.Println(absPath)

		s := New(l, absPath)

		e := repository.Entry{
			Key:   "key1",
			Value: "value1",
		}

		requestByte, _ := json.Marshal(e)

		request, _ := http.NewRequest(http.MethodPost, "/entry", bytes.NewReader(requestByte))
		response := httptest.NewRecorder()

		s.AddEntry(response, request)

		got := response.Body.String()
		want := `{"message":"Key already exists, value will be upated","status":200}`

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

//TestGet is simple test function for testing key1's presence. Make sure that adding key1 before.
func TestGet(t *testing.T) {
	t.Run("Get key1 value", func(t *testing.T) {
		l := log.New(os.Stdout, "test", log.LstdFlags)

		absPath, _ := filepath.Abs("../../")
		l.Println(absPath)

		s := New(l, absPath)

		request, _ := http.NewRequest(http.MethodGet, "/entry?key=key1", nil)
		response := httptest.NewRecorder()

		s.GetEntry(response, request)

		got := response.Body.String()
		want := `{"message":"{\"key\":\"key1\",\"value\":\"value1\"}","status":200}`

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

//TestDelete controls if map content cleared after flush
func TestDelete(t *testing.T) {
	t.Run("Flush store", func(t *testing.T) {
		l := log.New(os.Stdout, "test", log.LstdFlags)

		absPath, _ := filepath.Abs("../../")
		l.Println(absPath)

		s := New(l, absPath)

		repository.Flush(s.jsonFilePath)

		got := len(repository.GetStore())
		want := 0

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
