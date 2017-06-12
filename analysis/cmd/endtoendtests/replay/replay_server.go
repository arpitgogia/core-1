package replay

import (
	"bytes"
	"coralreefci/engine/ingestor"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	gzip "github.com/klauspost/pgzip"
	"io"
	"io/ioutil"
	"net/http"
	//"net/url"
	// "reflect"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	//"coralreefci/engine/onboarder"
	// "coralreefci/models"
	"runtime/debug"
	"time"
)

const (
	secretKey = "test"
	localPath = "http://localhost:8000/"
)

//var modelList = []*onboarder.ArchModel{}

var webhooksplit float32 = 0.5

type Issue struct {
	Payload json.RawMessage `json:issue`
}

type BacktestServer struct {
	client http.Client
	DB     *ingestor.Database
	server http.Server
	//onboarder.RepoServer
	events        []*ingestor.Event
	WebhookEvents []*ingestor.Event
}

func (b *BacktestServer) routes() *mux.Router {
	gorilla := mux.NewRouter()
	gorilla.HandleFunc("/repos/{org}/{repo}/issues", b.getIssues)
	gorilla.HandleFunc("/repos/{org}/{repo}/pulls", b.getPulls)
	gorilla.HandleFunc("/stream", b.streamWebhooks)
	return gorilla
}

func (b *BacktestServer) Start() {
	b.server = http.Server{Addr: "127.0.0.1:8000", Handler: b.routes()}
	err := b.server.ListenAndServe()
	if err != nil {
		fmt.Println("BacktestServer", err)
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	b.client = http.Client{Transport: tr}
}

func (b *BacktestServer) LoadArchive(path string) {
	fn, _ := os.Stat(path)
	switch mode := fn.Mode(); {
	case mode.IsDir():
		files, _ := ioutil.ReadDir(path)
		loadedFiles := 0
		//inserts := 0
		totalFiles := len(files)
		err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			parseErr := b.parseFile(path)
			if parseErr == nil {
				loadedFiles++
				fmt.Printf("Loaded %d out of %d files\n", loadedFiles, totalFiles)
			} else {
				fmt.Println(parseErr)
			}
			/*
				if inserts%2 == 0 {
					debug.FreeOSMemory()
					b.DB.FlushBackTestTable()
				}*/

			if loadedFiles > 1 && loadedFiles%100 == 0 {
				debug.FreeOSMemory()
				time.Sleep(1 * time.Minute)
			}
			if loadedFiles > 0 && loadedFiles%5 == 0 {
				debug.FreeOSMemory()
				b.DB.FlushBackTestTable()
				b.DB.BulkInsertBacktestEvents(b.events)
				debug.FreeOSMemory()

				for i := 0; i < len(b.events); i++ {
					RecycleEvent(b.events[i])
					//inserts++
				}
				time.Sleep(1 * time.Second)
				b.events = []*ingestor.Event{}
				fmt.Printf("Inserted %d out of %d files\n", loadedFiles, totalFiles)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking directory: %v", err)
		}
	case mode.IsRegular():
		b.parseFile(path)
	default:
		fmt.Println("Unrecognized argument; provide a file or directory")
	}
	if len(b.events) > 0 {
		fmt.Printf("Inserted remaining records", len(b.events))
		b.DB.BulkInsertBacktestEvents(b.events)
		for i := 0; i < len(b.events); i++ {
			RecycleEvent(b.events[i])
		}
		b.events = []*ingestor.Event{}
	}
}

func (b *BacktestServer) parseFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()
	jd := json.NewDecoder(gr)
	jd.UseNumber()
	for {
		e := GetEvent()
		if err := jd.Decode(&e); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		switch e.Type {
		case "IssuesEvent", "PullRequestEvent":
			m := e.Payload.(map[string]interface{})
			e.Action = m["action"].(string) // Workaround
			b.events = append(b.events, e)
		default:
			RecycleEvent(e)
		}
	}
	return nil
}

/*
func (b *BacktestServer) loadRepos() {
	// TODO: parsing some input (likely a file w/ JSON content)
	// - defines the desired repos to run (specific GitHub Repositories)
	repos := []github.Repository{} // TEMPORARY

	client := github.NewClient(nil)
	u, _ := url.Parse(localPath)
	client.BaseURL = u
	client.UploadURL = u

	for i := 0; i < len(repos); i++ {
		b.Repos[i] = &onboarder.ArchRepo{
			Repo:   &repos[i],
			Hive:   &onboarder.ArchHive{Models: modelList},
			Client: client,
		}
	}
}

func (b *BacktestServer) backtestHandler(w http.ResponseWriter, r *http.Request) {
	jsonRepo := r.FormValue("repository")
	repo := github.Repository{}
	if err := json.Unmarshal([]byte(jsonRepo), &repo); err != nil {
		fmt.Println("Error unmarshalling JSON into repo")
	}

	// TODO: There still needs to be logic in what is pulled out of the Archive
	//       and how best to return it to the writer (likely via fmt.Fprint).
	//       The archive has potential to be relatively large and looping
	//       through it could be expensive particularly if this handler is
	//       being hit for every repo in a given backtest run.
}*/

