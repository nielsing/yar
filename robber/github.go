package robber

import (
	"context"
	"github.com/google/go-github/github"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func handleGithubError(m *Middleware, err error, name string) {
	if err == nil {
		return
	}
	if _, ok := err.(*github.RateLimitError); ok {
		m.Logger.LogWarn("Hit Github rate limit.\n")
		return
	}
	if strings.Contains(err.Error(), "Bad credentials") {
		m.Logger.LogFail("Github token is invalid!\n")
	}
	if strings.Contains(err.Error(), "Not Found") {
		m.Logger.LogFail("%s does not exist.\n", name)
	}
	m.Logger.LogFail("%s\n", err)
}

func getCachedUserOrOrg(name string) []*string {
	repos := []*string{}
	folderPath := filepath.Join(os.TempDir(), "yar", name)
	files, err := ioutil.ReadDir(folderPath)

	if err != nil {
		return []*string{}
	}

	for _, file := range files {
		gitFolder := filepath.Join(folderPath, file.Name())
		repos = append(repos, &gitFolder)
	}
	return repos
}

// GetUserRepos returns all non forked public repositories for a given user.
func GetUserRepos(m *Middleware, username string) []*string {
	cloneUrls := getCachedUserOrOrg(username)
	if len(cloneUrls) != 0 {
		return cloneUrls
	}

	cloneURLs := []*string{}
	opt := &github.RepositoryListOptions{Type: "public", ListOptions: github.ListOptions{PerPage: 100}}
	for {
		repos, resp, err := m.Client.Repositories.List(context.Background(), username, opt)
		handleGithubError(m, err, username)

		for _, repo := range repos {
			if *repo.Fork && !*m.Flags.Forks {
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
func GetOrgRepos(m *Middleware, orgname string) []*string {
	cloneUrls := getCachedUserOrOrg(orgname)
	if len(cloneUrls) != 0 {
		return cloneUrls
	}

	cloneURLs := []*string{}
	opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 100}}
	for {
		repos, resp, err := m.Client.Repositories.ListByOrg(context.Background(), orgname, opt)
		handleGithubError(m, err, orgname)

		for _, repo := range repos {
			if *repo.Fork && !*m.Flags.Forks {
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
func GetOrgMembers(m *Middleware, orgname string) []*string {
	usernames := []*string{}
	opt := &github.ListMembersOptions{ListOptions: github.ListOptions{PerPage: 100}}
	for {
		members, resp, err := m.Client.Organizations.ListMembers(context.Background(), orgname, opt)
		handleGithubError(m, err, orgname)

		for _, member := range members {
			usernames = append(usernames, member.Login)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return usernames
}
