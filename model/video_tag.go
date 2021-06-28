package model

import (
	"fmt"
	"sprider_sex/utils"
)

type VideoTag struct {
	Id       int `gorm:"primary_key"`
	Video_id int `json:"video_id" gorm:"comment:'视频id'"`
	Tag_id   int `json:"tag_id" gorm:"comment:'标签id'"`
}

func CreateVideoTag(video_id, tag_id int) error {
	video_tag := VideoTag{
		Video_id: video_id,
		Tag_id:   tag_id,
	}
	err := DB.Create(&video_tag).Error
	return err
}

func CheckVideoTag(video_id, tag_id int) (map[string]string, error) {

	redisKey := "videos_tag:server_id_web_name:" + fmt.Sprintf("%d_%d", video_id, tag_id)
	//优先查询redis 拿map
	dMap, err := Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var videos_tag VideoTag

		err = DB.Table("video_tag").Where("video_id = ? and tag_id = ?", video_id, tag_id).First(&videos_tag).Error
		if err == nil && videos_tag.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val = utils.StructToMap(videos_tag)

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
