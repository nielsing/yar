package robber

import (
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"strings"
)

const (
	B64chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	Hexchars = "1234567890abcdefABCDEF"
)

// AnalyzeEntropyDiff breaks a given diff into words and finds valid base64 and hex
// strings within a word and finally runs an entropy check on the valid string.
// Code taken from https://github.com/dxa4481/truffleHog
func AnalyzeEntropyDiff(m *Middleware, commit *object.Commit, diff string, reponame string, filepath string) {
	words := strings.Fields(diff)
	for _, word := range words {
		b64strings := FindValidStrings(word, B64chars)
		hexstrings := FindValidStrings(word, Hexchars)
		PrintEntropyFinding(b64strings, m, diff, reponame, commit, 4.5, filepath)
		PrintEntropyFinding(hexstrings, m, diff, reponame, commit, 3, filepath)
	}
}

// AnalyzeRegexDiff runs line by line on a given diff and runs each given regex rule on the line.
func AnalyzeRegexDiff(m *Middleware, commit *object.Commit, diff string, reponame string, filepath string) {
	lines := strings.Split(diff, "\n")
	numOfLines := len(lines)

	for lineNum, line := range lines {
		for _, rule := range m.Rules {
			if found := rule.Regex.Match([]byte(line)); found {
				start, end := Max(0, lineNum-*m.Flags.Context), Min(numOfLines, lineNum+*m.Flags.Context+1)
				context := lines[start:end]
				newDiff := strings.Join(context, "\n")
				secret := rule.Regex.FindIndex([]byte(newDiff))

				newSecret := false
				secretString := newDiff[secret[0]:secret[1]]
				if !m.SecretExists(reponame, secretString) {
					m.AddSecret(reponame, secretString)
					newSecret = true
				}
				if newSecret {
					finding := NewFinding(rule.Reason, secret, commit, reponame, filepath)
					m.Logger.LogFinding(finding, m, newDiff)
					break
				}
			}
		}
	}
}

// AnalyzeRepo opens a given repository and extracts all diffs from it for later analysis.
func AnalyzeRepo(m *Middleware, reponame string) {
	repo, err := OpenRepo(reponame)
	if err != nil {
		if err == transport.ErrEmptyRemoteRepository {
			m.Logger.LogWarn("%s is empty\n", reponame)
			return
		}
		m.Logger.LogFail("Unable to open repo %s: %s\n", reponame, err)
	}

	commits, err := GetCommits(repo)
	if err != nil {
		m.Logger.LogWarn("Unable to fetch commits for %s: %s\n", reponame, err)
	}

	// Get all changes in correct order of commit history
	for index := range commits {
		commit := commits[len(commits)-index-1]
		changes, err := GetCommitChanges(commit)
		if err != nil {
			m.Logger.LogWarn("Unable to get commit changes for hash %s: %s\n", commit.Hash, err)
			continue
		}

		for _, change := range changes {
			diffs, filepath, err := GetDiffs(change)
			if err != nil {
				m.Logger.LogWarn("Unable to get diffs of %s: %s\n", change, err)
				break
			}
			for _, diff := range diffs {
				if *m.Flags.Both {
					AnalyzeRegexDiff(m, commit, diff, reponame, filepath)
					AnalyzeEntropyDiff(m, commit, diff, reponame, filepath)
				} else if *m.Flags.Entropy {
					AnalyzeEntropyDiff(m, commit, diff, reponame, filepath)
				} else {
					AnalyzeRegexDiff(m, commit, diff, reponame, filepath)
				}
			}
		}
	}
}

// AnalyzeUser simply sends a GET request on githubs API for a given username
// and starts and analysis of each of the user's repositories.
func AnalyzeUser(m *Middleware, username string) {
	repos := GetUserRepos(m, username)
	for _, repo := range repos {
		AnalyzeRepo(m, *repo)
	}
}

// AnalyzeOrg simply sends two GET requests to githubs API, one for a given organizations
// repositories and one for its' members.
func AnalyzeOrg(m *Middleware, orgname string) {
	repos := GetOrgRepos(m, orgname)
	members := GetOrgMembers(m, orgname)
	for _, repo := range repos {
		AnalyzeRepo(m, *repo)
	}
	for _, member := range members {
		AnalyzeUser(m, *member)
	}
}
