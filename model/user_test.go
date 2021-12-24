package model

import (
	"context"
	"testing"
)

func TestUserDao_InsertNewUser(t *testing.T) {
	dao := NewUserDao()
	err := dao.InsertNewUser(context.Background(), &User{Username: "z00"})
	t.Log(err)
}
