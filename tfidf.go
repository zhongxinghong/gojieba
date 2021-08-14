package gojieba

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zhongxinghong/gojieba/internal/common"
)

var (
	defaultStopWords = []string{
		"the", "of", "is", "and", "to", "in", "that", "we", "for", "an", "are",
		"by", "be", "as", "on", "with", "can", "if", "from", "which", "you", "it",
		"this", "then", "at", "have", "all", "not", "one", "has", "or", "that",
	}
)

// TFIDF is the core struct of this file.
type TFIDF struct {
	tokenizer *Tokenizer
	idf       *IDFTable
	stopWords common.StringSet
}

// NewTFIDF create an instance of TFIDF with the given tokenizer and idf file.
func NewTFIDF(t *Tokenizer, idf string) *TFIDF {
	tfidf := &TFIDF{
		tokenizer: t,
		idf:       NewIDFTable(idf),
		stopWords: common.NewStringSet(),
	}

	for _, sw := range defaultStopWords {
		tfidf.AddStopWord(sw)
	}
	return tfidf
}

func (t *TFIDF) String() string {
	return fmt.Sprintf("TFIDF(idf='%s', tokenizer=%v)", t.idf.GetFile(), t.tokenizer)
}

// AddStopWord add a stop word to stopWords set.
func (t *TFIDF) AddStopWord(word string) bool {
	return t.stopWords.Add(strings.ToLower(word))
}

// Extract extract topK keywords from the given sentence.
func (t *TFIDF) Extract(sentence string, topK int) []Keyword {
	wtotal := float64(0)
	wfreqs := make(map[string]float64)

	for _, word := range t.tokenizer.Cut(sentence, true) {
		r := []rune(strings.TrimSpace(word))
		if len(r) < 2 || t.stopWords.Has(strings.ToLower(word)) {
			continue
		}
		wfreqs[word]++
		wtotal++
	}

	for word := range wfreqs {
		freq, ok := t.idf.GetFreq(word)
		if !ok {
			freq = t.idf.GetMedian()
		}
		wfreqs[word] *= freq / wtotal
	}

	kws := make([]Keyword, 0, len(wfreqs))
	for word, freq := range wfreqs {
		kws = append(kws, Keyword{word, freq})
	}

	sort.Slice(kws, func(i, j int) bool {
		return kws[i].score > kws[j].score
	})

	return kws[:topK]
}
