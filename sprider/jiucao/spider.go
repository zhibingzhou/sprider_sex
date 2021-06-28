package jiucao

import (
	"fmt"
	"regexp"
	"sprider_sex/model"
	"sprider_sex/utils"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis"
)

func (j JiuCao) Delcash() error {
	return model.Pool.Del(j.Name).Err()
}

//从首页获取 类型存入 redis 和 数据库
func (j JiuCao) GetTypeFromUrl() {

	Info, _ := utils.Fetch(j.Url)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(Info)))
	if err != nil {
		return
	}

	dom.Find(".block li").Each(func(i int, t *goquery.Selection) {
		t.Find("a").Each(func(i int, t *goquery.Selection) {
			if t.Text() != "首页" && t.Text() != "小说" {
				herf, _ := t.Attr("href")
				fmt.Println(herf, t.Text())
				j.PutTypeData(herf)
				rep, _ := model.CheckTagList(t.Text()) //查看是否有这个类型
				if len(rep["title"]) < 1 {             //没有则添加
					model.AddTagList(t.Text(), 0, 0)
				}
			}
		})
	})
}

//爬虫去重 ，添加的是有序集合 加到redis
func (j JiuCao) PutTypeData(url string) error {

	// 添加有序集合 插入成功为1 插入失败为0,去重用
	value, err := model.Pool.ZAdd(j.Name, redis.Z{Score: 10, Member: url}).Result()

	if err != nil {
		utils.GVA_LOG.Debug(j.Name+"PuttypeData", err)
	}

	if value == 1 { //说明没有这个key
		fmt.Println(j.Name+"类型url添加", url)
		onlyid := model.ServerInfo.RedisConf.Head + model.GetKey(16)
		//存对应的data到有序集合
		value, err := model.Pool.ZAdd(j.Name+"_data", redis.Z{Score: 10, Member: onlyid}).Result()
		fmt.Println(value, err)
		ma := map[string]interface{}{}
		ma["url"] = url
		//再存入map参数
		err = model.Pool.HMSet(onlyid, ma).Err()
		if err != nil {
			fmt.Println(j.Name+"PuttypeData", err)
		}
	}
	return err
}

//从redis获取类型
func (j JiuCao) GetTypeFromRedis() (url_list []string, err error) {

	//设置最大和最小值  返回有序集合的所有元素和分数
	vals, err := model.Pool.ZRangeByScoreWithScores(j.Name+"_data", redis.ZRangeBy{
		Min:    "0",
		Max:    "50",
		Offset: 0,
		Count:  1,
	}).Result()

	for _, value := range vals {
		key := value.Member.(string)
		dMap, err := model.Pool.HGetAll(key).Result()
		if err != nil {
			return url_list, err
		}
		url_list = append(url_list, dMap["url"])
		model.Pool.Del(key).Err()
		//删除集合中的一个指定元素
		model.Pool.ZRem(j.Name+"_data", key)
	}

	return url_list, err
}

func GetFilmUrl(url string, i int) Jiu_PareResult {
	var rep Jiu_PareResult
	EndPage := false
	for {

		if EndPage {
			fmt.Println("爬完当前类型")
			return rep
		}
		EndPage = true
		fmt.Println("爬取页面：" + JiuCaoSex.Url + url)
		Info, _ := utils.Fetch(JiuCaoSex.Url + url)
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(Info)))
		if err != nil {
			return rep
		}

		//拿到类型
		var video_type int
		//拿类型
		dom.Find(".detail_right_div").Each(func(i int, t *goquery.Selection) {

			//判断是否是最后一页
			var next_page []string
			t.Find(".nextPage").Each(func(i int, t *goquery.Selection) {
				t.Find("a").Each(func(i int, t *goquery.Selection) {
					href, ok := t.Attr("href")
					if ok {
						next_page = append(next_page, href)
					}
				})
				t.Find(".select").Each(func(i int, t *goquery.Selection) {
					EndPage = false
				})
			})
			if len(next_page) < 1 {
				return
			}

			url = next_page[len(next_page)-1]

			t.Find("h3 span").Each(func(i int, t *goquery.Selection) {
				rep, _ := model.CheckTagList(t.Text()) //查看是否有这个类型
				if len(rep["title"]) > 0 {
					video_type, _ = strconv.Atoi(rep["id"])
				}
			})

			//拿当前页的信息
			t.Find("ul li").Each(func(i int, t *goquery.Selection) {
				img_src, _ := t.Find("img").Attr("data-original")
				href, _ := t.Find("a").Attr("href")
				create_time := t.Find("i").Text()
				watch_times := t.Find("strong").Text()
				title := ""
				t.Find("p").Each(func(i int, t *goquery.Selection) {
					if i == 1 {
						title = t.Text()
						fmt.Println(title)
					}
				})
				video_list := model.VideosList{}
				video_list.Title = title
				video_list.Img_h5 = img_src
				video_list.Img_pc = img_src
				stamp, _ := time.ParseInLocation("2006-1-2", create_time, time.Local)
				video_list.Create_time = stamp.Unix()
				video_list.Watch_times, _ = strconv.Atoi(strings.Replace(watch_times, "观看", "", -1))
				server_list := strings.Split(href, "/")

				if len(server_list) > 1 {
					server_id, _ := strconv.Atoi(server_list[len(server_list)-1]) //拿视频id
					video_list.Server_id = server_id
					video_list.Web_name = JiuCaoSex.Name
					rep.Requests = append(rep.Requests, Jiu_Request{Url: href, PareFunc: func(s string, i int) Jiu_PareResult {
						return Getinformation(JiuCaoSex.Url+s, i, video_list, video_type)
					}})
				}

			})
		})

	}
}

func Getinformation(url string, i int, v model.VideosList, video_type int) Jiu_PareResult {

	var result Jiu_PareResult

	Info, _ := utils.Fetch(url)
	reage := `playurl = \'([^\']*)\'`
	re := regexp.MustCompile(reage)
	math := re.FindAllStringSubmatch(string(Info), -1)
	if len(math[0]) > 1 {
		v.Video_url = string(math[0][1])
	}
	fmt.Println(v)
	//数据库插入信息
	res, _ := model.CheckVideosList(v.Server_id, JiuCaoSex.Name)

	if len(res["id"]) < 1 { //说明之前没有这个数据
		video, err := model.CreateVideoList(v)
		if err != nil {
			return result
		}
		tag, _ := model.CheckVideoTag(video.Id, video_type)
		if len(tag["id"]) < 1 {
			model.CreateVideoTag(video.Id, video_type)
		}
	}

	return result
}
