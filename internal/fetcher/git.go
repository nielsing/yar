package fetcher

import (
	"os"

	"github.com/nielsing/yar/internal/robber"
	"github.com/nielsing/yar/internal/utils"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// getCloneOptions returns either an authenticated clone of a repo or an
// anonymous clone of a repo based on whether an AccessToken was given or not.
func getCloneOptions(depth int, username, secret, url string) *git.CloneOptions {
	if username != "" {
		return &git.CloneOptions{
			URL:   url,
			Depth: depth + 1, // There is an off by one error in Depth field.
			Auth: &http.BasicAuth{
				Username: username,
				Password: secret,
			},
		}
	}
	if secret != "" {
		return &git.CloneOptions{
			URL:   url,
			Depth: depth + 1, // There is an off by one error in Depth field.
			Auth: &http.BasicAuth{
				Username: "NotEmpty", // https://pkg.go.dev/github.com/go-git/go-git/v5?tab=doc#PlainClone
				Password: secret,
			},
		}
	}
	return &git.CloneOptions{
		URL:   url,
		Depth: depth + 1, // There is an off by one error in Depth field.
	}
}

// cloneRepo clones the given URL to the given folder.
func cloneRepo(r *robber.Robber, url, cloneFolder string) (*git.Repository, error) {
	opt := getCloneOptions(r.Args.Git.Depth, "", "", url)
	repo, err := git.PlainClone(cloneFolder, !r.Args.Git.NoBare, opt)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// OpenRepo clones or opens a given repo based on whether the repo already exists on the system or
// not.
func OpenRepo(r *robber.Robber, location string) (*git.Repository, error) {
	dir, exists := utils.GetCacheLocation(r, location)
	if !r.Args.NoCache && !r.Args.Git.NoBare && exists {
		repo, err := git.PlainOpen(dir)
		if err != nil {
			return nil, err
		}
		return repo, nil
	}

	if exists && (r.Args.Git.NoBare || r.Args.NoCache) {
		os.RemoveAll(dir)
	}
	repo, err := cloneRepo(r, location, dir)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
