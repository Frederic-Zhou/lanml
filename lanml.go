package lanml

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"encoding/json"

	"github.com/yanyiwu/gojieba"
)

func GetData(filepath string) (dataList [][]string) {
	data, err := ioutil.ReadFile(filepath)
	errCheck("read file error", err, true)
	list := strings.Split(string(data), "\n")
	for _, v := range list {
		dataList = append(dataList, strings.Split(v, "	"))
	}

	return
}

// 获取类别,数量和边缘概率
func GetClassis(dataList [][]string) (classis Classis) {
	classis = Classis{}
	datalen := len(dataList)
	for _, v := range dataList {
		classis[v[0]] = Class{classis[v[0]].Count + 1, 0.0}
	}

	for k, v := range classis {
		classis[k] = Class{v.Count, float64(v.Count) / float64(datalen)}
	}
	return
}

func GetWords(dataList [][]string, DictPath string, classis Classis) (words Words) {
	dataListLen := len(dataList)
	gojieba.USER_DICT_PATH = DictPath
	x := gojieba.NewJieba()
	words = Words{}
	//获取词
	for _, sentence := range dataList {
		keys := x.Cut(sentence[1], true)
		for _, k := range keys {
			words[k] = Classis{}
		}
	}

	//获得词语的出现次数
	for word := range words {
		for _, sentence := range dataList {
			if strings.Contains(sentence[1], word) {
				words[word][sentence[0]] = Class{Count: words[word][sentence[0]].Count + 1, Prob: 0}
			}
		}
	}

	//去掉高频词，计算每个词的 P(w|c)
	for word, c := range words {
		count := 0
		for k, v := range c {
			count += v.Count
			words[word][k] = Class{Count: v.Count, Prob: float64(v.Count) / float64(classis[k].Count)}
		}

		if float64(count)/float64(dataListLen) > 0.2 {
			delete(words, word)
		}

	}

	return
}

func WriteWords(words Words, filePath string) (err error) {
	_, err = os.Stat(filePath)
	if os.IsExist(err) {
		os.Remove(filePath)
	}

	data, err := json.Marshal(words)
	errCheck("数据转换错误", err, true)

	err = ioutil.WriteFile(filePath, data, 0666)

	return
}

func ReadWords(filePath string) (words Words, err error) {
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		errCheck("无数据文件", err, true)
	}

	words = Words{}

	data, err := ioutil.ReadFile(filePath)
	errCheck("读取文件错误", err, true)
	err = json.Unmarshal(data, &words)

	return
}

func GetResult(text string, words Words, classis Classis) (result []Result) {

	resultMap := map[string]float64{}
	for name := range classis {
		resultMap[name] = 1
	}

	for word, c := range words {
		if strings.Contains(text, word) {
			for k := range resultMap {
				if _, ok := c[k]; ok {
					resultMap[k] = resultMap[k] * c[k].Prob
				} else {
					resultMap[k] = resultMap[k] * (0.1 / float64(classis[k].Count))
				}

			}
		}
	}

	for name, v := range classis {
		resultMap[name] = resultMap[name] * v.Prob
	}

	result = []Result{}
	for name, prob := range resultMap {
		result = append(result, Result{ClassName: name, Prob: prob})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Prob < result[j].Prob })
	return
}

type Result struct {
	ClassName string
	Prob      float64
}

type Words map[string]Classis
type Classis map[string]Class
type Class struct {
	Count int
	Prob  float64
}

func errCheck(flag string, err error, stop bool) {
	if err != nil {
		log.Println(flag, err.Error())
		if stop {
			os.Exit(0)
		}
	}

}
