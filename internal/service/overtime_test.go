package service_test

import (
	"testing"

	"github.com/your-overtime/api/internal/service"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/tests"
)

var e = pkg.User{
	Name:                     "Dieter",
	Surname:                  "Tester",
	Login:                    "dieter",
	Password:                 "secret",
	WeekWorkingTimeInMinutes: 1920,
	NumWorkingDays:           5,
}

func setUp(t *testing.T) (pkg.OvertimeService, *pkg.User) {
	db := tests.SetupDb(t)
	s := service.Init(&db)

	ePtr, err := s.SaveUser(e, "")
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	return s, ePtr
}

func TestOverviewDynamicWOrkinkdays(t *testing.T) {
	s, ePtr := setUp(t)
	o, err := s.CalcOverview(*ePtr, tests.ParseDayTime("2021-01-08 23:59"))
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

	start := tests.ParseDayTime("2021-01-09 07:04")
	end := tests.ParseDayTime("2021-01-09 08:08")
	_, err = s.AddActivity(pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start:       &start,
			End:         &end,
			Description: "Tests",
		},
		UserID: ePtr.ID,
	}, *ePtr)

	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	o, err = s.CalcOverview(*ePtr, tests.ParseDayTime("2021-01-09 23:59"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	if o.OvertimeThisDayInMinutes != -320 {
		t.Error("expect -320 but got ", o.OvertimeThisDayInMinutes)
	}

	if o.ActiveTimeThisDayInMinutes != 64 {
		t.Error("expect 64 but got ", o.ActiveTimeThisDayInMinutes)
	}

	// 384 * 4 - 64 = 1472
	if o.OvertimeThisWeekInMinutes != -1472 {
		t.Error("expect -1472 but got ", o.OvertimeThisWeekInMinutes)
	}

	if o.ActiveTimeThisWeekInMinutes != 64 {
		t.Error("expect 64 but got ", o.ActiveTimeThisWeekInMinutes)
	}

	if o.OvertimeThisMonthInMinutes != -2240 {
		t.Error("expect -2250 but got ", o.OvertimeThisMonthInMinutes)
	}

	if o.ActiveTimeThisMonthInMinutes != 64 {
		t.Error("expect 64 but got ", o.ActiveTimeThisMonthInMinutes)
	}

	if o.OvertimeThisYearInMinutes != -2240 {
		t.Error("expect -2240 but got ", o.OvertimeThisYearInMinutes)
	}

	if o.ActiveTimeThisYearInMinutes != 64 {
		t.Error("expect 64 but got ", o.ActiveTimeThisYearInMinutes)
	}

	if o.ActiveActivity != nil {
		t.Error("expect nil but got ", o.ActiveActivity)
	}
}

func TestOverviewStatic(t *testing.T) {
	s, ePtr := setUp(t)
	ePtr.WorkingDays = "Monday,Tuesday,Wednesday,Thursday,Friday"
	o, err := s.CalcOverview(*ePtr, tests.ParseDayTime("2021-01-08 23:59"))
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

	// 384 * 5 = 1920
	if o.OvertimeThisWeekInMinutes != -1920 {
		t.Error("expect -1920 but got ", o.OvertimeThisWeekInMinutes)
	}

	// 6 workingdays = 6 * 384 = 2304
	if o.OvertimeThisMonthInMinutes != -2304 {
		t.Error("expect -1920 but got ", o.OvertimeThisMonthInMinutes)
	}

	if o.OvertimeThisYearInMinutes != -2304 {
		t.Error("expect -2304 but got ", o.OvertimeThisYearInMinutes)
	}

	if o.ActiveActivity != nil {
		t.Error("expect nil but got ", o.ActiveActivity)
	}

	start := tests.ParseDayTime("2021-01-09 07:04")
	end := tests.ParseDayTime("2021-01-09 08:08")
	_, err = s.AddActivity(pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start:       &start,
			End:         &end,
			Description: "Tests",
		},
		UserID: ePtr.ID,
	}, *ePtr)

	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	o, err = s.CalcOverview(*ePtr, tests.ParseDayTime("2021-01-09 23:59"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}
	// the 2021-01-09 is a Saturday and not in "Monday,Tuesday,Wednesday,Thursday,Friday"
	if o.OvertimeThisDayInMinutes != 64 {
		t.Error("expect 64 but got ", o.OvertimeThisDayInMinutes)
	}

	if o.ActiveTimeThisDayInMinutes != 64 {
		t.Error("expect 64 but got ", o.ActiveTimeThisDayInMinutes)
	}

	// 384 * 5 - 64 = 1856
	if o.OvertimeThisWeekInMinutes != -1856 {
		t.Error("expect -1856 but got ", o.OvertimeThisWeekInMinutes)
	}

	if o.ActiveTimeThisWeekInMinutes != 64 {
		t.Error("expect 64 but got ", o.ActiveTimeThisWeekInMinutes)
	}

	if o.OvertimeThisMonthInMinutes != -2240 {
		t.Error("expect -2240 but got ", o.OvertimeThisMonthInMinutes)
	}

	if o.ActiveTimeThisMonthInMinutes != 64 {
		t.Error("expect 64 but got ", o.ActiveTimeThisMonthInMinutes)
	}

	if o.OvertimeThisYearInMinutes != -2240 {
		t.Error("expect -2240 but got ", o.OvertimeThisYearInMinutes)
	}

	if o.ActiveTimeThisYearInMinutes != 64 {
		t.Error("expect 64 but got ", o.ActiveTimeThisYearInMinutes)
	}

	if o.ActiveActivity != nil {
		t.Error("expect nil but got ", o.ActiveActivity)
	}
}
