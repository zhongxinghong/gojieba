package gojieba

import (
	"fmt"
	"regexp"
)

var (
	reHan  = regexp.MustCompile(`[\x{4E00}-\x{9FD5}a-zA-Z0-9+#&\._%\-]+`)
	reEng  = regexp.MustCompile(`[a-zA-Z0-9]`)
	reSkip = regexp.MustCompile(`\r\n|\pZ|\s`)
)

type routeEntry struct {
	y float64
	x int
}

// Token is the element struct of the result list of Tokenizer::Tokenize method.
type Token struct {
	word  string
	start int
	end   int
}

// GetWord returns the word field.
func (t *Token) GetWord() string {
	return t.word
}

// GetStart returns the start field.
func (t *Token) GetStart() int {
	return t.start
}

// GetEnd returns the end field.
func (t *Token) GetEnd() int {
	return t.end
}

func (t *Token) String() string {
	return fmt.Sprintf("Token('%s', %d, %d)", t.word, t.start, t.end)
}

// Keyword is the element struct of the result list of TFIDF::Extract method.
type Keyword struct {
	word  string
	score float64
}

// GetWord returns the word field.
func (kw *Keyword) GetWord() string {
	return kw.word
}

// GetScore returns the score field.
func (kw *Keyword) GetScore() float64 {
	return kw.score
}

func (kw *Keyword) String() string {
	return fmt.Sprintf("Keyword('%s', %g)", kw.word, kw.score)
}
