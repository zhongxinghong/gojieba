## GoJieba

[Jieba 0.42.1](https://github.com/fxsjy/jieba/tree/67fa2e36e72f69d9134b8a1037b83fbb070b9775) 的部分复刻版，支持原版 Jieba 的如下功能：

- 基于 DAG 和 HMM 的中文分词 (无词性标注)
- 基于 TFIDF 的关键词提取

## 写在前面

这是一个我为了熟悉 Golang 而拿来练手的项目，类似于 [Java-Jieba](https://github.com/zhongxinghong/Java-Jieba)，可以保证已实现的功能是正确的，但是未来不会再继续维护

## 兼容性

使用到了 `//go:embed`，因此至少需要 Go 1.16

## 基本用法

基本用法参考 `tokenizer_test.go` 和 `tfidf_test.go`，准确重现了 Jieba [README](https://github.com/fxsjy/jieba/blob/master/README.md) 和 [test](https://github.com/fxsjy/jieba/tree/master/test) 中的案例

## 基础测试

基本的功能测试参考 `tokenizer_test.go` 和 `tfidf_test.go`，大规模的测试参考 `test/gojieba_test.go`

## 一致性说明

GoJieba 复刻的模块的运行结果与原版 Jieba 基本上一致。由于 Python 的正则表达式和 Golang 的正则表达式对 Unicode 的支持存在差别，用 Golang 的 regexp 分割空白字符，其结果与 Python 的 re 在 re.U 选项下的分割结果存在出入，但是在 Trim 后进行比较，两者的结果是**完全一致**的，比较的方法参看 `test/gojieba_test.go` 的 `isEqual`，通过 `script/diff_test_result.py` 比较的结果也是**完全一致**的。因此，基本上可以认为 GoJieba 与原版 Jieba 的分词结果仅仅在**空白字符**的划分上存在差别

GoJieba 的 TFIDF 关键词提取功能在小规模数据上的测试结果与原版 Jieba **完全一致**

## 性能测试

### 测试环境

- go version go1.16.7 windows/amd64
- Windows 10 Pro 1903 Build 18362.1256
- Intel(R) Core(TM) i5-3320M CPU @ 2.60GHz

### GoJieba

以 `test/gojieba_test.go` 中的 `BenchmarkCut` 函数为例，测试在精确模式下对《围城》全文进行分词的性能，结果如下，平均耗时 0.40 秒

```powershell
PS E:\go\src\github.com\zhongxinghong\gojieba\test> go test -cpu 1 -benchmem -benchtime 20x -run ^$ -bench BenchmarkCut -v
goos: windows
goarch: amd64
pkg: github.com/zhongxinghong/gojieba/test
cpu: Intel(R) Core(TM) i5-3320M CPU @ 2.60GHz
BenchmarkCut
BenchmarkCut          20         403216190 ns/op        54515138 B/op    1417561 allocs/op
PASS
ok      github.com/zhongxinghong/gojieba/test   11.663s
```

### Jieba

类似地，在 IPython shell 中运行下述代码，做相同的测试，结果如下，平均耗时 2.25 秒

```ipython
In [22]: with open('../test/artical03.in.txt', 'r', encoding='utf-8') as fp:
    ...:     content = fp.read()
    ...:

In [23]: %timeit jieba.lcut(content)
2.25 s ± 67.4 ms per loop (mean ± std. dev. of 7 runs, 1 loop each)
```

### 小结

基于上述两个性能测试，可以给出一个大致的结论，GoJieba 在精确模式下的分词速度约为 Jieba 的 5 - 6 倍

## 证书

- [Jieba](https://github.com/fxsjy/jieba/blob/master/LICENSE)
- [GoJieba](https://github.com/zhongxinghong/gojieba/blob/master/LICENSE)
