package fetcher

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nielsing/yar/internal/robber"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func handleGithubError(r *robber.Robber, err error, name string) {
	if err == nil {
		return
	}
	if _, ok := err.(*github.RateLimitError); ok {
		r.Logger.LogWarn("Hit Github rate limit.\n")
		return
	}
	if strings.Contains(err.Error(), "Bad credentials") {
		r.Logger.LogFail("Github token is invalid!\n")
	}
	if strings.Contains(err.Error(), "Not Found") {
		r.Logger.LogFail("%s does not exist.\n", name)
	}
	r.Logger.LogFail("%s\n", err)
}

// getCachedUserOrOrg retrieves cached repos under user or org.
// First tries to read the .git folder under a folder named "name",
// and if that doesn't exist it assumes that the folder is .git folder
func getCachedUserOrOrg(r *robber.Robber, name string) []*string {
	var folderPath string
	repos := []*string{}
	folderPath = filepath.Join(os.TempDir(), "yar", name)
	files, err := ioutil.ReadDir(folderPath)

	if err != nil {
		return repos
	}

	if r.Args.Git.NoBare || r.Args.NoCache {
		os.RemoveAll(filepath.Join(os.TempDir(), "yar", name))
		return repos
	}

	for _, file := range files {
		if file.Name() == "members.txt" {
			continue
		}
		gitFolder := filepath.Join(folderPath, file.Name(), ".git")
		if _, err := os.Stat(gitFolder); err != nil {
			gitFolder = filepath.Dir(gitFolder)
			repos = append(repos, &gitFolder)
		} else {
			repos = append(repos, &gitFolder)
		}
	}
	return repos
}

func getCachedOrgMembers(orgname string) []*string {
	members := []*string{}
	filename := filepath.Join(os.TempDir(), "yar", orgname, "members.txt")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return members
	}

	for _, member := range strings.Split(string(content), "\n") {
		if member != "" {
			members = append(members, &member)
		}
	}
	return members
}

// CreateGithubClient takes a given accesstoken and returns either an authenticated github client
// or, if the `token` is an empty string an unauthenticated github client.
func CreateGithubClient(token string) *github.Client {
	if token == "" {
		return github.NewClient(nil)
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

// GetUserRepos returns all repositories of a given user.
func GetUserRepos(r *robber.Robber, c *github.Client, username string) []*string {
	cache := getCachedUserOrOrg(r, username)
	if !r.Args.NoCache && !r.Args.Git.NoBare && len(cache) != 0 {
		return cache
	}

	cloneURLs := []*string{}
	opt := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 100}}
	for {
		repos, resp, err := c.Repositories.List(context.Background(), username, opt)
		handleGithubError(r, err, username)

		for _, repo := range repos {
			if *repo.Fork && r.Args.Github.Forks {
				continue
			}
			cloneURLs = append(cloneURLs, repo.CloneURL)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return cloneURLs
}

// GetOrgRepos returns all repositories of a given organization.
func GetOrgRepos(r *robber.Robber, c *github.Client, orgname string) []*string {
	cache := getCachedUserOrOrg(r, orgname)
	if !r.Args.NoCache && !r.Args.Git.NoBare && len(cache) != 0 {
		return cache
	}

	cloneURLs := []*string{}
	opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 100}}
	for {
		repos, resp, err := c.Repositories.ListByOrg(context.Background(), orgname, opt)
		handleGithubError(r, err, orgname)

		for _, repo := range repos {
			if *repo.Fork && r.Args.Github.Forks {
				continue
			}
			cloneURLs = append(cloneURLs, repo.CloneURL)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return cloneURLs
}

// GetOrgMembers returns all members of a given organization.
func GetOrgMembers(r *robber.Robber, c *github.Client, orgname string) []*string {
	cache := getCachedOrgMembers(orgname)
	if !r.Args.NoCache && !r.Args.Git.NoBare && len(cache) != 0 {
		return cache
	}

	usernames := []*string{}
	opt := &github.ListMembersOptions{ListOptions: github.ListOptions{PerPage: 100}}
	for {
		members, resp, err := c.Organizations.ListMembers(context.Background(), orgname, opt)
		handleGithubError(r, err, orgname)

		for _, member := range members {
			usernames = append(usernames, member.Login)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	folderPath := filepath.Join(os.TempDir(), "yar", orgname)
	os.MkdirAll(folderPath, 0777)
	err := WriteToFile(filepath.Join(folderPath, "members.txt"), usernames)
	if err != nil {
		r.Logger.LogWarn("Failed to save org members of %s due to: %s\n", orgname, err)
	}
	return usernames
}
