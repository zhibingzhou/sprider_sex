package model

type TagList struct {
	Id    int    `gorm:"primary_key"`
	Title string `json:"sort" gorm:"comment:'名称'"`
	Sort  int    `json:"sort" gorm:"comment:'排序'"`
	Hot   int    `json:"sort" gorm:"comment:'热度'"`
}

func AddTagList(title string, sort, hot int) (err error) {
	tag := TagList{
		Title: title,
		Sort:  sort,
		Hot:   hot,
	}
	err = DB.Create(&tag).Error
	return err
}

func CheckTagList(title string) (map[string]string, error) {

	redisKey := "TagList:title:" + title
	//优先查询redis 拿map
	dMap, err := Pool.HGetAll(redisKey).Result()
	if err == nil && len(dMap["id"]) < 1 {
		var taglist TagList

		err = DB.Table("tag_list").Where("title = ?", title).First(&taglist).Error
		if err == nil && taglist.Id > 0 {
			// 查询数据库 得 map
			val := map[string]interface{}{}
			val["id"] = taglist.Id
			val["title"] = taglist.Title

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
