package test

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/zhongxinghong/gojieba"
)

var (
	reInput = regexp.MustCompile(`^(.*?)\.in\.txt$`)
)

func TestCut(t *testing.T) {
	dt := gojieba.NewTokenizer("../dict/dict.txt")
	t.Logf("%v\n", dt)

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rd, err := os.ReadDir(pwd)
	if err != nil {
		panic(err)
	}

	for _, de := range rd {
		if de.IsDir() {
			continue
		}
		match := reInput.FindAllStringSubmatch(de.Name(), 1)
		if len(match) == 0 {
			continue
		}

		basename := match[0][1]

		inFile := filepath.Join(pwd, de.Name())
		content := loadInput(inFile)

		var got, want []string
		var outFile, stdoutFile string

		outFile = filepath.Join(pwd, basename+".accurate.out.txt")
		stdoutFile = filepath.Join(pwd, basename+".accurate.stdout.txt")

		want = loadWantResult(stdoutFile)

		got = dt.Cut(content, true)
		dumpGotResult(got, outFile)

		if !isEqual(got, want) {
			t.Errorf("case: %s, got != want, mode: accurate", basename)
		}

		outFile = filepath.Join(pwd, basename+".all.out.txt")
		stdoutFile = filepath.Join(pwd, basename+".all.stdout.txt")

		want = loadWantResult(stdoutFile)

		got = dt.CutAll(content, true)
		dumpGotResult(got, outFile)

		if !isEqual(got, want) {
			t.Errorf("case: %s, got != want, mode: all", basename)
		}

		outFile = filepath.Join(pwd, basename+".search.out.txt")
		stdoutFile = filepath.Join(pwd, basename+".search.stdout.txt")

		want = loadWantResult(stdoutFile)

		got = dt.CutForSearch(content, true)
		dumpGotResult(got, outFile)

		if !isEqual(got, want) {
			t.Errorf("case: %s, got != want, mode: search", basename)
		}
	}
}

func BenchmarkCut(b *testing.B) {
	dt := gojieba.NewTokenizer("../dict/dict.txt")

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	inFile := filepath.Join(pwd, "artical03.in.txt")
	content := loadInput(inFile)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dt.Cut(content, true)
	}
}

func isEqual(got, want []string) bool {
	i, n := 0, len(got)
	j, m := 0, len(want)

	var x, y string

	for i < n || j < m {
		if i < n {
			x = strings.TrimSpace(got[i])
			if len(x) == 0 {
				i++
				continue
			}
		}
		if j < m {
			y = strings.TrimSpace(want[j])
			if len(y) == 0 {
				j++
				continue
			}
		}
		if x != y {
			return false
		}
		i++
		j++
	}
	return i == n && j == m
}

func loadInput(file string) string {
	fp, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func loadWantResult(file string) []string {
	fp, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	res := []string{}

	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return res
}

func dumpGotResult(got []string, file string) {
	fp, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	wt := bufio.NewWriter(fp)

	for _, word := range got {
		wt.WriteString(word)
		wt.WriteByte('\n')
	}

	wt.Flush()
}
