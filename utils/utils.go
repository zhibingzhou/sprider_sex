package utils

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func RandSleep(sec int) {
	time.Sleep(time.Duration(rand.Intn(sec)) * time.Second)
}

func ZhToUnicode(sText string) ([]byte, error) {
	textQuoted := strconv.QuoteToASCII(sText)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(textUnquoted)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// 利用反射将化为map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

var ratelimit = time.Tick(200 * time.Millisecond)

//解析url
func Fetch(url string) ([]byte, error) {
	<-ratelimit //等待时间
	re, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if re.StatusCode == http.StatusOK {
		all, err := ioutil.ReadAll(re.Body)
		if err != nil {
			return nil, err
		}
		return all, nil
	}
	defer re.Body.Close()

	return nil, fmt.Errorf("wrong!!")
}

func DeleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "  ", " ", -1)      //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}


