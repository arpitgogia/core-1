package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
)

func TestNewHook(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	client.UploadURL = url

	fileName := "hooker_test.db"
	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
	}

	testDB, err := bolt.Open(fileName, 0644, nil)
	if err != nil {
		t.Errorf("error opening test database: %v", err)
	}

	testServer := FrontendServer{Database: BoltDB{DB: testDB}}

	mux.HandleFunc("/repos/nihilus/hunger/hooks", func(w http.ResponseWriter, r *http.Request) {
		v := new(github.Hook)
		json.NewDecoder(r.Body).Decode(v)
		fmt.Fprint(w, `{"id":1}`)
	})

	login := "nihilus"
	user := &github.User{Login: &login}
	name := "hunger"
	id := 1
	testRepo := github.Repository{
		Name:  &name,
		Owner: user,
		ID:    &id,
	}

	err = testServer.NewHook(&testRepo, client)
	if err != nil {
		t.Errorf("newhook test failed: %v", err)
	}
}
