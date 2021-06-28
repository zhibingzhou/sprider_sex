package model

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Server struct {
	RedisConf RedisConf `mapstructure:"redis" json:"redis" yaml:"redis"`
	MysqlConf MysqlConf `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
}

var (
	ServerInfo Server
	Viper      *viper.Viper
)

const defaultConfigFile = "config.yaml"

// viper 系统信息读取
func init() {

	v := viper.New()
	v.SetConfigFile(defaultConfigFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&ServerInfo); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&ServerInfo); err != nil {
		fmt.Println(err)
	}

	Viper = v

}

func GetKey(length int) string {
	sec := strconv.FormatInt(time.Now().Unix(), 10)
	redKey := "model_get_key:" + sec
	randLen := length
	exTime := 1
	preId := ""

	if length > 10 {
		randLen = length - 10
		preId = sec
	}
	randStr := ""
	for i := 0; i < 50; i++ {
		randStr = Random("smallnumber", randLen)
		//新增无序集合 所有的key头存在无序集合里面
		res, err := Pool.SAdd(redKey, randStr, exTime).Result()
		if err == nil && res > 0 {
			break
		}
	}

	keyStr := preId + randStr
	return keyStr
}

func Random(param string, length int) string {
	str := ""
	if length < 1 {
		return str
	}
	tmp := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	switch param {
	case "number":
		tmp = "1234567890"
	case "small":
		tmp = "abcdefghijklmnopqrstuvwxyz"
	case "big":
		tmp = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "smallnumber":
		tmp = "1234567890abcdefghijklmnopqrstuvwxyz"
	case "bignumber":
		tmp = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "bigsmall":
		tmp = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	leng := len(tmp)
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		s_ind := ran.Intn(leng)
		str = str + Substr(tmp, s_ind, 1)
	}

	return str
}

/**
*  start：正数 - 在字符串的指定位置开始,超出字符串长度强制把start变为字符串长度
*  负数 - 在从字符串结尾的指定位置开始
*  0 - 在字符串中的第一个字符处开始
*  length:正数 - 从 start 参数所在的位置返回
*  负数 - 从字符串末端返回
 */
func Substr(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	rune_str := []rune(str)
	len_str := len(rune_str)

	if start < 0 {
		start = len_str + start
	}
	if start > len_str {
		start = len_str
	}
	end := start + length
	if end > len_str {
		end = len_str
	}
	if length < 0 {
		end = len_str + length
	}
	if start > end {
		start, end = end, start
	}
	return string(rune_str[start:end])
}
