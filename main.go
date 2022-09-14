package main

import (
	"ginblogtest/model"
	"ginblogtest/routes"
)

func main() {
	// fmt.Println(utils.AppMode)
	model.InitDb()
	routes.InitRouter()
}
