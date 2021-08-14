package hmm

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"io"
)

//go:embed model
var modelFs embed.FS

var (
	startP    [4]float64
	transP    [4][4]float64
	emitP     [4]map[rune]float64
	prevState [4][]int
)

const (
	minFloat = -3.14e100

	stateB = 0
	stateM = 1
	stateE = 2
	stateS = 3
)

func byte2state(state byte) int {
	switch state {
	case 'B':
		return stateB
	case 'M':
		return stateM
	case 'E':
		return stateE
	case 'S':
		return stateS
	default:
		panic(state)
	}
}

func loadStartP() {
	fp, err := modelFs.Open("model/prob_start.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	rd := bufio.NewReader(fp)

	var (
		state byte
		prob  float64
	)

	for {
		_, err = fmt.Fscanf(rd, "%c %g\n", &state, &prob)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		state := byte2state(state)
		startP[state] = prob
	}
}

func loadTransP() {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			transP[i][j] = minFloat
		}
	}

	fp, err := modelFs.Open("model/prob_trans.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	rd := bufio.NewReader(fp)

	var (
		prev  byte
		state byte
		prob  float64
	)

	for {
		_, err = fmt.Fscanf(rd, "%c %c %g\n", &prev, &state, &prob)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		prev := byte2state(prev)
		state := byte2state(state)
		transP[prev][state] = prob
	}
}

func loadEmitP() {
	for i := 0; i < 4; i++ {
		emitP[i] = make(map[rune]float64)
	}

	fp, err := modelFs.Open("model/prob_emit.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	rd := bufio.NewReader(fp)

	var (
		state byte
		word  rune
		prob  float64
	)

	for {
		_, err := fmt.Fscanf(rd, "%c %c %g\n", &state, &word, &prob)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		state := byte2state(state)
		emitP[state][word] = prob
	}
}

func buildPrevState() {
	prevState[stateB] = []int{stateE, stateS}
	prevState[stateM] = []int{stateM, stateB}
	prevState[stateS] = []int{stateS, stateE}
	prevState[stateE] = []int{stateB, stateM}
}

func init() {
	loadStartP()
	loadTransP()
	loadEmitP()
	buildPrevState()
}
