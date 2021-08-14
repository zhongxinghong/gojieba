package gojieba

import (
	"reflect"
	"strings"
	"testing"
)

func TestCut(t *testing.T) {
	var got, want []string

	dt := NewTokenizer("./dict/dict.txt")
	t.Logf("%v\n", dt)

	got = dt.CutAll("我来到北京清华大学", true)
	want = []string{"我", "来到", "北京", "清华", "清华大学", "华大", "大学"}
	logResult(t, "Full mode", got)
	checkResult(t, got, want)

	got = dt.Cut("我来到北京清华大学", true)
	want = []string{"我", "来到", "北京", "清华大学"}
	logResult(t, "Accurate mode", got)
	checkResult(t, got, want)

	got = dt.Cut("他来到了网易杭研大厦", true)
	want = []string{"他", "来到", "了", "网易", "杭研", "大厦"}
	logResult(t, "New word '杭研'", got)
	checkResult(t, got, want)

	got = dt.Cut("小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true)
	want = []string{"小明", "硕士", "毕业", "于", "中国科学院", "计算所", "，",
		"后", "在", "日本京都大学", "深造"}
	logResult(t, "Accurate mode", got)
	checkResult(t, got, want)

	got = dt.CutForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true)
	want = []string{"小明", "硕士", "毕业", "于", "中国", "科学", "学院", "科学院",
		"中国科学院", "计算", "计算所", "，", "后", "在", "日本", "京都", "大学",
		"日本京都大学", "深造"}
	logResult(t, "Search mode", got)
	checkResult(t, got, want)

	got = dt.Cut("但是有考据癖的人也当然不肯错过索隐的杨会、放弃附会的权利的。", true)
	want = []string{"但是", "有", "考据", "癖", "的", "人", "也", "当然", "不肯",
		"错过", "索隐", "的", "杨会", "、", "放弃", "附会", "的", "权利", "的", "。"}
	logResult(t, "Accurate mode", got)
	checkResult(t, got, want)

	got = dt.CutAll("署名Ｐａｔｒｉｃ　Ｍａｈｏｎｅｙ", true)
	want = []string{"署名", "Ｐａｔｒｉｃ", "\u3000", "Ｍａｈｏｎｅｙ"}
	logResult(t, "Full mode", got)
	checkResult(t, got, want)
}

func TestTokenize(t *testing.T) {
	s := "小明硕士毕业于中国科学院计算所，后在日本京都大学深造"

	var got, want []string
	var tokens []Token

	dt := NewTokenizer("./dict/dict.txt")
	t.Logf("%v\n", dt)

	tokens = dt.Tokenize(s, true)
	t.Logf("%v\n", tokens)

	got = make([]string, 0, len(tokens))

	for _, token := range tokens {
		if token.GetWord() != s[token.GetStart():token.GetEnd()] {
			t.Errorf("Incorrect Token: %v\n", token)
		}
		got = append(got, token.GetWord())
	}

	want = []string{"小明", "硕士", "毕业", "于", "中国科学院", "计算所", "，",
		"后", "在", "日本京都大学", "深造"}
	checkResult(t, got, want)

	tokens = dt.TokenizeForSearch(s, true)
	t.Logf("%v\n", tokens)

	got = make([]string, 0, len(tokens))

	for _, token := range tokens {
		if token.GetWord() != s[token.GetStart():token.GetEnd()] {
			t.Errorf("Incorrect Token: %v\n", token)
		}
		got = append(got, token.GetWord())
	}

	want = []string{"小明", "硕士", "毕业", "于", "中国", "科学", "学院", "科学院",
		"中国科学院", "计算", "计算所", "，", "后", "在", "日本", "京都", "大学",
		"日本京都大学", "深造"}
	checkResult(t, got, want)
}

func TestAddDelWord(t *testing.T) {
	var got, want []string

	dt := NewTokenizer("./dict/dict.txt")
	t.Logf("%v\n", dt)

	got = dt.Cut("李小福是创新办主任也是云计算方面的专家", true)
	want = []string{"李小福", "是", "创新", "办", "主任", "也", "是", "云",
		"计算", "方面", "的", "专家"}
	logResult(t, "Before add word", got)
	checkResult(t, got, want)

	dt.AddWord("创新办", 3, "i")
	dt.AddWord("云计算", 5, "")

	got = dt.Cut("李小福是创新办主任也是云计算方面的专家", true)
	want = []string{"李小福", "是", "创新办", "主任", "也", "是", "云计算",
		"方面", "的", "专家"}
	logResult(t, "After add word", got)
	checkResult(t, got, want)

	dt.DelWord("创新办")
	dt.DelWord("云计算")

	got = dt.Cut("李小福是创新办主任也是云计算方面的专家", true)
	want = []string{"李小福", "是", "创新", "办", "主任", "也", "是", "云",
		"计算", "方面", "的", "专家"}
	logResult(t, "After del word", got)
	checkResult(t, got, want)
}

func TestSuggestFreq(t *testing.T) {
	var got, want []string
	var freq int64

	dt := NewTokenizer("./dict/dict.txt")
	t.Logf("%v\n", dt)

	got = dt.Cut("如果放到post中将出错。", false)
	want = []string{"如果", "放到", "post", "中将", "出错", "。"}
	logResult(t, "Before suggest", got)
	checkResult(t, got, want)

	freq = dt.SuggestFreqForSplit([]string{"中", "将"}, true)
	t.Logf("dt.SuggestFreqForSplit: %v\n", freq)

	got = dt.Cut("如果放到post中将出错。", false)
	want = []string{"如果", "放到", "post", "中", "将", "出错", "。"}
	logResult(t, "After suggest", got)
	checkResult(t, got, want)

	got = dt.Cut("「台中」正确应该不会被切开", false)
	want = []string{"「", "台", "中", "」", "正确", "应该", "不会", "被", "切开"}
	logResult(t, "Before suggest", got)
	checkResult(t, got, want)

	freq = dt.SuggestFreqForMerge("台中", true)
	t.Logf("dt.SuggestFreqForMerge: %v\n", freq)

	got = dt.Cut("「台中」正确应该不会被切开", false)
	want = []string{"「", "台中", "」", "正确", "应该", "不会", "被", "切开"}
	logResult(t, "After suggest", got)
	checkResult(t, got, want)
}

func logResult(t *testing.T, prefix string, res []string) {
	buf := strings.Builder{}
	if len(prefix) > 0 {
		buf.WriteString(prefix)
		buf.WriteString(": ")
	}
	buf.WriteString("[")
	n := len(res)
	for i, word := range res {
		buf.WriteString(word)
		if i != n-1 {
			buf.WriteString(" / ")
		}
	}
	buf.WriteString("]\n")
	t.Logf(buf.String())
}

func checkResult(t *testing.T, got, want []string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got != want, got: %v, want: %v", got, want)
	}
}
