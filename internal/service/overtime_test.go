package service_test

import (
	"encoding/json"
	"testing"

	"github.com/your-overtime/api/internal/service"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/tests"
)

func TestOverview(t *testing.T) {
	db := tests.SetupDb(t)
	s := service.Init(&db)

	e := pkg.Employee{
		User: &pkg.User{
			Name:     "Dieter",
			Surname:  "Tester",
			Login:    "dieter",
			Password: "secret",
		},
		WeekWorkingTimeInMinutes: 1920,
		NumWorkingDays:           5,
	}

	ePtr, err := s.SaveEmployee(e)
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	_, err = s.AddHoliday(pkg.Holiday{
		Start:       tests.ParseDay("2020-12-28"),
		End:         tests.ParseDay("2020-12-31"),
		Description: "fix first week", // TODO: fix this in CalcOverview function @jasperem
	}, *ePtr)

	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	o, err := s.CalcOverview(*ePtr, tests.ParseDay("2021-01-01"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	x, _ := json.MarshalIndent(&o, " ", "")
	t.Log(string(x))
}
