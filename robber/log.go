package robber

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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

var validColors = map[string]*color.Color{
	"black":     color.New(color.FgBlack),
	"blue":      color.New(color.FgBlue),
	"cyan":      color.New(color.FgCyan),
	"green":     color.New(color.FgGreen),
	"magenta":   color.New(color.FgMagenta),
	"red":       color.New(color.FgRed),
	"white":     color.New(color.FgWhite),
	"yellow":    color.New(color.FgYellow),
	"hiBlack":   color.New(color.FgHiBlack),
	"hiBlue":    color.New(color.FgHiBlue),
	"hiCyan":    color.New(color.FgHiCyan),
	"hiGreen":   color.New(color.FgHiGreen),
	"hiMagenta": color.New(color.FgHiMagenta),
	"hiRed":     color.New(color.FgHiRed),
	"hiWhite":   color.New(color.FgHiWhite),
	"hiYellow":  color.New(color.FgHiYellow),
}

// Default colors are set
var logColors = map[int]*color.Color{
	debug:  color.New(color.FgBlue),
	secret: color.New(color.FgHiYellow).Add(color.Bold),
	info:   color.New(color.FgHiWhite),
	data:   color.New(color.FgHiBlue),
	succ:   color.New(color.FgGreen),
	warn:   color.New(color.FgRed),
	fail:   color.New(color.FgRed).Add(color.Bold),
}

// Finding struct contains data of a given secret finding, used for later output of a finding.
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
	Filepath      string
}

// Logger handles all logging to the output.
type Logger struct {
	sync.Mutex
	Debug bool
}

func setColors() {
	colors := GetEnvColors()
	for colorType := debug; colorType <= fail; colorType++ {
		if empty, _ := colors[colorType]; empty == "" {
			continue
		}
		fields := strings.Fields(colors[colorType])
		if val, ok := validColors[fields[0]]; ok {
			if len(fields) > 1 && fields[1] == "bold" {
				logColors[colorType] = val.Add(color.Bold)
				continue
			}
			logColors[colorType] = val
		}
	}
}

// NewLogger sets all colors as specified and returns a new logger.
func NewLogger(debug bool) *Logger {
	setColors()
	return &Logger{
		Debug: debug,
	}
}

// NewFinding simply returns a new finding struct.
func NewFinding(reason string, secret []int, commit *object.Commit, reponame string, filepath string) *Finding {
	finding := &Finding{
		CommitHash:    commit.Hash.String(),
		CommitMessage: commit.Message,
		Committer:     commit.Committer.Name,
		DateOfCommit:  commit.Committer.When.Format(time.RFC1123),
		Email:         commit.Committer.Email,
		Reason:        reason,
		Secret:        secret,
		RepoName:      reponame,
		Filepath:      filepath,
	}
	return finding
}

func (l *Logger) log(level int, format string, a ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if level == debug && l.Debug == false {
		return
	}

	if c, ok := logColors[level]; ok {
		c.Printf(format, a...)
	} else {
		fmt.Printf(format, a...)
	}

	if level == fail {
		os.Exit(1)
	}
}

func (l *Logger) logSecret(diff string, booty []int, contextNum int) {
	data, _ := logColors[data]
	secret, _ := logColors[secret]

	data.Printf("%s", diff[:booty[0]])
	secret.Printf("%s", diff[booty[0]:booty[1]])
	data.Printf("%s\n\n", diff[booty[1]:])
}

// LogFinding is used to output Findings
func (l *Logger) LogFinding(f *Finding, m *Middleware, diff string) {
	l.Lock()
	defer l.Unlock()
	info, _ := logColors[info]
	data, _ := logColors[data]
	secret, _ := logColors[secret]

	info.Println(seperator)
	info.Printf("Reason: ")
	data.Println(f.Reason)
	if f.Filepath != "" {
		info.Printf("Filepath: ")
		data.Println(f.Filepath)
	}
	info.Printf("Repo name: ")
	data.Println(f.RepoName)
	info.Printf("Committer: ")
	data.Printf("%s (%s)\n", f.Committer, f.Email)
	info.Printf("Commit hash: ")
	data.Println(f.CommitHash)
	info.Printf("Date of commit: ")
	data.Println(f.DateOfCommit)
	info.Printf("Commit message: ")
	data.Printf("%s\n\n", strings.Trim(f.CommitMessage, "\n"))
	if *m.Flags.NoContext {
		secret.Printf("%s\n\n", diff[f.Secret[0]:f.Secret[1]])
	} else {
		l.logSecret(diff, f.Secret, *m.Flags.Context)
	}
}

// LogDebug prints to output using 'debug' colors
func (l *Logger) LogDebug(format string, a ...interface{}) {
	l.log(debug, format, a...)
}

// LogSecret prints to output using 'secret' colors
func (l *Logger) LogSecret(format string, a ...interface{}) {
	l.log(secret, format, a...)
}

// LogInfo prints to output using 'info' colors
func (l *Logger) LogInfo(format string, a ...interface{}) {
	l.log(info, format, a...)
}

// LogSucc prints to output using 'succ' colors
func (l *Logger) LogSucc(format string, a ...interface{}) {
	l.log(succ, format, a...)
}

// LogWarn prints to output using 'warn' colors
func (l *Logger) LogWarn(format string, a ...interface{}) {
	l.log(warn, format, a...)
}

// LogFail prints to output using 'fail' colors
func (l *Logger) LogFail(format string, a ...interface{}) {
	l.log(fail, format, a...)
}
