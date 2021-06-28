package main

import (
	"sprider_sex/router"
	"sprider_sex/sprider"
	"sprider_sex/utils"
)

func main() {
	utils.GVA_LOG.Debug("sprider start!!")
	go sprider.NewSprider_go().Register().Do()
	router.Router.Run(":8081")
}
