package robber

import (
	"context"
	"github.com/google/go-github/github"
	"strings"
)

func handleGithubError(m *Middleware, err error) {
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
	m.Logger.LogFail("%s\n", err)
}

// GetUserRepos returns all non forked public repositories for a given user.
func GetUserRepos(m *Middleware, username string) []*string {
	cloneURLs := []*string{}
	opt := &github.RepositoryListOptions{Type: "public", ListOptions: github.ListOptions{PerPage: 100}}
	for {
		repos, resp, err := m.Client.Repositories.List(context.Background(), username, opt)
		handleGithubError(m, err)

		for _, repo := range repos {
			if !*repo.Fork {
				cloneURLs = append(cloneURLs, repo.CloneURL)
			}
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
	cloneURLs := []*string{}
	opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 100}}
	for {
		repos, resp, err := m.Client.Repositories.ListByOrg(context.Background(), orgname, opt)
		handleGithubError(m, err)

		for _, repo := range repos {
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
		handleGithubError(m, err)

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
