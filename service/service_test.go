package service

import (
	"testing"
	"user/config"
	"user/cred"
	"user/model"
)

func TestMain(m *testing.M) {
	if err := config.InitConfig(config.CfgFileNested); err != nil {
		panic(err)
	}
	if err := model.InitModel(); err != nil {
		panic(err)
	}
	if err := cred.InitJwt(); err != nil {
		panic(err)
	}
	if code := m.Run(); code != 0 {
		panic(code)
	}
}
