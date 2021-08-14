package gojieba

import "testing"

func TestDictInit(t *testing.T) {
	NewDictionary("dict/dict.txt")
}

func TestIDFTable(t *testing.T) {
	NewIDFTable("dict/idf.txt")
}
