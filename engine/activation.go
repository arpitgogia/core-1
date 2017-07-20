package engine

import (
	"net/http"
	"net/url"

	"coralreefci/engine/frontend"
	"coralreefci/utils"
)

var (
	destinationBase  = "http://127.0.0.1"
	destinationPorts = []string{":8020", ":8030"}
	destinationEnd   = "/activate-repos-ingestor"
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
	token := r.FormValue("token")
	for i := range destinationPorts {
		resp, err := http.PostForm(destinationBase+destinationPorts[i]+destinationEnd, url.Values{
			"state": {secret},
			"repos": {repoID},
			"token": {token},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			defer resp.Body.Close()
		}
	}
}

func (as *ActivationServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/activate", as.activationServerHandler)

	as.Server = http.Server{
		Addr:    "127.0.0.1:8010",
		Handler: mux,
	}
	as.Server.ListenAndServe()
}
