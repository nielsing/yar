package robber

import (
	"sync"

	"github.com/google/go-github/github"
)

// Middleware handles all flags, rules, secrets and logging.
// It essentially holds all values which will be accessed by multiple go routines.
type Middleware struct {
	sync.Mutex
	Logger  *Logger
	Flags   *Flags
	Rules   []*Rule
	Secrets map[string]map[string]bool
	Client  *github.Client
}

// NewMiddleware creates a new Middleware and returns it.
func NewMiddleware() *Middleware {
	m := &Middleware{
		Secrets: make(map[string]map[string]bool),
		Flags:   ParseFlags(),
	}
	m.Logger = NewLogger(*m.Flags.Verbose)
	// If CleanUp flag is given, handle immediately
	if *m.Flags.CleanUp {
		CleanUp(m)
	}
	ParseRegex(m)
	m.Client = github.NewClient(GetAccessToken(m))
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
