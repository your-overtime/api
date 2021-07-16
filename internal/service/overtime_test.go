package service_test

import (
	"encoding/json"
	"testing"

	"github.com/your-overtime/api/internal/service"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/tests"
)

var e = pkg.Employee{
	User: &pkg.User{
		Name:     "Dieter",
		Surname:  "Tester",
		Login:    "dieter",
		Password: "secret",
	},
	WeekWorkingTimeInMinutes: 1920,
	NumWorkingDays:           5,
}

func setUp(t *testing.T) (pkg.OvertimeService, *pkg.Employee) {
	db := tests.SetupDb(t)
	s := service.Init(&db)

	ePtr, err := s.SaveEmployee(e, "")
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	return s, ePtr
}

func TestOverview(t *testing.T) {
	s, ePtr := setUp(t)
	o, err := s.CalcOverview(*ePtr, tests.ParseDay("2021-01-08"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	if o.ActiveTimeThisDayInMinutes != 0 ||
		o.ActiveTimeThisWeekInMinutes != 0 ||
		o.ActiveTimeThisMonthInMinutes != 0 ||
		o.ActiveTimeThisYearInMinutes != 0 {
		t.Error("expect that active time is 0")
	}

	if o.OvertimeThisDayInMinutes != -384 {
		t.Error("expect -384 but got ", o.OvertimeThisDayInMinutes)
	}

	// 384 * 3 = 1152
	if o.OvertimeThisWeekInMinutes != -1152 {
		t.Error("expect -1152 but got ", o.OvertimeThisWeekInMinutes)
	}

	// 7 full workdays without activity => -1920
	if o.OvertimeThisMonthInMinutes != -1920 {
		t.Error("expect -1920 but got ", o.OvertimeThisMonthInMinutes)
	}

	if o.OvertimeThisYearInMinutes != -1920 {
		t.Error("expect -1920 but got ", o.OvertimeThisYearInMinutes)
	}

	if o.ActiveActivity != nil {
		t.Error("expect nil but got ", o.ActiveActivity)
	}

	x, _ := json.MarshalIndent(&o, " ", "")
	t.Log(string(x))
}
