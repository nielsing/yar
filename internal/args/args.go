package args

import (
    arg "github.com/alexflint/go-arg"
)

type CleanCmd struct {
    Cache string `arg:"positional"`
}

type GithubCmd struct {
    Org string `arg:"-o"`
    User string `arg:"-u"`
    Repo string `arg:"-r"`
    Forks bool `arg:"-f"`
    Depth int `arg:"-d"`
    NoBare bool `arg:"--no-bare"`
    IncludeMembers bool `arg:"--include-members"`
}

// TODO: Implemented!
type GitlabCmd struct {}

// TODO: Implemented!
type BitbucketCmd struct {}

type Args struct {
    // General flags
    Both bool `arg:"-b"`
    Save bool `arg:"-s"`
    Noise string `arg:"-n"`
    Config string `arg:"-C"`
    Context string `arg:"-c"`
    Entropy string `arg:"-e"`
    NoCache bool `arg:"--no-cache"`
    NoContext bool `arg:"--no-context"`

    // Commands
    Clean *CleanCmd `arg:"subcommand:clean"`
    Github *GithubCmd `arg:"subcommand:github" help:"Plunder github"`
    Gitlab *GitlabCmd `arg:"subcommand:gitlab" help:"Unimplemented!"`
    Bitbucket *BitbucketCmd `arg:"subcommand:github" help:"Unimplemented!"`
}

func ParseArgs() Args {
    parsedArgs := Args{}
    arg.MustParse(&parsedArgs)
    return parsedArgs
}
