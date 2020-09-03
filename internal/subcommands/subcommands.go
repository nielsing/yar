package subcommands

import (
	"os"
	"path/filepath"

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
	r.Logger.LogFail("Unimplemented!\n")
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
