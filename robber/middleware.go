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
	Findings    []*Finding
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

// Append appends finding to Middlewares Findings array if save mode is enabled.
func (m *Middleware) Append(finding *Finding) {
	if *m.Flags.SavePresent {
		m.Findings = append(m.Findings, finding)
	}
}

// Start handles the CLI args and starts yar accordingly.
func (m *Middleware) Start(kill chan bool, finished chan<- bool, cleanup <-chan bool) {
	wg := new(sync.WaitGroup)
	cpuCount := runtime.NumCPU()

	quit := make(chan bool, cpuCount) // Channel is buffered to ensure deadlock doesn't happen
	done := make(chan bool)
	repoch := make(chan string, cpuCount)
	// Start all workers
	for proc := 1; proc <= cpuCount; proc++ {
		wg.Add(1)
		go AnalyzeRepo(m, proc, repoch, quit, done, wg)
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

	// Clean up all workers
	select {
	case <-quit:
		for proc := 0; proc < cpuCount; proc++ {
			done <- true
		}
	case <-kill:
		for proc := 0; proc < cpuCount; proc++ {
			done <- true
		}
		finished <- true
		<-cleanup
	}
	wg.Wait()
}
