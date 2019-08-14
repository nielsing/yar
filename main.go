package main

import (
	"github.com/Furduhlutur/yar/robber"
)

// Entry function handles CLI arguments and executes accordingly
func main() {
	m := robber.NewMiddleware()
	if *m.Flags.Org != "" {
		robber.AnalyzeOrg(m, *m.Flags.Org)
	}
	if *m.Flags.User != "" {
		robber.AnalyzeUser(m, *m.Flags.User)
	}
	if *m.Flags.Repo != "" {
		robber.AnalyzeRepo(m, *m.Flags.Repo)
	}
	robber.CleanUp(m)
}
