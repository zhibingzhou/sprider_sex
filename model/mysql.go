package model

import (
	"fmt"
	"time"
    _ "github.com/go-sql-driver/mysql" 
	"github.com/jinzhu/gorm"
)

type MysqlConf struct {
	Network   string `mapstructure:"network" json:"network" yaml:"network"`
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
	Port      int    `mapstructure:"port" json:"port" yaml:"port"`
	User      string `mapstructure:"user" json:"user" yaml:"user"`
	Pwd       string `mapstructure:"pwd" json:"pwd" yaml:"pwd"`
	Db_name   string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	Life_time string `mapstructure:"life_time" json:"life_time" yaml:"life_time"`
	Max_open  int    `mapstructure:"max_open" json:"max_open" yaml:"max_open"`
	Max_idle  int    `mapstructure:"max_idle" json:"max_idle" yaml:"max_idle"`
}

var DB *gorm.DB

//mysql 初始化
func init() {
	DB = ReloadConfSQL()
}

func ReloadConfSQL() *gorm.DB {

	conn_str := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8",ServerInfo.MysqlConf.User, ServerInfo.MysqlConf.Pwd, ServerInfo.MysqlConf.Network, ServerInfo.MysqlConf.Host, ServerInfo.MysqlConf.Port, ServerInfo.MysqlConf.Db_name)

	fmt.Println(conn_str)
	db, err := gorm.Open("mysql", conn_str)
	if err != nil {
		fmt.Println("conn_str->", conn_str)
		panic(err)
	}

	life_time, _ := time.ParseDuration(ServerInfo.MysqlConf.Life_time)
	//最大生命周期
	db.DB().SetConnMaxLifetime(life_time)
	//连接池的最大打开连接数
	db.DB().SetMaxOpenConns(ServerInfo.MysqlConf.Max_open)
	//连接池的最大空闲连接数
	db.DB().SetMaxIdleConns(ServerInfo.MysqlConf.Max_idle)
	db.SingularTable(true)
	//启用Logger，显示详细日志
	db.LogMode(true)

	// 禁用日志记录器，不显示任何日志
	//db.LogMode(false)
	fmt.Println("连接 mysql 成功")
	return db
}
