package model

import (
	"fmt"
	"sprider_sex/utils"
)

type VideosList struct {
	Id            int    `gorm:"primary_key"`
	Title         string `json:"title" gorm:"comment:'名称'"`
	Description   string `json:"description" gorm:"comment:'描述'"`
	Series_digits string `json:"series_digits" gorm:"comment:'番号'"`
	Status        int    `json:"status" gorm:"comment:'状态'"`
	Video_time    int    `json:"video_time" gorm:"comment:'视频时间'"`
	Img_pc        string `json:"img_pc" gorm:"comment:'图片'"`
	Img_h5        string `json:"img_h5" gorm:"comment:'图片'"`
	Video_url     string `json:"video_url" gorm:"comment:'视频地址'"`
	Video_type    int    `json:"video_type" gorm:"comment:'类型'"`
	Video_class   int    `json:"video_class" gorm:"comment:'类型'"`
	Create_time   int64  `json:"create_time" gorm:"comment:'创建时间'"`
	Server_id     int    `json:"server_id" gorm:"comment:'网站视频id'"`
	Watch_times   int    `json:"watch_times" gorm:"comment:'观看人数'"`
	Pro_num       int    `json:"pro_num" gorm:"comment:'点赞人数'"`
	Update_time   int    `json:"update_time" gorm:"comment:'更新时间'"`
	Web_name      string `json:"web_name" gorm:"comment:'网站名称'"`
}

func CreateVideoList(v VideosList) (rep VideosList, err error) {
	err = DB.Create(&v).Error
	rep = v
	return rep, err
}

func CheckVideosList(server_id int, web_name string) (map[string]string, error) {

	redisKey := "videos_list:server_id_web_name:" + fmt.Sprintf("%d_", server_id) + web_name
	//优先查询redis 拿map
	dMap, err := Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var videos_list VideosList

		err = DB.Table("videos_list").Where("server_id = ? and web_name = ?", server_id, web_name).First(&videos_list).Error
		if err == nil && videos_list.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val = utils.StructToMap(videos_list)

			err = Pool.HMSet(redisKey, val).Err()
			if err != nil {
				return dMap, err
			}

			//新增无序集合 所有的key头存在无序集合里面
			err = Pool.SAdd(ServerInfo.RedisConf.Head, redisKey).Err()
			if err != nil {
				return dMap, err
			}

			dMap, err = Pool.HGetAll(redisKey).Result()
			if err != nil {
				return dMap, err
			}
		}

	}

	return dMap, err
}

//获取视频
func GetVideoListByType(v_type, status, limit int) (v []VideosList, err error) {
	err = DB.Table("videos_list").Where("video_type = ? and status = ?", v_type, status).Find(&v).Limit(limit).Error
	return v, err
}

//更新下载状态和url
func UpdateVideoListById(url, imag_path string, id, video_type, min int) error {

	err := DB.Table("videos_list").Where("id = ?", id).Update(map[string]interface{}{
		"video_url":  url,
		"video_type": video_type, //已经下载
		"img_pc":     imag_path,
		"img_h5":     imag_path,
		"video_time": min,
	}).Error
	return err

}
