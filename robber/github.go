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

// getCachedUserOrOrg retrieves cached repos under user or org.
// First tries to read the .git folder under a folder named "name",
// and if that doesn't exist it assumes that the folder is .git folder
func getCachedUserOrOrg(name string) []*string {
	var folderPath string
	repos := []*string{}
	folderPath = filepath.Join(os.TempDir(), "yar", name)
	files, err := ioutil.ReadDir(folderPath)

	if err != nil {
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

// GetUserRepos returns all non forked public repositories for a given user.
func GetUserRepos(m *Middleware, username string) []*string {
	cache := getCachedUserOrOrg(username)
	if !*m.Flags.NoCache && !*m.Flags.NoBare && len(cache) != 0 {
		return cache
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
	cache := getCachedUserOrOrg(orgname)
	if !*m.Flags.NoCache && !*m.Flags.NoBare && len(cache) != 0 {
		return cache
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
	cache := getCachedOrgMembers(orgname)
	if !*m.Flags.NoCache && !*m.Flags.NoBare && len(cache) != 0 {
		return cache
	}

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
	folderPath := filepath.Join(os.TempDir(), "yar", orgname)
	os.MkdirAll(folderPath, 0777)
	err := WriteToFile(filepath.Join(folderPath, "members.txt"), usernames)
	if err != nil {
		m.Logger.LogWarn("Failed to save org members of %s due to: %s\n", orgname, err)
	}
	return usernames
}
