package hmm

import (
	"math"

	"github.com/zhongxinghong/gojieba/internal/common"
)

// HMM is the core struct of this package. It implement jieba.finalseg module.
type HMM struct {
	forceSplitWords common.StringSet
}

// NewHMM create an instance of HMM.
func NewHMM() *HMM {
	return &HMM{
		forceSplitWords: common.NewStringSet(),
	}
}

// AddForceSplit implement jieba.finalseg.add_force_split method.
func (hmm *HMM) AddForceSplit(word string) bool {
	return hmm.forceSplitWords.Add(word)
}

// DelForceSplit implement the inverse operation of jieba.finalseg.add_force_split
// method.
func (hmm *HMM) DelForceSplit(word string) bool {
	return hmm.forceSplitWords.Del(word)
}

// Cut implement jieba.finalseg.cut method.
func (hmm *HMM) Cut(sentence string) []string {
	res := []string{}

	for i, block := range common.PyReSplit(reHan, sentence) {
		if i%2 == 1 { // reHan.MatchString(block)
			for _, word := range hmm.blockCut(block) {
				if !hmm.forceSplitWords.Has(word) {
					res = append(res, word)
				} else {
					for _, r := range word {
						res = append(res, string(r))
					}
				}
			}
		} else {
			for _, word := range common.PyReSplit(reSkip, block) {
				if len(word) > 0 {
					res = append(res, word)
				}
			}
		}
	}
	return res
}

func (hmm *HMM) blockCut(block string) []string {
	states := []int{stateB, stateM, stateE, stateS}
	r := []rune(block)
	n := len(r)
	V := make([][4]float64, n)
	path := [4][]int{}

	for _, y := range states {
		ep, ok := emitP[y][r[0]]
		if !ok {
			ep = minFloat
		}
		V[0][y] = startP[y] + ep
		path[y] = []int{y}
	}

	for t := 1; t < n; t++ {
		newPath := [4][]int{}
		for _, y := range states {
			ep, ok := emitP[y][r[t]]
			if !ok {
				ep = minFloat
			}
			maxProb := math.Inf(-1)
			maxState := -1
			for _, y0 := range prevState[y] {
				prob := V[t-1][y0] + transP[y0][y] + ep
				if prob > maxProb || (prob == maxProb && y0 > maxState) {
					maxProb = prob
					maxState = y0
				}
			}
			V[t][y] = maxProb
			// WARNING: MUST BE copy, NOT newPath[y] = append(path[maxState], y)
			newPath[y] = make([]int, len(path[maxState])+1)
			copy(newPath[y], append(path[maxState], y))
		}
		path = newPath
	}

	states = []int{stateE, stateS}
	maxProb := math.Inf(-1)
	maxState := -1

	for _, y := range states {
		prob := V[n-1][y]
		if prob > maxProb || (prob == maxProb && y > maxState) {
			maxProb = prob
			maxState = y
		}
	}

	posList := path[maxState]

	begin, ni := 0, 0
	res := []string{}

	for i := range r {
		switch posList[i] {
		case stateB:
			begin = i
		case stateE:
			res = append(res, string(r[begin:i+1]))
			ni = i + 1
		case stateS:
			res = append(res, string(r[i]))
			ni = i + 1
		}
	}
	if ni < n {
		res = append(res, string(r[ni:]))
	}
	return res
}
