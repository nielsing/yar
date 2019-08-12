package robber

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"sync"
	"time"
)

const (
	debug = iota
	secret
	info
	data
	succ
	warn
	fail
)

const seperator = "--------------------------------------------------------"

var LogColors = map[int]*color.Color{
	debug:  color.New(color.FgBlue),
	secret: color.New(color.FgHiYellow).Add(color.Bold),
	info:   color.New(color.FgHiWhite),
	data:   color.New(color.FgHiBlue),
	succ:   color.New(color.FgGreen),
	warn:   color.New(color.FgRed),
	fail:   color.New(color.FgRed).Add(color.Bold),
}

type Finding struct {
	CommitHash    string
	CommitMessage string
	Committer     string
	DateOfCommit  string
	Email         string
	Reason        string
	Secret        []int
	Diff          string
	RepoName      string
}

type Logger struct {
	sync.Mutex
	Debug bool
}

func NewFinding(reason string, secret []int, commit *object.Commit, reponame string) *Finding {
	finding := &Finding{
		CommitHash:    commit.Hash.String(),
		CommitMessage: commit.Message,
		Committer:     commit.Committer.Name,
		DateOfCommit:  commit.Committer.When.Format(time.RFC1123),
		Email:         commit.Committer.Email,
		Reason:        reason,
		Secret:        secret,
		RepoName:      reponame,
	}
	return finding
}

func (l *Logger) log(level int, format string, a ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if level == debug && l.Debug == false {
		return
	}

	if c, ok := LogColors[level]; ok {
		c.Printf(format, a...)
	} else {
		fmt.Printf(format, a...)
	}

	if level == fail {
		os.Exit(1)
	}
}

func (l *Logger) logSecret(diff string, booty []int, contextNum int) {
	data, _ := LogColors[data]
	secret, _ := LogColors[secret]

	data.Printf("%s", diff[:booty[0]])
	secret.Printf("%s", diff[booty[0]:booty[1]])
	data.Printf("%s\n\n", diff[booty[1]:])
}

func (l *Logger) LogFinding(f *Finding, m *Middleware, diff string) {
	l.Lock()
	defer l.Unlock()
	info, _ := LogColors[info]
	data, _ := LogColors[data]
	secret, _ := LogColors[secret]

	info.Println(seperator)
	info.Printf("Reason: ")
	data.Println(f.Reason)
	info.Printf("Repo name: ")
	data.Println(f.RepoName)
	info.Printf("Committer: ")
	data.Printf("%s (%s)\n", f.Committer, f.Email)
	info.Printf("Commit hash: ")
	data.Println(f.CommitHash)
	info.Printf("Date of commit: ")
	data.Println(f.DateOfCommit)
	info.Printf("Commit message: ")
	data.Println(f.CommitMessage)
	if *m.Flags.NoContext {
		secret.Printf("%s\n\n", diff[f.Secret[0]:f.Secret[1]])
	} else {
		l.logSecret(diff, f.Secret, *m.Flags.Context)
	}
}

func (l *Logger) LogDebug(format string, a ...interface{}) {
	l.log(debug, format, a...)
}

func (l *Logger) LogSecret(format string, a ...interface{}) {
	l.log(secret, format, a...)
}

func (l *Logger) LogInfo(format string, a ...interface{}) {
	l.log(info, format, a...)
}

func (l *Logger) LogSucc(format string, a ...interface{}) {
	l.log(succ, format, a...)
}

func (l *Logger) LogWarn(format string, a ...interface{}) {
	l.log(warn, format, a...)
}

func (l *Logger) LogFail(format string, a ...interface{}) {
	l.log(fail, format, a...)
}
