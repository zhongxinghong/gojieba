package common

import "regexp"

// PyReSplit implement re.split in Python when the whole pattern is grouped.
// It will be used by reHan/reSkip/... regexp.
func PyReSplit(re *regexp.Regexp, s string) []string {
	indices := re.FindAllStringIndex(s, -1)
	res := make([]string, 0, len(indices)*2+1)
	pos := 0
	for _, span := range indices {
		res = append(res, s[pos:span[0]], s[span[0]:span[1]])
		pos = span[1]
	}
	res = append(res, s[pos:])
	return res
}
