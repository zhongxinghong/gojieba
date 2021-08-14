package common

import (
	"testing"
)

func TestStringSet(t *testing.T) {
	s := NewStringSet()
	t.Logf("%v\n", s)

	s.Add("你好")
	s.Add("世界")

	if s.Add("你好") {
		t.Errorf("s.Add")
	}

	t.Logf("%v\n", s)

	if !s.Has("你好") {
		t.Errorf("s.Has")
	}

	if !s.Del("你好") {
		t.Errorf("t.Del")
	}
	if s.Del("你好") {
		t.Errorf("t.Del")
	}

	if s.Len() != 1 {
		t.Errorf("s.Len")
	}

	t.Logf("%v\n", s)
}
