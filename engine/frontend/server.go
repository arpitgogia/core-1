package frontend

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
)

type FrontendServer struct {
	Server   http.Server
	Database BoltDB
}

func (fs *FrontendServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", mainHandler)
	mux.HandleFunc("/login", githubLoginHandler)
	mux.HandleFunc("/github_oauth_cb", fs.githubCallbackHandler)
	mux.HandleFunc("/setup_complete", completeHandle)
	return mux
}

func (fs *FrontendServer) Start() {
	fs.Server = http.Server{Addr: "127.0.0.1:8080", Handler: fs.routes()}
	// TODO: Add in logging and remove print statement.
	err := fs.Server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (fs *FrontendServer) Stop() {
	// TODO: Closing the server down is a needed operation that will be added.
}

func (fs *FrontendServer) OpenBolt() error {
	boltDB, err := bolt.Open("storage.db", 0644, nil)
	if err != nil {
		return err
	}
	fs.Database = BoltDB{DB: boltDB}
	return nil
}

func (fs *FrontendServer) CloseBolt() {
	fs.Database.DB.Close()
}
