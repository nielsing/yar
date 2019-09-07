package robber

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"
)

const (
	envTokenVariable = "YAR_GITHUB_TOKEN"
)

type jsonFinding []struct {
	Reason        string `json:"Reason"`
	Filepath      string `json:"Filepath"`
	RepoName      string `json:"RepoName"`
	Commiter      string `json:"Commiter"`
	CommitHash    string `json:"CommitHash"`
	DateOfCommit  string `json:"DateOfCommit"`
	CommitMessage string `json:"CommitMessage"`
	Link          string `json:"Link"`
	Secret        string `json:"Secret"`
}

// CleanUp deletes all temp directories which were created for cloning of repositories.
func CleanUp(m *Middleware) {
	err := os.RemoveAll(path.Join(os.TempDir(), "yar"))
	if err != nil {
		m.Logger.LogWarn("Unable to remove the cache folder!")
	}
	os.Exit(0)
}

// HandleSigInt captures the SIGINT signal and removes the cache folder.
// This is done to avoid nil pointers for future runs of yar.
func HandleSigInt(m *Middleware, sigc chan os.Signal, kill chan<- bool, finished <-chan bool, cleanup chan<- bool) {
	<-sigc
	m.Logger.LogInfo("Signal reveived, killing threads!\n")
	kill <- true
	<-finished
	CleanUp(m)
	cleanup <- true
}

// GetDir returns the respective directory of a given cloneurl and whether it exists.
func GetDir(cloneurl string) (string, bool) {
	if _, err := os.Stat(cloneurl); err == nil {
		return cloneurl, true
	}
	names := strings.Split(cloneurl, "/")
	parentFolder := names[len(names)-2]
	childFolder := strings.Replace(names[len(names)-1], ".git", "", -1)
	dir := filepath.Join(os.TempDir(), "yar", parentFolder, childFolder)
	_, err := os.Stat(dir)
	return dir, os.IsNotExist(err)
}

// FindValidStrings finds parts of a word which are valid in respect
// to a given charset
func FindValidStrings(word string, charSet string) []string {
	count := 0
	value := ""
	values := []string{}
	for _, char := range word {
		if strings.Contains(charSet, string(char)) {
			value += string(char)
			count++
		} else {
			if count > 15 {
				values = append(values, value)
			}
			value, count = "", 0
		}
	}
	if count > 15 {
		values = append(values, value)
	}
	return values
}

// EntropyCheck runs Shannon's Entropy on a given word
// H(X) = - \sigma{i=1}{n} P(x_i) log_bP(x_i)
// P(X = x) = P({s \in S: X(s) = x})
func EntropyCheck(data string, values string) float64 {
	if data == "" {
		return 0.0
	}

	var entropy float64
	for _, letter := range values {
		pX := float64(strings.Count(data, string(letter))) / float64(len(data))
		if pX > 0 {
			entropy += -(pX * math.Log2(pX))
		}
	}
	return entropy
}

// FindContext finds context lines of an entropy finding.
func FindContext(m *Middleware, diff string, secret string) (string, []int) {
	lines := strings.Split(diff, "\n")
	numOfLines := len(lines)

	for lineNum, line := range lines {
		if strings.Contains(line, secret) {
			start, end := Max(0, lineNum-*m.Flags.Context), Min(numOfLines, lineNum+*m.Flags.Context+1)
			context := lines[start:end]
			newDiff := strings.Join(context, "\n")
			index := strings.Index(newDiff, secret)
			return newDiff, []int{index, index + len(secret)}
		}
	}
	return "", nil
}

// PrintEntropyFinding checks for a given validString set whether the threshold is broken and if it is
// finds the context around the secret of the diff and prints it along with the secret.
func PrintEntropyFinding(validStrings []string, m *Middleware, diffObject *DiffObject, threshold float64) {
	for _, validString := range validStrings {
		entropy := EntropyCheck(validString, B64chars)
		if entropy > threshold {
			context, indexes := FindContext(m, *diffObject.Diff, validString)
			diffObject.Diff = &context
			secretString := context[indexes[0]:indexes[1]]
			if !m.SecretExists(*diffObject.Reponame, secretString) {
				m.AddSecret(*diffObject.Reponame, secretString)
				finding := NewFinding("Entropy Check", indexes, diffObject)
				m.Findings = append(m.Findings, finding)
				m.Logger.LogFinding(finding, m)
			}
		}
	}
}

// SaveFindings saves all findings to a JSON file named findings.json
func SaveFindings(findings []*Finding) {
	var savedFindings jsonFinding
	for _, finding := range findings {
		savedFindings = append(savedFindings, jsonFinding{{
			Reason:        finding.Reason,
			Filepath:      finding.Filepath,
			RepoName:      finding.RepoName,
			Commiter:      finding.Committer,
			CommitHash:    finding.CommitHash,
			DateOfCommit:  finding.DateOfCommit,
			CommitMessage: finding.CommitMessage,
			Link:          strings.Join([]string{finding.RepoName, "commit", finding.CommitHash}, "/"),
			Secret:        finding.Diff[finding.Secret[0]:finding.Secret[1]],
		}}...)
	}
	content, _ := json.MarshalIndent(savedFindings, "", "  ")
	_ = ioutil.WriteFile("findings.json", content, 0644)
}

// GetAccessToken retreives access token from env variables and returns an oauth2 client.
func GetAccessToken(m *Middleware) (string, *http.Client) {
	accessToken := os.Getenv(envTokenVariable)
	if accessToken == "" {
		return "", nil
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return accessToken, tc
}

// GetGoPath returns user's GOPATH env variable.
func GetGoPath() string {
	gopath := os.Getenv("GOPATH")
	return gopath
}

// GetEnvColors retreives color settings from env variables and returns them.
func GetEnvColors() map[int]string {
	colors := map[int]string{}
	values := []string{"VERBOSE", "SECRET", "INFO", "DATA", "SUCC", "WARN", "FAIL"}
	baseValue := "YAR_COLOR_"

	for index, value := range values {
		colors[index] = os.Getenv(baseValue + value)
	}
	return colors
}

// Max returns the larger of two given ints
func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// Min returns the smaller of two given ints
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
