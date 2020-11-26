package main

import (
	"fmt"
	"github.com/nielsing/yar/internal/robber"
	"github.com/nielsing/yar/internal/subcommands"
)

func main() {
	r := robber.NewRobber()
	fmt.Println(r.Args.Workers)
	if r.Args.Clear != nil {
		subcommands.Clear(r)
	} else if r.Args.Git != nil {
		subcommands.Git(r)
	} else if r.Args.Github != nil {
		subcommands.Github(r)
	} else if r.Args.Gitlab != nil {
		subcommands.Gitlab(r)
	} else if r.Args.Bitbucket != nil {
		subcommands.Bitbucket(r)
	}
}
