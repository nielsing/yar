package subcommands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/nielsing/yar/internal/analyzer"
	"github.com/nielsing/yar/internal/robber"
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
	var workers = make([]<-chan string, numWorkers, numWorkers)

	// Start all workers
	for worker := 0; worker < numWorkers; worker++ {
		workers[worker] = analyzer.AnalyzeRepos(r, input)
	}

	// Fetch all repos
	for _, repo := range r.Args.Git.Repo {
		input <- repo
	}
	close(input)

	c := fanIn(workers...)

	for value := range c {
		fmt.Println(value)
	}
}

func fanIn(workers ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(workers))

	for _, c := range workers {
		go func(c <-chan string) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
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
