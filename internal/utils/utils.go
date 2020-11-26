package utils

import (
	"io/ioutil"
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
		base = "git"
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
		base = "github"
		user := filepath.Base(filepath.Dir(location))
		repo := strings.Replace(filepath.Base(location), ".git", "", -1)
		filepath.Join(end, user, repo)
	}
	if r.Args.Gitlab != nil {
		base = "gitlab"
		end = "Unimplemented!"
	}
	if r.Args.Bitbucket != nil {
		base = "bitbucket"
		end = "Unimplemented!"
	}
	return base, end
}

// GetCacheLocation returns the cache location and whether it exists or not.
func GetCacheLocation(r *robber.Robber, location string) (string, bool) {
	if _, err := os.Stat(location); !os.IsNotExist(err) {
		return location, true
	}
	baseFolder, endFolder := getCacheHelper(r, location)
	cache := filepath.Join(os.TempDir(), "yar", baseFolder, endFolder)
	_, err := os.Stat(cache)
	return cache, !os.IsNotExist(err)
}

// ChunkString chunks a given string `s` into an array of string chunks of size `chunkSize`.
func ChunkString(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}

	var chunks []string
	chunk := make([]rune, chunkSize)
	currSize := 0

	for _, r := range s {
		chunk[currSize] = r
		currSize++
		if currSize == chunkSize {
			chunks = append(chunks, string(chunk))
			currSize = 0
		}
	}
	if currSize > 0 {
		chunks = append(chunks, string(chunk[:currSize]))
	}
	return chunks
}

// SplitAndChunkString splits a given string `s` into substrings separated by `sep` and then chunks
// the splits into chunks of size `chunkSize`.
func SplitAndChunkString(s, sep string, chunkSize int) []string {
	var chunks []string
	splits := strings.Split(s, sep)

	for _, split := range splits {
		chunks = append(chunks, ChunkString(split, chunkSize)...)
	}
	return chunks
}

func WriteToFile(filename string, values []*string) error {
	unRefValues := []string{}
	for _, refValue := range values {
		unRefValues = append(unRefValues, *refValue)
	}

	value := []byte(strings.Join(unRefValues, "\n"))
	err := ioutil.WriteFile(filename, value, 0644)
	return err
}
