package hmm

import "regexp"

var (
	reHan  = regexp.MustCompile(`[\x{4E00}-\x{9FD5}]+`)
	reSkip = regexp.MustCompile(`[a-zA-Z0-9]+(?:\.\d+)?%?`)
)
