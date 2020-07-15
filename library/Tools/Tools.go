package Tools

import (
	"fmt"
	"github.com/Chain-Zhang/pinyin"
	"regexp"
	"strconv"
	"strings"
)

func Compare(str string, str2 string) bool {
	str, err := pinyin.New(str).Split("").Mode(pinyin.WithoutTone).Convert()
	if err != nil {
		return false
	} else {
		str2, err := pinyin.New(str2).Split("").Mode(pinyin.WithoutTone).Convert()
		if err != nil {
			return false
		} else {

			matched, _ := regexp.MatchString(strings.ToLower(str2), strings.ToLower(str))
			if matched {
				fmt.Println(str2, str)
				return true
			}
		}
	}
	return false
}
func CheckOrder(cm string, cms []string) bool {
	for _, cm2 := range cms {
		if Compare(cm, cm2) {
			return true
		}
	}
	return false
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start) // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

// 格式化数值    1,234,567,898.55
func NumberFormat(number int) string {
	str := strconv.Itoa(number)
	length := len(str)
	if length < 4 {
		return str
	}
	arr := strings.Split(str, ".") //用小数点符号分割字符串,为数组接收
	length1 := len(arr[0])
	if length1 < 4 {
		return str
	}
	count := (length1 - 1) / 3
	for i := 0; i < count; i++ {
		arr[0] = arr[0][:length1-(i+1)*3] + "," + arr[0][length1-(i+1)*3:]
	}
	return strings.Join(arr, ".") //将一系列字符串连接为一个字符串，之间用sep来分隔。
}
