package jiucao

import (
	"sprider_sex/channel_go"
	"sprider_sex/model"
	"sprider_sex/utils"
	"sync"
	"time"
)

type JiuCao struct {
	Name     string
	Url      string
	Count    int  //分配的协程数量
	Open     bool //是否开启
	Download bool //是否开启下载
}

var JiuCaoSex JiuCao

//爬的初始网页
func NewJiuCao(name string, url string, count int, open, download bool) JiuCao {
	JiuCaoSex = JiuCao{Name: name, Url: url, Count: count, Open: open, Download: download}
	return JiuCaoSex
}

func (j JiuCao) Run() {

	if !j.Open { //如果未开启，直接返回
		return
	}

	j.RunDownload() //开启下载

	channel_go.Manager = channel_go.NewManagerChannel(j.Count)
	channel_go.Manager.Run()

	j.GetTypeFromUrl() //拿首页的所有类型链接

	go func() {
		for {

			re := <-channel_go.Manager.Result
			if request, ok := re.(Jiu_PareResult); ok {
				for _, value := range request.Requests {
					channel_go.Manager.Request <- value
				}
			}

		}
	}()

	for {

		sex_type_url, err := j.GetTypeFromRedis() //取任意一种类型链接
		if err != nil {
			utils.GVA_LOG.Error()
		}

		if len(sex_type_url) < 1 {
			utils.GVA_LOG.Debug(" 类型视频 全部爬完！！")
			err := j.Delcash()
			if err != nil {
				return
			}
			j.GetTypeFromUrl()           //拿首页的所有类型链接
			time.Sleep(time.Second * 10) //30 分钟后再爬
			continue
		}

		channel_go.Manager.Request <- Jiu_Request{Url: sex_type_url[0], PareFunc: GetFilmUrl} //拿到当前链接下的所有地址

		<-channel_go.Manager.EndJob

	}

}

//运行下载
func (j JiuCao) RunDownload() {

	if j.Download == false {
		return
	}

	var group sync.WaitGroup

	for {

		var video_array []model.VideosList
		video_array, _ = model.GetVideoListByType(0, 0, j.Count)

		if len(video_array) < 1 {
			time.Sleep(30 * time.Minute) //30 分钟后再下载
			continue
		}

		for _, value := range video_array {
			go func() {
				group.Add(1)
				Download(value.Video_url, value.Img_pc, value.Id)
				group.Done()
			}()
		}

		group.Wait()

	}

}
