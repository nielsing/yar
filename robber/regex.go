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
	Rules         []*Rule
	FileBlacklist []string `json:"FileBlacklist"`
}

// Rule struct holds a given regex rule with its' reason for matching.
type Rule struct {
	Reason string
	Regex  *regexp.Regexp
	// Noiselevel int
}

// ParseRegex reads regex rules from a given JSON file.
// If no file is given, it reads the default rule file (rules.json).
// Compiles and adds each regex rule to the middleware.
func ParseRegex(m *Middleware) {
	var rules []*Rule
	var values map[string]string

	// Read contents of JSON file
	reader := bufio.NewReader(m.Flags.Rules)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		m.Logger.LogFail("Unable to read file %s: %s", m.Flags.Rules.Name(), err)
	}

	// Parse JSON file and compile regex rules
	json.Unmarshal([]byte(content), &values)
	for key, value := range values {
		regex, err := regexp.Compile(value)
		if err != nil {
			m.Logger.LogFail(regexErrorMessage, key, value, err)
		}
		rule := &Rule{
			Reason: key,
			Regex:  regex,
		}
		rules = append(rules, rule)
	}

	m.Rules = rules
	m.Flags.Rules.Close()
}
