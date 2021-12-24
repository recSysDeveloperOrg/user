package main

import (
	"user/config"
	"user/cred"
	"user/model"
)

func main() {
	if err := config.InitConfig(config.CfgFileMain); err != nil {
		panic(err)
	}
	if err := model.InitModel(); err != nil {
		panic(err)
	}
	if err := cred.InitJwt(); err != nil {
		panic(err)
	}
}
