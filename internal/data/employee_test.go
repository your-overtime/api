package data_test

import (
	"testing"

	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/tests"
)

func TestSaveEmployee(t *testing.T) {
	db := tests.SetupDb(t)

	newEmployee := pkg.Employee{
		User: &pkg.User{
			Name:     "Name",
			Surname:  "Surname",
			Login:    "login",
			Password: "supersecurepassword",
		},
		WeekWorkingTimeInMinutes: 1280,
		NumWorkingDays:           3,
	}
	if err := db.SaveEmployee(&newEmployee); err != nil {
		t.Fatal(err)
	}
	if newEmployee.ID == 0 {
		t.Fatal("expected auto increment id after insert")
	}

	if newEmployee.ID != 1 {
		t.Fatalf("expected first inserted employee has ID 1 but got %d", newEmployee.ID)
	}
}

func TestGetEmployee(t *testing.T) {
	db := tests.SetupDb(t)
	employee := pkg.Employee{
		User: &pkg.User{
			Name:     "Name",
			Surname:  "Surname",
			Login:    "login",
			Password: "supersecurepassword",
		},
		WeekWorkingTimeInMinutes: 1280,
		NumWorkingDays:           3,
	}
	db.SaveEmployee(&employee)

	actual, err := db.GetEmployee(1)
	if err != nil {
		t.Fatal(err)
	}
	if !tests.Equals(t, actual, employee) {
		t.Fatalf("%v does not equal %v", actual, employee)
	}

	actual, err = db.GetEmployeeByLogin("login")
	if err != nil {
		t.Fatal(err)
	}
	if !tests.Equals(t, actual, employee) {
		t.Fatalf("%v does not equal %v", actual, employee)
	}
}

func TestEmployeeTokens(t *testing.T) {
	db := tests.SetupDb(t)
	employee := pkg.Employee{
		User: &pkg.User{
			Name:     "Name",
			Surname:  "Surname",
			Login:    "login",
			Password: "supersecurepassword",
		},
		WeekWorkingTimeInMinutes: 1280,
		NumWorkingDays:           3,
	}
	db.SaveEmployee(&employee)
	tokens, err := db.GetTokens(employee)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 0 {
		t.Fatalf("expected no tokens for employee but got tokens with len %d", len(tokens))
	}

	token := pkg.Token{
		UserID: 1,
		Name:   "nice token name",
		Token:  "asd",
	}
	db.Conn.Create(&token)

	tokens, err = db.GetTokens(employee)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token for employee but got tokens with len %d", len(tokens))
	}
	if !tests.Equals(t, token, tokens[0]) {
		t.Fatalf("%v does not equal %v", token, tokens[0])
	}

	actual, err := db.GetEmployeeByToken("asd")
	if err != nil {
		t.Fatal(err)
	}

	if !tests.Equals(t, actual, employee) {
		t.Fatalf("expected employee by token %v equals %v", actual, employee)
	}
}
