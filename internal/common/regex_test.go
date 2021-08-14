package common

import (
	"reflect"
	"regexp"
	"testing"
)

var (
	reHan = regexp.MustCompile(`[\x{4E00}-\x{9FD5}a-zA-Z0-9+#&\._%\-]+`)

	s1 = "你好，我来到Beijing天安门，我来到 Peking Univ，我看到THU，   ohhh ~"
)

func TestRegexpSplit(t *testing.T) {
	t.Logf("%v\n", reHan.Split(s1, -1))
}

func TestReHanSplit(t *testing.T) {
	got := PyReSplit(reHan, s1)
	want := []string{"", "你好", "，", "我来到Beijing天安门", "，", "我来到", " ",
		"Peking", " ", "Univ", "，", "我看到THU", "，   ", "ohhh", " ~"}
	for _, e := range got {
		t.Logf("'%s'\n", e)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
