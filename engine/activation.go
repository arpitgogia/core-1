package engine

import (
	"net/http"
	"net/url"

	"coralreefci/engine/frontend"
	"coralreefci/utils"
)

const (
	destinationBase = "http://127.0.0.1"
	destinationPort = ":8090"
	destinationEnd  = "/activate-repos-ingestor"
)

type ActivationServer struct {
	Server http.Server
}

func (as *ActivationServer) activationServerHandler(w http.ResponseWriter, r *http.Request) {
	secret := frontend.BackendSecret
	if r.FormValue("state") != secret {
		utils.AppLog.Error("failed validating frontend-backend secret")
		http.Error(w, "invalid secret", http.StatusForbidden)
		return
	}
	repoID := r.FormValue("repos")
	resp, err := http.PostForm(destinationBase+destinationPort+destinationEnd, url.Values{
		"state": {secret},
		"repos": {repoID},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		defer resp.Body.Close()
	}

}

func (as *ActivationServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/activate", as.activationServerHandler)

	as.Server = http.Server{
		Addr:    "127.0.0.1:8090",
		Handler: mux,
	}
	as.Server.ListenAndServe()
}
