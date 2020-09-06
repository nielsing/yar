package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/nielsing/yar/internal/robber"

	"github.com/whilp/git-urls"
)

// Min returns the minimum of a and b
func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// Max returns the maximum of a and b
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// BlacklistedFile returns whether a given filename is blacklisted or not.
func BlacklistedFile(r *robber.Robber, filename string) bool {
	for _, rule := range r.Config.Blacklist {
		if rule.Match([]byte(filename)) {
			return true
		}
	}
	return false
}

func getCacheHelper(r *robber.Robber, location string) (string, string) {
	end := ""
	base := ""
	if r.Args.Git != nil {
		base = "Git"
		website := ""
		url, err := giturls.Parse(location)
		if err != nil {
			website = "Default"
		} else {
			website = url.Hostname()
		}
		gitFolder := strings.Replace(filepath.Base(location), ".git", "", -1)
		filepath.Join(end, website, gitFolder)
	}
	if r.Args.Github != nil {
		base = "Github"
		user := filepath.Base(filepath.Dir(location))
		repo := strings.Replace(filepath.Base(location), ".git", "", -1)
		filepath.Join(end, user, repo)
	}
	if r.Args.Gitlab != nil {
		base = "Gitlab"
		end = "Unimplemented!"
	}
	if r.Args.Bitbucket != nil {
		base = "Bitbucket"
		end = "Unimplemented!"
	}
	return base, end
}

// GetCacheLocation returns the location of a given location
func GetCacheLocation(r *robber.Robber, location string) (string, bool) {
	if _, err := os.Stat(location); !os.IsNotExist(err) {
		return location, true
	}
	baseFolder, endFolder := getCacheHelper(r, location)
	cache := filepath.Join(os.TempDir(), "yar", baseFolder, endFolder)
	_, err := os.Stat(cache)
	return cache, !os.IsNotExist(err)
}
