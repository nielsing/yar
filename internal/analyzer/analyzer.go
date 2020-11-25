package analyzer

import (
	"github.com/nielsing/yar/internal/robber"
	"math"
)

func entropyCheck(data string) float64 {
	if data == "" {
		return 0.0
	}

	var entropy float64
	charCounts := make(map[rune]int, len(data))
	for _, char := range data {
		charCounts[char]++
	}

	invLength := 1.0 / float64(len(data))
	for _, count := range charCounts {
		freq := float64(count) * invLength
		entropy -= freq * math.Log2(freq)
	}
	return entropy
}

// TODO: Comment
func RegexSearch(r *robber.Robber, line string) string {
	for _, rule := range r.Config.Rules {
		found := rule.Regex.FindString(line)
		if found == "" {
			continue
		}
		return found
	}
	return ""
}

// TODO: Comment
func EntropySearch(r *robber.Robber, line string) string {
	return ""
}
