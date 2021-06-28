package sprider

import "sprider_sex/sprider/jiucao"

type Sprider_go struct {
	Calls map[string]func()
}

type Sprider_Do interface {
	GetType()
	Download()
}

func NewSprider_go() *Sprider_go {
	return &Sprider_go{Calls: make(map[string]func())}
}

//注册爬虫网站
func (s *Sprider_go) Register() *Sprider_go {

	jiu := jiucao.NewJiuCao("jiucao", "https://jcxx77.com", 1, true, false) //定义 色站爬取入口网址,协和数量,是否开启,是否下载

	s.Calls = map[string]func(){
		"jiucao": jiu.Run,
	}
	return s
}

//开启所有色站爬取
func (s *Sprider_go) Do() {
	for _, run := range s.Calls {
		go run()
	}
}
