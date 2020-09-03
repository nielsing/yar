package robber

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
)

const (
	verbose = iota
	secret
	info
	data
	succ
	warn
	fail
)

const seperator = "--------------------------------------------------------"

var validOptions = map[string]int{
	"YAR_COLOR_VERBOSE": verbose,
	"YAR_COLOR_SECRET":  secret,
	"YAR_COLOR_INFO":    info,
	"YAR_COLOR_DATA":    data,
	"YAR_COLOR_SUCC":    succ,
	"YAR_COLOR_WARN":    warn,
	"YAR_COLOR_FAIL":    fail,
}

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
	verbose: color.New(color.FgBlue),
	secret:  color.New(color.FgHiYellow).Add(color.Bold),
	info:    color.New(color.FgHiWhite),
	data:    color.New(color.FgHiBlue),
	succ:    color.New(color.FgGreen),
	warn:    color.New(color.FgRed),
	fail:    color.New(color.FgRed).Add(color.Bold),
}

// TODO: Comment and test
func parseColors(r *Robber, config jsonConfig) {
	for _, colorOption := range config.Colors {
		option, validOption := validOptions[colorOption.Name]
		value, validColor := validColors[colorOption.Value]
		if !validOption || !validColor {
			continue
		}
		logColors[option] = value
	}
}

// Logger handles all logging.
type Logger struct {
	sync.Mutex
}

func newLogger() *Logger {
	return &Logger{}
}

func (l *Logger) log(level int, format string, a ...interface{}) {
	l.Lock()
	defer l.Unlock()

	if c, ok := logColors[level]; ok {
		c.Printf(format, a...)
	} else {
		fmt.Printf(format, a...)
	}

	if level == fail {
		os.Exit(1)
	}
}

// LogVerbose prints to output using 'verbose' colors
func (l *Logger) LogVerbose(format string, a ...interface{}) {
	l.log(verbose, format, a...)
}

// LogSecret prints to output using 'secret' colors
func (l *Logger) LogSecret(format string, a ...interface{}) {
	l.log(secret, format, a...)
}

// LogInfo prints to output using 'info' colors
func (l *Logger) LogInfo(format string, a ...interface{}) {
	l.log(info, "[+] "+format, a...)
}

// LogSucc prints to output using 'succ' colors
func (l *Logger) LogSucc(format string, a ...interface{}) {
	l.log(succ, "[+] "+format, a...)
}

// LogWarn prints to output using 'warn' colors
func (l *Logger) LogWarn(format string, a ...interface{}) {
	l.log(warn, "[-] "+format, a...)
}

// LogFail prints to output using 'fail' colors
func (l *Logger) LogFail(format string, a ...interface{}) {
	l.log(fail, "[!] "+format, a...)
}
