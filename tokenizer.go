package gojieba

import (
	"fmt"
	"math"
	"strings"

	"github.com/zhongxinghong/gojieba/internal/common"
	"github.com/zhongxinghong/gojieba/internal/hmm"
)

// Tokenizer is the core struct of this file.
type Tokenizer struct {
	dict *Dictionary
	hmm  *hmm.HMM
}

// NewTokenizer create an instance of Tokenizer with the given dict file.
func NewTokenizer(dict string) *Tokenizer {
	return &Tokenizer{
		dict: NewDictionary(dict),
		hmm:  hmm.NewHMM(),
	}
}

func (t *Tokenizer) String() string {
	return fmt.Sprintf("Tokenizer(dict: '%s')", t.dict.GetFile())
}

// Cut implement the default mode (accurate mode) of jieba.cut method.
func (t *Tokenizer) Cut(sentence string, HMM bool) []string {
	return t.cut(sentence, false, false, HMM)
}

// CutAll implement the full mode of jieba.cut method.
func (t *Tokenizer) CutAll(setence string, HMM bool) []string {
	return t.cut(setence, true, false, HMM)
}

// CutForSearch implement jieba.cut_for_search method.
func (t *Tokenizer) CutForSearch(sentence string, HMM bool) []string {
	return t.cut(sentence, false, true, HMM)
}

// Tokenize implement the default mode (accurate mode) jieba.tokenize method.
func (t *Tokenizer) Tokenize(sentence string, HMM bool) []Token {
	pos := 0
	res := []Token{}
	for _, word := range t.Cut(sentence, HMM) {
		width := len(word)
		res = append(res, Token{word, pos, pos + width})
		pos += width
	}
	return res
}

// TokenizeForSearch implement the search mode of jieba.tokenize method.
func (t *Tokenizer) TokenizeForSearch(sentence string, HMM bool) []Token {
	pos := 0
	res := []Token{}
	for _, word := range t.Cut(sentence, HMM) {
		rix := make([]int, 0, len(word)+1)
		rlen := 0
		for i := range word {
			rix = append(rix, i)
			rlen++
		}
		rix = append(rix, len(word))

		if rlen > 2 {
			for i, n := 0, rlen-1; i < n; i++ {
				st, ed := rix[i], rix[i+2]
				gram2 := word[st:ed]
				if freq, ok := t.dict.GetFreq(gram2); ok && freq != 0 {
					res = append(res, Token{gram2, pos + st, pos + ed})
				}
			}
		}
		if rlen > 3 {
			for i, n := 0, rlen-2; i < n; i++ {
				st, ed := rix[i], rix[i+3]
				gram3 := word[st:ed]
				if freq, ok := t.dict.GetFreq(gram3); ok && freq != 0 {
					res = append(res, Token{gram3, pos + st, pos + ed})
				}
			}
		}

		width := len(word)
		res = append(res, Token{word, pos, pos + width})
		pos += width
	}
	return res
}

// AddWord implement jieba.add_word method.
func (t *Tokenizer) AddWord(word string, freq int64, tag string) {
	t.dict.AddWord(word, freq, tag)
	if freq == 0 {
		t.hmm.AddForceSplit(word)
	} else {
		t.hmm.DelForceSplit(word)
	}
}

// DelWord implement jieba.del_word method.
func (t *Tokenizer) DelWord(word string) {
	tag, ok := t.dict.GetTag(word)
	if !ok {
		tag = ""
	}
	t.AddWord(word, 0, tag)
}

// SuggestFreqForMerge implement jieba.suggest_freq method when param segment is
// the word that should be treated as a whole.
func (t *Tokenizer) SuggestFreqForMerge(word string, tune bool) int64 {
	fTotal := float64(t.dict.total)
	fFreq := float64(1)

	for _, seg := range t.Cut(word, false) {
		if freq, ok := t.dict.GetFreq(seg); ok {
			fFreq *= float64(freq)
		}
		fFreq /= fTotal
	}

	iFreq := int64(fFreq*fTotal) + 1
	freq, ok := t.dict.GetFreq(word)
	if !ok {
		freq = 1
	}
	if freq > iFreq {
		iFreq = freq
	}

	if tune {
		tag, ok := t.dict.GetTag(word)
		if !ok {
			tag = ""
		}
		t.AddWord(word, iFreq, tag)
	}
	return iFreq
}

// SuggestFreqForSplit implement jieba.suggest_freq method when param segment is
// the segments list that the word is expected to be cut into.
func (t *Tokenizer) SuggestFreqForSplit(segments []string, tune bool) int64 {
	fTotal := float64(t.dict.total)
	fFreq := float64(1)
	word := strings.Join(segments, "")

	for _, seg := range segments {
		if freq, ok := t.dict.GetFreq(seg); ok {
			fFreq *= float64(freq)
		}
		fFreq /= fTotal
	}

	iFreq := int64(fFreq * fTotal)
	freq, ok := t.dict.GetFreq(word)
	if !ok {
		freq = 0
	}
	if freq < iFreq {
		iFreq = freq
	}

	if tune {
		tag, ok := t.dict.GetTag(word)
		if !ok {
			tag = ""
		}
		t.AddWord(word, iFreq, tag)
	}
	return iFreq
}

