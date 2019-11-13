package robber

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"regexp"
)

const (
	regexErrorMessage = `Invalid regex rule in file!
    Key: %s
    Rule: %s
    Error: %s

Read here for more information: https://golang.org/pkg/regexp/syntax/
`
)

// Config struct holds all config from the given JSON file.
type Config struct {
	Rules []struct {
		Reason string `json:"Reason"`
		Rule   string `json:"Rule"`
		Noise  int    `json:"Noise"`
	} `json:"Rules"`
	FileBlacklist []string `json:"FileBlacklist"`
}

// Rule struct holds a given regex rule with its' reason for matching.
type Rule struct {
	Reason string
	Regex  *regexp.Regexp
}

// ParseConfig parses a given config file, if there was none given
// it will parse the default config file.
//
// ParseConfig first parses all rules in the config file below a given noiselevel
// the default max noiselevel being 3.
// Then it parses all regex rules for the file blacklist.
func ParseConfig(m *Middleware) {
	var config Config
	var rules []*Rule
	var blacklist []*regexp.Regexp

	// Read contents of JSON file
	reader := bufio.NewReader(m.Flags.Config)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		m.Logger.LogFail("Unable to read file %s: %s", m.Flags.Config.Name(), err)
	}

	// Parse JSON file and compile regex rules
	json.Unmarshal([]byte(content), &config)
	for _, rule := range config.Rules {
		if rule.Noise > m.Flags.NoiseLevel.Upper || rule.Noise < m.Flags.NoiseLevel.Lower {
			continue
		}
		regex, err := regexp.Compile(rule.Rule)
		if err != nil {
			m.Logger.LogFail(regexErrorMessage, rule.Reason, rule.Rule, err)
		}
		rule := &Rule{
			Reason: rule.Reason,
			Regex:  regex,
		}
		rules = append(rules, rule)
	}

	for _, fileRule := range config.FileBlacklist {
		regex, err := regexp.Compile(fileRule)
		if err != nil {
			m.Logger.LogFail(regexErrorMessage, "File blacklist", fileRule, err)
		}
		blacklist = append(blacklist, regex)
	}
	m.Rules = rules
	m.Blacklist = blacklist
	m.Flags.Config.Close()
}
