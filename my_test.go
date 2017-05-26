package lanml

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

import "log"

func TestMyFunc(t *testing.T) {

	dl := GetData("trained.txt")

	class := GetClassis(dl)

	ws1 := GetWords(dl, "./custom.dict", class)

	err := WriteWords(ws1, "words.dat")
	if err != nil {
		log.Println("write words", err.Error())
	}

	ws, err := ReadWords("words.dat")
	if err != nil {
		log.Println("read words", err.Error())
	}

	data, err := ioutil.ReadFile("./test.txt")
	errCheck("读取测试数据", err, true)
	datastring := string(data)
	datalist := strings.Split(datastring, "\n")
	sCount := 0
	for k, v := range datalist {
		l := strings.Split(v, "	")

		r := GetResult(l[1], ws, class)

		if r[len(r)-1].ClassName == l[0] {
			fmt.Println(k, r[len(r)-1].ClassName, l[0], l[1], "成功")
			sCount++
		} else {
			fmt.Println(k, r[len(r)-1].ClassName, l[0], l[1], "失败")
		}
	}

	fmt.Println("成功率:", float64(sCount)/float64(len(datalist))*100)

}
