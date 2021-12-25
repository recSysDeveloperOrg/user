package service

import (
	"math/rand"
	"testing"
	"time"
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

func testRandomString(size int) string {
	rand.Seed(time.Now().Unix())
	s := ""
	for i := 0; i < size; i++ {
		s += string(rune('a' + rand.Intn(26)))
	}

	return s
}
