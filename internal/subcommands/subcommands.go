package subcommands

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/nielsing/yar/internal/analyzer"
	"github.com/nielsing/yar/internal/fetcher"
	"github.com/nielsing/yar/internal/processor"
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
	numAnalyzers := r.Args.Workers
	if numAnalyzers == 0 {
		numAnalyzers = runtime.NumCPU()
	}
	input := make(chan *processor.DiffObject, 1000)
	var wg sync.WaitGroup
	var analyzers = make([]chan *processor.DiffObject, numAnalyzers, numAnalyzers)

	// Start all analyzers
	for a := 0; a < numAnalyzers; a++ {
		analyzers[a] = analyzer.AnalyzeDiffs(r, input)
	}
	secrets := processor.Multiplex(analyzers...)
	go func() {
		wg.Add(1)
		for value := range secrets {
			// The secrets channel will send secrets instead of diffs once the secret analysis is ready
			r.Logger.LogInfo("Received Diff %s with %d lines\n", *value.Diff, len(strings.Split(*value.Diff, "\n")))
		}
		wg.Done()
	}()

	// Fetcher
	repos := fetcher.RepoGenerator(r, r.Args.Git.Repo)

	// Processor
	for repo := range repos {
		diffs, _ := processor.GetDiffObjects(r, repo, "TODO: FIX!")
		// Send off processed data to analyzers
		for _, diff := range diffs {
			input <- diff
		}
	}
	close(input)
	wg.Wait()
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
