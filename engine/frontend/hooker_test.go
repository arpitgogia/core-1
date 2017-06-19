package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

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

	testServer := FrontendServer{}
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

	defer testServer.CloseBolt()
	err := testServer.OpenBolt()
	if err != nil {
		t.Error(err) // TODO: Flesh out message
	}

	err = testServer.NewHook(&testRepo, client)
	if err != nil {
		t.Error(err) // TODO: Flesh out message
	}
}
