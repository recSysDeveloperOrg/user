package model

import (
	"testing"
	"user/config"
)

func TestMain(m *testing.M) {
	if err := config.InitConfig(config.CfgFileNested); err != nil {
		panic(err)
	}
	if err := InitModel(); err != nil {
		panic(err)
	}
	if code := m.Run(); code != 0 {
		panic(code)
	}
}
