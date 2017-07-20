package ingestor

import (
	"net/http"

	"go.uber.org/zap"

	"coralreefci/engine/frontend"
	"coralreefci/utils"
)

type IngestorServer struct {
	Server          http.Server
	RepoInitializer RepoInitializer
}

func (i *IngestorServer) activateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != frontend.BackendSecret {
			utils.AppLog.Error("failed validating frontend-backend secret")
			http.Redirect(w, r, "/", http.StatusForbidden)
			return
		}
	})
}

func (i *IngestorServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hook", collectorHandler())
	mux.Handle("/activate-repos-ingestor", i.activateHandler())
	return mux
}

func (i *IngestorServer) Start() {
	i.RepoInitializer = RepoInitializer{}
	i.Server = http.Server{Addr: "127.0.0.1:8030", Handler: i.routes()}
	err := i.Server.ListenAndServe()
	if err != nil {
		utils.AppLog.Error("ingestor server failed to start", zap.Error(err))
	}
}

func (i *IngestorServer) Restart() {}

func (i *IngestorServer) ContinuityCheck() {}

func (i *IngestorServer) Stop() {
	//TODO: Closing the server down is a needed operation that will be added.
}
