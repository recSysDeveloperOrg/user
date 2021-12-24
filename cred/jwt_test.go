package cred

import "testing"

func TestMain(m *testing.M) {
	if err := InitJwt(); err != nil {
		panic(err)
	}
	if code := m.Run(); code != 0 {
		panic(code)
	}
}

func TestIssueJWT(t *testing.T) {
	userID := "100000000000000000000001"
	token, err := IssueJWT(userID, ExpireDayAccessToken)
	if err != nil {
		t.Fatal(err)
	}

	parsedUserID, err := ParseJWT(token)
	if parsedUserID != userID {
		t.Fatal("not equal")
	}
}