func (t *Tokenizer) cut(sentence string, all, search, HMM bool) []string {
	var blockCut func(string) []string

	switch {
	case all:
		blockCut = t.blockCutAll
	case HMM:
		blockCut = t.blockCutDAG
	default:
		blockCut = t.blockCutDAGNoHMM
	}

	res := []string{}

	for i, block := range common.PyReSplit(reHan, sentence) {
		if len(block) == 0 {
			continue
		}
		if i%2 == 1 { // reHan.MatchString(block)
			for _, word := range blockCut(block) {
				res = append(res, word)
			}
		} else {
			for j, word := range common.PyReSplit(reSkip, block) {
				if j%2 == 1 { // reSkip.MatchString(word)
					res = append(res, word)
				} else if !all {
					for _, r := range word {
						res = append(res, string(r))
					}
				} else {
					res = append(res, word)
				}
			}
		}
	}

	if search {
		words := res
		res = make([]string, 0, len(words))

		for _, word := range words {
			r := []rune(word)
			if len(r) > 2 {
				for i, n := 0, len(r)-1; i < n; i++ {
					gram2 := string(r[i : i+2])
					if freq, ok := t.dict.GetFreq(gram2); ok && freq != 0 {
						res = append(res, gram2)
					}
				}
			}
			if len(r) > 3 {
				for i, n := 0, len(r)-2; i < n; i++ {
					gram3 := string(r[i : i+3])
					if freq, ok := t.dict.GetFreq(gram3); ok && freq != 0 {
						res = append(res, gram3)
					}
				}
			}
			res = append(res, word)
		}
	}

	return res
}

func (t *Tokenizer) getDAG(r []rune) [][]int {
	n := len(r)
	dag := make([][]int, n)

	for k := 0; k < n; k++ {
		lst := []int{}
		frag := string(r[k])
		for i := k; i < n; {
			freq, ok := t.dict.GetFreq(frag)
			if !ok {
				break
			}
			if freq != 0 {
				lst = append(lst, i)
			}
			i++
			if i < n {
				frag = string(r[k : i+1])
			}
		}
		if len(lst) == 0 {
			lst = append(lst, k)
		}
		dag[k] = lst
	}
	return dag
}

func (t *Tokenizer) getRoute(r []rune, dag [][]int) []routeEntry {
	n := len(r)
	route := make([]routeEntry, n+1)
	route[n] = routeEntry{0, 0}
	logTotal := math.Log(float64(t.dict.GetTotal()))

	for i := n - 1; i >= 0; i-- {
		maxy := math.Inf(-1)
		maxx := -1
		for _, x := range dag[i] {
			freq, ok := t.dict.GetFreq(string(r[i : x+1]))
			if !ok || freq == 0 {
				freq = 1
			}
			y := math.Log(float64(freq)) - logTotal + route[x+1].y
			if y > maxy || (y == maxy && x > maxx) {
				maxy = y
				maxx = x
			}
		}
		route[i] = routeEntry{maxy, maxx}
	}
	return route
}

func (t *Tokenizer) blockCutAll(block string) []string {
	r := []rune(block)
	dag := t.getDAG(r)
	j0 := -1
	engScan := false
	engBuf := strings.Builder{}
	res := []string{}

	for k, lst := range dag {
		if engScan && !reEng.MatchString(string(r[k])) {
			engScan = false
			res = append(res, engBuf.String())
		}
		if len(lst) == 1 && k > j0 {
			word := string(r[k : lst[0]+1])
			if reEng.MatchString(word) {
				if !engScan {
					engScan = true
					engBuf.Reset()
				}
				engBuf.WriteString(word)
			}
			if !engScan {
				res = append(res, word)
			}
			j0 = lst[0]
		} else {
			for _, j := range lst {
				if j > k {
					res = append(res, string(r[k:j+1]))
					j0 = j
				}
			}
		}
	}
	if engScan {
		res = append(res, engBuf.String())
	}
	return res
}

func (t *Tokenizer) blockCutDAG(block string) []string {
	r := []rune(block)
	n := len(r)
	dag := t.getDAG(r)
	route := t.getRoute(r, dag)
	buf := strings.Builder{}
	res := []string{}

	processBuf := func(buf *strings.Builder) {
		if buf.Len() > 0 {
			sbuf := buf.String()
			rbuf := []rune(sbuf)
			if len(rbuf) == 1 {
				res = append(res, sbuf)
			} else if freq, ok := t.dict.GetFreq(sbuf); !ok || freq == 0 {
				for _, word := range t.hmm.Cut(sbuf) {
					res = append(res, word)
				}
			} else {
				for _, r := range rbuf {
					res = append(res, string(r))
				}
			}
			buf.Reset()
		}
	}

	for x := 0; x < n; {
		y := route[x].x + 1
		word := string(r[x:y])
		if y-x == 1 {
			buf.WriteString(word)
		} else {
			processBuf(&buf)
			res = append(res, word)
		}
		x = y
	}
	processBuf(&buf)

	return res
}

func (t *Tokenizer) blockCutDAGNoHMM(block string) []string {
	r := []rune(block)
	n := len(r)
	dag := t.getDAG(r)
	route := t.getRoute(r, dag)
	buf := strings.Builder{}
	res := []string{}

	processBuf := func(buf *strings.Builder) {
		if buf.Len() > 0 {
			res = append(res, buf.String())
			buf.Reset()
		}
	}

	for x := 0; x < n; {
		y := route[x].x + 1
		word := string(r[x:y])
		if y-x == 1 && reEng.MatchString(word) {
			buf.WriteString(word)
		} else {
			processBuf(&buf)
			res = append(res, word)
		}
		x = y
	}
	processBuf(&buf)

	return res
}
