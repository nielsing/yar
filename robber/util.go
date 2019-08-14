package robber

import (
	"context"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/oauth2"
)

const (
	envTokenVariable = "YAR_GITHUB_TOKEN"
)

// CleanUp deletes all temp directories which were created for cloning of repositories.
func CleanUp(m *Middleware) {
	files, err := ioutil.ReadDir(os.TempDir())
	if err != nil {
		m.Logger.LogFail("Something extremely bad is going on!\n")
		os.Exit(1)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "yar") {
			err := os.RemoveAll(path.Join(os.TempDir(), file.Name()))
			if err != nil {
				m.Logger.LogInfo("Unable to remove %s\n", file.Name())
			}
		}
	}
	os.Exit(0)
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
			secretString := context[indexes[0]:indexes[1]]
			if !m.SecretExists(*diffObject.Reponame, secretString) {
				m.AddSecret(*diffObject.Reponame, secretString)
				m.Logger.LogFinding(NewFinding("Entropy Check", indexes, diffObject), m, context)
			}
		}
	}
}

// GetAccessToken retreives access token from env variables and returns an oauth2 client.
func GetAccessToken(m *Middleware) *http.Client {
	accessToken := os.Getenv(envTokenVariable)
	if accessToken == "" {
		m.Logger.LogWarn("No access token found for GitHub, consider adding it by running 'export YAR_GITHUB_TOKEN=YOUR_TOKEN'.\n")
		return nil
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return tc
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
