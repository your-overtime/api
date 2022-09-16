package data_test

import (
	"testing"

	"github.com/your-overtime/api/v2/internal/data"
	"github.com/your-overtime/api/v2/pkg"
	"github.com/your-overtime/api/v2/tests"
)

func TestSaveUser(t *testing.T) {
	db := tests.SetupDb(t)

	newUser := data.UserDB{
		User: pkg.User{
			Name:                     "Name",
			Surname:                  "Surname",
			Login:                    "login",
			Password:                 "supersecurepassword",
			WeekWorkingTimeInMinutes: 1280,
			NumWorkingDays:           3,
		},
	}
	if err := db.SaveUser(&newUser); err != nil {
		t.Fatal(err)
	}
	if newUser.ID == 0 {
		t.Fatal("expected auto increment id after insert")
	}

	if newUser.ID != 1 {
		t.Fatalf("expected first inserted user has ID 1 but got %d", newUser.ID)
	}
}

func TestGetUser(t *testing.T) {
	db := tests.SetupDb(t)
	user := data.UserDB{

		User: pkg.User{
			Name:                     "Name",
			Surname:                  "Surname",
			Login:                    "login",
			Password:                 "supersecurepassword",
			WeekWorkingTimeInMinutes: 1280,
			NumWorkingDays:           3,
		},
	}
	db.SaveUser(&user)

	actual, err := db.GetUser(1)
	if err != nil {
		t.Fatal(err)
	}
	if !tests.Equals(t, actual.ID, user.ID) && !tests.Equals(t, actual.Login, user.Login) {
		t.Fatalf("%v does not equal %v", actual, user)
	}

	actual, err = db.GetUserByLogin("login")
	if err != nil {
		t.Fatal(err)
	}
	if !tests.Equals(t, actual.ID, user.ID) && !tests.Equals(t, actual.Login, user.Login) {
		t.Fatalf("%v does not equal %v", actual, user)
	}
}

func TestUserTokens(t *testing.T) {
	db := tests.SetupDb(t)
	user := data.UserDB{
		User: pkg.User{
			Name:                     "Name",
			Surname:                  "Surname",
			Login:                    "login",
			Password:                 "supersecurepassword",
			WeekWorkingTimeInMinutes: 1280,
			NumWorkingDays:           3,
		},
	}
	db.SaveUser(&user)
	tokens, err := db.GetTokens(user)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 0 {
		t.Fatalf("expected no tokens for user but got tokens with len %d", len(tokens))
	}

	token := data.TokenDB{
		Token: pkg.Token{
			UserID: 1,
			InputToken: pkg.InputToken{
				Name: "nice token name",
			},
			Token: "asd",
		},
	}
	db.Conn.Create(&token)

	tokens, err = db.GetTokens(user)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token for user but got tokens with len %d", len(tokens))
	}
	if !tests.Equals(t, token, tokens[0]) {
		t.Fatalf("%v does not equal %v", token, tokens[0])
	}

	actual, err := db.GetUserByToken("asd")
	if err != nil {
		t.Fatal(err)
	}

	if !tests.Equals(t, actual.ID, user.ID) && !tests.Equals(t, actual.Login, user.Login) {
		t.Fatalf("expected user by token %v equals %v", actual, user)
	}
}
