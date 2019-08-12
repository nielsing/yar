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

// Grunt function for analyzing a diff using Shannon's Entropy
// Code taken from https://github.com/dxa4481/truffleHog
func AnalyzeEntropyDiff(m *Middleware, commit *object.Commit, diff string, reponame string) {
	words := strings.Fields(diff)
	for _, word := range words {
		b64strings := FindValidStrings(word, B64chars)
		hexstrings := FindValidStrings(word, Hexchars)
		PrintEntropyFinding(b64strings, m, diff, reponame, commit, 4.5)
		PrintEntropyFinding(hexstrings, m, diff, reponame, commit, 3)
	}
}

// Grunt function for analyzing a diff using regex rules
func AnalyzeRegexDiff(m *Middleware, commit *object.Commit, diff string, reponame string) {
	lines := strings.Split(diff, "\n")
	numOfLines := len(lines)

	for lineNum, line := range lines {
		for _, rule := range m.Rules {
			if found := rule.Regex.Match([]byte(line)); found {
				context := []string{lines[lineNum]}
				for i := 1; i <= *m.Flags.Context; i++ {
					if lineNum-i >= 0 {
						context = append([]string{lines[lineNum-i]}, context...)
					}
					if lineNum+i < numOfLines {
						context = append(context, lines[lineNum+i])
					}
				}
				newDiff := strings.Join(context, "\n")
				secret := rule.Regex.FindIndex([]byte(newDiff))

				newSecret := false
				secretString := newDiff[secret[0]:secret[1]]
				if !m.SecretExists(reponame, secretString) {
					m.AddSecret(reponame, secretString)
					newSecret = true
				}
				if newSecret {
					finding := NewFinding(rule.Reason, secret, commit, reponame)
					m.Logger.LogFinding(finding, m, newDiff)
				}
			}
		}
	}
}

// Grunt function for analyzing a repo
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
			diffs, err := GetDiffs(change)
			if err != nil {
				m.Logger.LogWarn("Unable to get diffs of %s: %s\n", change, err)
				break
			}
			for _, diff := range diffs {
				if *m.Flags.Both {
					AnalyzeRegexDiff(m, commit, diff, reponame)
					AnalyzeEntropyDiff(m, commit, diff, reponame)
				} else if *m.Flags.Entropy {
					AnalyzeEntropyDiff(m, commit, diff, reponame)
				} else {
					AnalyzeRegexDiff(m, commit, diff, reponame)
				}
			}
		}
	}
}

// Grunt function for analyzing a user
func AnalyzeUser(m *Middleware, username string) {
	repos := GetUserRepos(m, username)
	for _, repo := range repos {
		AnalyzeRepo(m, *repo)
	}
}

// Grunt function for analyzing an organization
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