func (b *BacktestServer) backtestPredict(w http.ResponseWriter, r *http.Request) {

}

//TODO Refactor
func (b *BacktestServer) getIssues(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := vars["org"]
	repo := vars["repo"]
	queryParams := ingestor.EventQuery{Type: ingestor.Issue, Repo: org + "/" + repo}
	events, _ := b.DB.ReadBacktestEvents(queryParams)
	issues := make([]interface{}, int(float32(len(events))*webhooksplit))
	if webhooksplit == 1 {
		for i := 0; i < len(events); i++ {
			m := events[i].Payload.(map[string]interface{})
			issue := m["issue"]
			issues[i] = issue
		}
	} else {
		for i := 0; i < int(float32(len(events))*webhooksplit); i++ {
			m := events[i].Payload.(map[string]interface{})
			issue := m["issue"]
			issues[i] = issue
		}
		for i := int(float32(len(events)) * webhooksplit); i < len(events); i++ {
			event := events[i]
			b.WebhookEvents = append(b.WebhookEvents, &event)
		}
	}
	payload, _ := json.Marshal(&issues)
	w.Write(payload)
}

//TODO Refactor
func (b *BacktestServer) getPulls(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := vars["org"]
	repo := vars["repo"]
	queryParams := ingestor.EventQuery{Type: ingestor.PullRequest, Repo: org + "/" + repo}
	events, _ := b.DB.ReadBacktestEvents(queryParams)
	pulls := make([]interface{}, int(float32(len(events))*webhooksplit))
	if webhooksplit == 1 {
		for i := 0; i < len(events); i++ {
			m := events[i].Payload.(map[string]interface{})
			pull := m["pull_request"]
			pulls[i] = pull
		}
	} else {
		for i := 0; i < int(float32(len(events))*webhooksplit); i++ {
			m := events[i].Payload.(map[string]interface{})
			pull := m["pull_request"]
			pulls[i] = pull
		}
		for i := int(float32(len(events)) * webhooksplit); i < len(events); i++ {
			event := events[i]
			b.WebhookEvents = append(b.WebhookEvents, &event)
		}
		fmt.Println("Streaming", len(b.WebhookEvents))
	}
	payload, _ := json.Marshal(&pulls)
	w.Write(payload)
}

func (b *BacktestServer) streamWebhooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Streaming WebHooks", len(b.WebhookEvents))

	for i := 0; i < len(b.WebhookEvents); i++ {
		m := b.WebhookEvents[i].Payload.(map[string]interface{})
		m["repository"] = &b.WebhookEvents[i].Repo // Workaround
		payload, _ := json.Marshal(m)
		var event string
		if b.WebhookEvents[i].Type == "PullRequestEvent" {
			event = "pull_request"
		} else {
			event = "issues"
		}
		event = event
		payload = payload
		bs.HTTPPost(bytes.NewBuffer(payload), event)
	}
}

func (b *BacktestServer) StreamWebhookEvents() {
	fmt.Println("Streaming", len(b.WebhookEvents))
	for i := 0; i < len(b.WebhookEvents); i++ {
		m := b.WebhookEvents[i].Payload.(map[string]interface{})
		m["repository"] = &b.WebhookEvents[i].Repo // Workaround
		payload, _ := json.Marshal(m)
		var event string
		if b.WebhookEvents[i].Type == "PullRequestEvent" {
			event = "pull_request"
		} else {
			event = "issues"
		}
		bs.HTTPPost(bytes.NewBuffer(payload), event)
	}
}

func (b *BacktestServer) HTTPPost(payload *bytes.Buffer, event string) {
	req, err := http.NewRequest("POST", "http://localhost:8080/hook", payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("X-Github-Event", event)
	req.Header.Set("X-GitHub-Delivery", "placeholder")
	req.Header.Set("content-type", "application/json")
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write(payload.Bytes())
	sig := "sha1=" + hex.EncodeToString(mac.Sum(nil))
	req.Header.Set("X-Hub-Signature", sig)

	b.client.Do(req)
}
