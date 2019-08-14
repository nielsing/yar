package robber

import (
	"github.com/google/go-github/github"
	"sync"
)

type Middleware struct {
	sync.Mutex
	Logger  *Logger
	Flags   *Flags
	Rules   []*Rule
	Secrets map[string]map[string]bool
	Client  *github.Client
	Repos   []string
	Users   []string
}

func NewMiddleware() *Middleware {
	m := &Middleware{
		Secrets: make(map[string]map[string]bool),
		Flags:   ParseFlags(),
		Repos:   []string{""},
		Users:   []string{""},
	}
	// If CleanUp flag is given, handle immediately
	if *m.Flags.CleanUp {
		CleanUp()
	}
	m.Logger = NewLogger(*m.Flags.Debug)
	ParseRegex(m)
	m.Client = github.NewClient(GetAccessToken(m))
	return m
}

func (m *Middleware) AddRepo(reponame string) {
	m.Lock()
	defer m.Unlock()
	m.Repos = append(m.Repos, reponame)
}

func (m *Middleware) AddUser(username string) {
	m.Lock()
	defer m.Unlock()
	m.Users = append(m.Users, username)
}

func (m *Middleware) RemoveRepo() string {
	m.Lock()
	defer m.Unlock()
	var reponame string
	reponame, m.Repos = m.Repos[len(m.Repos)-1], m.Repos[:len(m.Repos)-1]
	if reponame == "" {
		m.Repos = append(m.Repos, "")
	}
	return reponame
}

func (m *Middleware) RemoveUser() string {
	m.Lock()
	defer m.Unlock()
	var username string
	username, m.Users = m.Users[len(m.Users)-1], m.Users[:len(m.Users)-1]
	if username == "" {
		m.Users = append(m.Users, "")
	}
	return username
}

func (m *Middleware) AddSecret(reponame string, secret string) {
	m.Lock()
	defer m.Unlock()
	if m.Secrets[reponame] == nil {
		m.Secrets[reponame] = make(map[string]bool)
	}
	m.Secrets[reponame][secret] = true
}

func (m *Middleware) SecretExists(reponame string, secret string) bool {
	m.Lock()
	defer m.Unlock()
	return m.Secrets[reponame][secret]
}

func (m *Middleware) Start() {}
