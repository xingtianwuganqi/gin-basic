package main

import (
	"pet-project/db"
	"pet-project/routers"
	"pet-project/settings"
)

func main() {
	if err := settings.LoadConfig(); err != nil {
		panic(err)
	}
	db.LinkInit()
	r := routers.RegisterRouter()
	r.Run()
}
