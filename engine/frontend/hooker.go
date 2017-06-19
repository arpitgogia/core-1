package frontend

import (
	"context"
	"strconv"

	"github.com/google/go-github/github"
)

const secretKey = "figrin-dan-and-the-modal-nodes"

func (fs *FrontendServer) NewHook(repo *github.Repository, client *github.Client) error {
	if check, err := fs.hookExists(repo, client); check {
		return err
	}
	name := *repo.Name
	owner := *repo.Owner.Login
	// TODO: This URL will change to a config parameter.
	url := "http://00ad0ac7.ngrok.io/hook"
	hook, _, err := client.Repositories.CreateHook(context.Background(), owner, name, &github.Hook{
		Name:   github.String("web"),
		Events: []string{"issues", "repository"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          url,
			"secret":       secretKey,
			"content_type": "json",
			"insecure_ssl": false,
		},
	})
	if err != nil {
		return err
	}
	if err = fs.Database.Store("hook", *repo.ID, []byte(strconv.Itoa(*hook.ID))); err != nil {
		return err
	}
	return nil
}

func (fs *FrontendServer) hookExists(repo *github.Repository, client *github.Client) (bool, error) {
	name, owner := "", ""
	if repo.Name != nil && repo.Owner.Login != nil {
		name = *repo.Name
		owner = *repo.Owner.Login
	}
	hook, err := fs.Database.Retrieve("hook", *repo.ID)
	if err != nil {
		return false, err
	}

	hookID, err := strconv.Atoi(string(hook))
	if err != nil {
		return false, err
	}
	_, _, err = client.Repositories.GetHook(context.Background(), owner, name, hookID)
	if err != nil {
		return false, err
	}
	return true, nil
}
