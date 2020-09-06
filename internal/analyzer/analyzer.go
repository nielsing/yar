package analyzer

import (
	"github.com/nielsing/yar/internal/robber"
	"math"
	"strings"
)

const (
	// B64chars is used for entropy finding of base64 strings.
	B64chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	// Hexchars is used for entropy finding of hex based strings.
	Hexchars = "1234567890abcdefABCDEF"
	// Threshold for b64 matching of entropy strings
	b64Threshold = 4.5
	// Threshold for hex matching of entropy strings
	hexThreshold = 3
)

// entropyCheck runs Shannon's Entropy on a given word
// H(X) = - \sigma{i=1}{n} P(x_i) log_bP(x_i)
// P(X = x) = P({s \in S: X(s) = x})
func entropyCheck(data string, values string) float64 {
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

func RegexSearch(r *robber.Robber, s string) string {
	return ""
}

func EntropySearch(r *robber.Robber, s string) string {
	return ""
}
