package backend

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"golang.org/x/oauth2"

	"coralreefci/engine/frontend"
)

type ActiveRepos struct {
	sync.RWMutex
	Actives map[int]*ArchRepo
}

type BackendServer struct {
	Server   http.Server
	Database MemSQL
	Repos    *ActiveRepos
}

// TODO: Possibly move into separate file.
func (bs *BackendServer) activateHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != frontend.BackendSecret {
		// TODO: Something to handle this scenario - no redirect.
		return
	}
	repoIDString := r.FormValue("repos")
	repoID, err := strconv.Atoi(repoIDString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if bs.Repos.Actives[repoID] == nil {
		db, err := bolt.Open("../frontend/storage.db", 0644, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer db.Close()

		boltDB := frontend.BoltDB{DB: db}

		byteToken, err := boltDB.Retrieve("token", repoID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		token := oauth2.Token{}
		if err := json.Unmarshal(byteToken, &token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		bs.NewArchRepo(repoID)
		bs.NewClient(repoID, &token)
		bs.NewModel(repoID)
	}
}

func (bs *BackendServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/activate-repos", bs.activateHandler)
	bs.Server = http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	bs.Server.ListenAndServe()

	db, err := bolt.Open("storage.db", 0644, nil)
	defer db.Close()
	boltDB := frontend.BoltDB{DB: db}
	keys, tokens, err := boltDB.RetrieveBulk("tokens")
	if err != nil {
		panic(err) // TODO: Implement proper error handling.
	}

	for i := 0; i < len(keys); i++ {
		key, err := strconv.Atoi(string(keys[i]))
		if err != nil {
			panic(err) // TODO: Implement proper error handling.
		}
		token := oauth2.Token{}
		if err := json.Unmarshal(tokens[i], &token); err != nil {
			panic(err) // TODO: Implement proper error handling.
		}
		bs.NewArchRepo(key)
		bs.NewClient(key, &token)
		bs.NewModel(key)
	}
}

func (bs *BackendServer) OpenSQL() {
	bs.Database.Open()
}

func (bs *BackendServer) CloseSQL() {
	bs.Database.Close()
}

func (bs *BackendServer) Timer() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {

			data, err := bs.Database.Read()
			if err != nil {
				panic(err) // TODO: Implement proper error handling.
			}

			bs.Dispatcher(10)
			Collector(data)

			// TODO: Complete this method.

		}
	}()
}
