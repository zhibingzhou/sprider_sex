package jiucao

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sprider_sex/model"
	"sprider_sex/utils"
	"strconv"
	"strings"
)

type M3u8WebFile struct {
	Url       string
	VideoUrl  string
	FileName  string
	ImageUrl  string
	NewPath   string
	ImagePath string
	Id        string
}

func NewM3u8File(url, img_url string, id int) *M3u8WebFile {
	var m M3u8WebFile
	resp, err := http.Get(url)
	if err != nil {
		utils.GVA_LOG.Error(err)
		return &m
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	fileInfo := strings.Split(string(body), "\n")

	var result string
	for _, line := range fileInfo {
		if strings.Contains(line, ".m3u8") {
			fmt.Println(line)
			result = line
			break
		}
	}

	m.NewPath = fmt.Sprintf("./static/%s/video/", JiuCaoSex.Name)
	m.ImagePath = fmt.Sprintf("./static/%s/img/", JiuCaoSex.Name)
	m.Url = strings.Replace(url, "index.m3u8", result, -1)
	m.VideoUrl = strings.Replace(m.Url, "index.m3u8", "", -1)
	m.ImageUrl = img_url
	m.Id = strconv.Itoa(id)
	return &m
}

//通过链接 拿到 m3u8 文件
func (m *M3u8WebFile) GetM3u8File() (string, int) {
	// Get the data
	timeSecond := 0.00
	resp, err := http.Get(m.Url)
	if err != nil {
		utils.GVA_LOG.Error(err)
		return m.Url, int(timeSecond)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	video_key := model.GetKey(5)
	m.FileName = m.NewPath + m.Id + "_" + video_key + "/"

	fileInfo := strings.Split(string(body), "\n")
	var filetext = ""
	keyurl := strings.Replace(m.Url, "index.m3u8", "key.key", -1)
	for _, line := range fileInfo {
		if strings.Contains(line, ".image") {
			fmt.Println(line)
			//下载.image文件
			DownLoadFromM3u8(m.FileName, m.VideoUrl+line, "")
		} else if strings.Contains(line, ".ts") {
			fmt.Println(line)
			//下载.ts文件
			DownLoadFromM3u8(m.FileName, m.VideoUrl+line, GetM3u8Key(keyurl))
		}

		if strings.Contains(line, "#EXTINF:") {
			searray := strings.Split(line, ":")
			if len(searray) >= 2 {
				second, _ := strconv.ParseFloat(strings.Replace(searray[1], ",", "", -1), 64)
				timeSecond += second
			}
		}

		if strings.Contains(line, "key.key") {
			continue
		}
		filetext = filetext + line + "\n"
	}

	//重写m3u8文件
	WriteToFile(m.FileName+video_key+".m3u8", filetext)

	if timeSecond > 0 {
		timeSecond = timeSecond / 60
	}
	return m.FileName + video_key + ".m3u8", int(timeSecond)
}

//通过链接 拿到 图片文件
func (m *M3u8WebFile) GetImgeFile() string {
	// Get the data
	resp, err := http.Get(m.ImageUrl)
	if err != nil {
		utils.GVA_LOG.Error(err)
		return m.ImageUrl
	}
	defer resp.Body.Close()

	video_key := model.GetKey(5)
	m.FileName = m.ImagePath + m.Id + "_" + video_key + "/"

	urlarray := strings.Split(m.ImageUrl, "/")
	if len(urlarray) < 2 {
		return m.ImageUrl
	}

	if ok := Exists(m.FileName); !ok { //创建文件夹
		err := os.MkdirAll(m.FileName, os.ModePerm)
		if err != nil {
			return m.ImageUrl
		}
	}

	// 创建一个文件用于保存
	out, err := os.Create(m.FileName + urlarray[len(urlarray)-1])
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	return m.FileName + urlarray[len(urlarray)-1]
}

func Download(url, img_url string, id int) {
	file := NewM3u8File(url, img_url, id)
	m3u8path, min := file.GetM3u8File()
	impath := file.GetImgeFile()

	err := model.UpdateVideoListById(m3u8path, impath, id, 1, min)
	if err != nil {
		utils.GVA_LOG.Error(id, "已经下载完成,更新报错", err)
	}
	fmt.Println(min, impath, m3u8path)
}

//通过链接下载 .ts文件或者.image 文件
func DownLoadFromM3u8(filepath, url, key string) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		utils.GVA_LOG.Error(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.GVA_LOG.Error(err)
		return
	}

	fmt.Println(url, key)
	result, _ := AesDecrypt(body, []byte(key))

	arrayU := strings.Split(url, "/")

	if len(arrayU) < 2 {
		return
	}

	fineName := arrayU[len(arrayU)-1]
	if ok := Exists(filepath); !ok { //创建文件夹
		err := os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			return
		}
	}

	// 创建一个文件用于保存
	out, err := os.Create(filepath + fineName)
	if err != nil {
		utils.GVA_LOG.Error(err)
		return
	}
	defer out.Close()

	io.WriteString(out, string(result))
}

func GetM3u8Key(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		utils.GVA_LOG.Error(err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.GVA_LOG.Error(err)
		return ""
	}
	fmt.Println(string(body))
	return string(body)

}

//aes 解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	//origData = PKCS5UnPadding(origData)
	return origData, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//切换路径
func ChangeUrlHead(url string) string {

	result := ""

	arrayU := strings.Split(url, "/")

	if len(arrayU) < 2 {
		return result
	}

	result = arrayU[len(arrayU)-1]

	return result
}

func WriteToFile(filePath, fileText string) {
	//创建一个新文件，写入内容 5 句 “http://c.biancheng.net/golang/”
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(fileText)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

// ffmpeg -i "concat:segment-1-v1-a1.ts|segment-2-v1-a1.ts" -acodec copy -vcodec copy -absf aac_adtstoasc output.mp4
func makeMp4(str, savedir, name string) {
	binary, lookErr := exec.LookPath("ffmpeg")
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{
		"-i",
		fmt.Sprintf("concat:%s", str),
		"-acodec",
		"copy",
		"-vcodec",
		"copy",
		"-absf",
		"aac_adtstoasc",
		fmt.Sprintf("%s/%s.mp4", savedir, name),
	}

	cmd := exec.Command(binary, args...)
	r, err := cmd.Output()

	fmt.Println(err)
	fmt.Println(string(r))
}
