package ingestor

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"

	"go.uber.org/zap"

	"core/utils"
)

// HeuprInstallationEvent is a workaround for the Github API limitation. This
// is required to wrap HeuprInstallation
type HeuprInstallationEvent struct {
	// The action that was performed. Can be either "created" or "deleted".
	Action            *string            `json:"action,omitempty"`
	Sender            *github.User       `json:"sender,omitempty"`
	HeuprInstallation *HeuprInstallation `json:"installation,omitempty"`
	Repositories      []HeuprRepository  `json:"repositories,omitempty"`
}

// HeuprInstallation is a workaround for the Github API limitation. The
// go-github library is missing a repositories field.
type HeuprInstallation struct {
	ID              *int64       `json:"id,omitempty"`
	Account         *github.User `json:"account,omitempty"`
	AppID           *int         `json:"app_id,omitempty"`
	AccessTokensURL *string      `json:"access_tokens_url,omitempty"`
	RepositoriesURL *string      `json:"repositories_url,omitempty"`
	HTMLURL         *string      `json:"html_url,omitempty"`
}

type HeuprRepository struct {
	ID       *int64  `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

type HeuprInstallationRepositoriesEvent struct {
	// The action that was performed. Can be either "added" or "removed".
	Action              *string              `json:"action,omitempty"`
	RepositoriesAdded   []*github.Repository `json:"repositories_added,omitempty"`
	RepositoriesRemoved []*github.Repository `json:"repositories_removed,omitempty"`
	Sender              *github.User         `json:"sender,omitempty"`
	HeuprInstallation   *HeuprInstallation   `json:"installation,omitempty"`
}

const secretKey = ""

var Workload = make(chan interface{}, 100)

func collectorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventType := r.Header.Get("X-Github-Event")
		if eventType != "issues" && eventType != "pull_request" && eventType != "installation" && eventType != "installation_repositories" && eventType != "issue_comment" {
			utils.AppLog.Warn("Ignoring event", zap.String("EventType", eventType))
			return
		}
		payload, err := github.ValidatePayload(r, []byte(secretKey))
		if err != nil {
			utils.AppLog.Error("could not validate secret: ", zap.Error(err))
			return
		}
		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			utils.AppLog.Error("could not parse webhook", zap.Error(err))
			return
		}
		switch v := event.(type) {
		case *github.IssuesEvent:
			Workload <- *v
		case *github.PullRequestEvent:
			Workload <- *v
		case *github.IssueCommentEvent:
			Workload <- *v
		case *github.InstallationEvent:
			e := &HeuprInstallationEvent{}
			err := json.Unmarshal(payload, &e)
			if err != nil {
				utils.AppLog.Error("could not parse webhook", zap.Error(err))
				return
			}
			Workload <- *e
		case *github.InstallationRepositoriesEvent:
			e := &HeuprInstallationRepositoriesEvent{}
			err := json.Unmarshal(payload, &e)
			if err != nil {
				utils.AppLog.Error("could not parse webhook", zap.Error(err))
				return
			}
			Workload <- *e
		default:
			utils.AppLog.Error("Unknown", zap.ByteString("GithubEvent", payload))
		}
	})
}
