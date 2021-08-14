package gojieba

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// Dictionary is the core struct of this file. It will be used by Tokenizer.
type Dictionary struct {
	freqs map[string]int64
	tags  map[string]string
	total int64
	file  string
}

func (d *Dictionary) String() string {
	return fmt.Sprintf("Dictionary(file='%s')", d.file)
}

// NewDictionary create an instance of Dictionary with the given dict file.
func NewDictionary(file string) *Dictionary {
	file, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}

	d := &Dictionary{
		freqs: make(map[string]int64),
		tags:  make(map[string]string),
		total: 0,
		file:  file,
	}

	fp, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	rd := bufio.NewReader(fp)

	var (
		word string
		freq int64
		tag  string
	)

	for {
		_, err := fmt.Fscanf(rd, "%s %d %s\n", &word, &freq, &tag)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		d.AddWord(word, freq, tag)
	}
	return d
}

// GetFile returns the dict field.
func (d *Dictionary) GetFile() string {
	return d.file
}

// GetTotal returns the total field.
func (d *Dictionary) GetTotal() int64 {
	return d.total
}

// GetFreq returns the freq of the given word by d.freqs[word].
func (d *Dictionary) GetFreq(word string) (int64, bool) {
	freq, ok := d.freqs[word]
	return freq, ok
}

// GetTag returns the tag of the given word by d.tags[word].
func (d *Dictionary) GetTag(word string) (string, bool) {
	tag, ok := d.tags[word]
	return tag, ok
}

// AddWord add a word with its freq and tag to dict.
func (d *Dictionary) AddWord(word string, freq int64, tag string) {
	d.freqs[word] = freq
	d.tags[word] = tag
	d.total += freq

	r := []rune(word)

	for i, n := 0, len(r); i < n; i++ {
		prefix := string(r[:i+1])
		if _, ok := d.freqs[prefix]; !ok {
			d.freqs[prefix] = 0
		}
	}
}

// IDFTable is the core struct of this file. It will be used by TFIDF.
type IDFTable struct {
	freqs  map[string]float64
	median float64
	file   string
}

func (idf *IDFTable) String() string {
	return fmt.Sprintf("IDFTable(file='%s')", idf.file)
}

// NewIDFTable create an instance of IDFTable with the given idf file.
func NewIDFTable(file string) *IDFTable {
	file, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}

	idf := &IDFTable{
		freqs:  make(map[string]float64),
		median: 0,
		file:   file,
	}

	fp, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	rd := bufio.NewReader(fp)

	var (
		word string
		freq float64
	)

	for {
		_, err := fmt.Fscanf(rd, "%s %g\n", &word, &freq)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}

		idf.freqs[word] = freq
	}

	freqs := make([]float64, 0, len(idf.freqs))
	for _, freq := range idf.freqs {
		freqs = append(freqs, freq)
	}

	sort.Float64s(freqs)
	idf.median = freqs[len(freqs)/2]

	return idf
}

// GetFile returns the file field.
func (idf *IDFTable) GetFile() string {
	return idf.file
}

// GetMedian returns the median field.
func (idf *IDFTable) GetMedian() float64 {
	return idf.median
}

// GetFreq returns the freq of the given word by idf.freqs[word].
func (idf *IDFTable) GetFreq(word string) (float64, bool) {
	freq, ok := idf.freqs[word]
	return freq, ok
}
