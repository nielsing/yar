package subcommands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/nielsing/yar/internal/analyzer"
	"github.com/nielsing/yar/internal/robber"
	"github.com/nielsing/yar/internal/utils"
)

// Clear handles the 'clear' subcommand
func Clear(r *robber.Robber) {
	err := os.RemoveAll(filepath.Join(os.TempDir(), "yar", r.Args.Clear.Cache))
	if err != nil {
		r.Logger.LogFail("Failed to clear cache!\n")
	}
	os.Exit(0)
}

// Git handles the 'git' subcommand
func Git(r *robber.Robber) {
	// Boilerplate
	numWorkers := r.Args.Workers
	if numWorkers == 0 {
		numWorkers = runtime.NumCPU()
	}
	input := make(chan string, 100)
	var workers = make([]chan string, numWorkers, numWorkers)

	// Start all workers
	for worker := 0; worker < numWorkers; worker++ {
		workers[worker] = analyzer.AnalyzeRepos(r, input)
	}

	// Fetch all repos
	for _, repo := range r.Args.Git.Repo {
		input <- repo
	}
	close(input)

	c := utils.Multiplex(workers...)

	for value := range c {
		fmt.Println(value)
	}
}

// Github handles the 'github' subcommand
func Github(r *robber.Robber) {
	r.Logger.LogFail("Unimplemented!\n")
}

// Gitlab handles the 'gitlab' subcommand
func Gitlab(r *robber.Robber) {
	r.Logger.LogFail("Unimplemented!\n")
}

// Bitbucket handles the 'bitbucket' subcommand
func Bitbucket(r *robber.Robber) {
	r.Logger.LogFail("Unimplemented!\n")
}
