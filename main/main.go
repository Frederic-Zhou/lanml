package main

import (
	"bufio"
	"io/ioutil"
	"strings"

	"log"

	"os"

	"fmt"

	"sort"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yanyiwu/gojieba"
)

const (
	DictPath = "../custom.dict"
	WFJD     = "违法解读"
	LKXX     = "路况信息"
	QT       = "其他"
)

var classis = map[string]([]float64){
	WFJD: []float64{0, 0.0},
	LKXX: []float64{0, 0.0},
	QT:   []float64{0, 0.0}}

func main() {

	gojieba.USER_DICT_PATH = DictPath

	x := gojieba.NewJieba()

	//将所有文本数据保存到数组中
	data, err := ioutil.ReadFile("../answers.txt")
	errCheck("read file error", err, true)
	dataList := strings.Split(string(data), "\n")
	dataListLen := float64(len(dataList))
	fmt.Println(dataListLen)

	//分词，得到所有关键词
	words := map[string]([]float64){}
	for _, sentence := range dataList {
		a := strings.Split(sentence, "	")
		if len(a) > 0 {
			classis[a[0]][0]++
			keys := x.Cut(a[1], true)
			for _, k := range keys {
				words[k] = []float64{0.5, 0.5, 0.5, 0.5}
			}
		}
	}

	for k, v := range classis {
		classis[k][1] = v[0] / dataListLen
	}

	fmt.Println(classis)

	//获得词语的出现次数,去掉出现频率太高的词
	for word := range words {
		for _, sentence := range dataList {
			a := strings.Split(sentence, "	")
			if strings.Contains(a[1], word) {
				words[word][0]++
				switch a[0] {
				case WFJD:
					words[word][1]++
				case LKXX:
					words[word][2]++
				default:
					words[word][3]++

				}
			}
		}

		//去掉出现频率太高的词,阀值0.1(10%)
		if words[word][0]/dataListLen > 0.1 {
			delete(words, word)
		}
	}

	fmt.Println(words)

	//测试
	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "stop" {
			running = false
		}
		log.Println("command", command)
		Test(words, classis, command)
	}

}

func Test(words map[string][]float64, classis map[string][]float64, text string) {

	var p1, p2, p3 float64 = 1, 1, 1
	for word, v := range words {
		if strings.Contains(text, word) {
			p1 = p1 * v[1] / classis[WFJD][0]
			p2 = p2 * v[2] / classis[LKXX][0]
			p3 = p3 * v[3] / classis[QT][0]
		}
	}

	p1 = p1 * classis[WFJD][1]
	p2 = p2 * classis[LKXX][1]
	p3 = p3 * classis[QT][1]

	fmt.Println("违法解读", p1, ";路况信息", p2, ";其他", p3)
	f := []float64{p1, p2, p3}
	sort.Float64s(f)

	fmt.Println(f)
}

func errCheck(flag string, err error, stop bool) {
	if err != nil {
		log.Println(flag, err.Error())
		if stop {
			os.Exit(0)
		}
	}

}
