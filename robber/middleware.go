package robber

import (
	"os"
	"regexp"
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
	Blacklist   []*regexp.Regexp
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
	m.Logger = NewLogger(false)
	// If CleanUp flag is given, handle immediately
	if *m.Flags.CleanUp {
		CleanUp(m)
	}
	ParseConfig(m)
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

// Start handles the CLI args and starts yar accordingly.
func (m *Middleware) Start() {
	cpuCount := runtime.NumCPU()
	quit := make(chan bool)
	repoch := make(chan string, cpuCount)
	for proc := 0; proc < cpuCount; proc++ {
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

	// Handle edge case of an org/user containing 0 repos.
	if atomic.LoadInt32(m.RepoCount) == 0 {
		os.Exit(0)
	}

	select {
	case <-quit:
		break
	}
}
