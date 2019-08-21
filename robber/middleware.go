package robber

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/google/go-github/github"
)

// Middleware handles all flags, rules, secrets and logging.
// It essentially holds all values which will be accessed by multiple go routines.
type Middleware struct {
	sync.Mutex
	Logger      *Logger
	Flags       *Flags
	Rules       []*Rule
	Secrets     map[string]map[string]bool
	Client      *github.Client
	AccessToken string
	RepoCount   *int32
}

// NewMiddleware creates a new Middleware and returns it.
func NewMiddleware() *Middleware {
	m := &Middleware{
		Secrets:   make(map[string]map[string]bool),
		Flags:     ParseFlags(),
		RepoCount: new(int32),
	}
	m.Logger = NewLogger(*m.Flags.Verbose)
	// If CleanUp flag is given, handle immediately
	if *m.Flags.CleanUp {
		CleanUp(m)
	}
	ParseRegex(m)
	accessToken, client := GetAccessToken(m)
	m.AccessToken = accessToken
	m.Client = github.NewClient(client)
	return m
}

// AddSecret adds a new secret for a given repo.
func (m *Middleware) AddSecret(reponame string, secret string) {
	m.Lock()
	defer m.Unlock()
	if m.Secrets[reponame] == nil {
		m.Secrets[reponame] = make(map[string]bool)
	}
	m.Secrets[reponame][secret] = true
}

// SecretExists checks to see whether a given secret string has been noticed before or not.
func (m *Middleware) SecretExists(reponame string, secret string) bool {
	m.Lock()
	defer m.Unlock()
	return m.Secrets[reponame][secret]
}

// Start starts yar
func (m *Middleware) Start() {
	quit := make(chan bool)
	repoch := make(chan string, runtime.NumCPU())
	for proc := 0; proc < runtime.NumCPU(); proc++ {
		go AnalyzeRepo(m, repoch, quit)
	}

	if *m.Flags.Org != "" {
		AnalyzeOrg(m, *m.Flags.Org, repoch)
	}
	if *m.Flags.User != "" {
		AnalyzeUser(m, *m.Flags.User, repoch)
	}
	if *m.Flags.Repo != "" {
		atomic.AddInt32(m.RepoCount, 1)
		repoch <- *m.Flags.Repo
	}

	select {
	case <-quit:
		break
	}
}
