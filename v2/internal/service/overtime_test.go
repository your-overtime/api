package service_test

import (
	"testing"

	"github.com/your-overtime/api/v2/internal/data"
	"github.com/your-overtime/api/v2/internal/service"
	"github.com/your-overtime/api/v2/pkg"
	"github.com/your-overtime/api/v2/pkg/utils"
	"github.com/your-overtime/api/v2/tests"
)

var e = pkg.User{
	Name:                     "Dieter",
	Surname:                  "Tester",
	Login:                    "dieter",
	Password:                 "secret",
	WeekWorkingTimeInMinutes: 1920,
	NumWorkingDays:           5,
}

func setUp(t *testing.T) (*service.Service, *pkg.User) {
	db := tests.SetupDb(t)
	err := db.SaveUser(&data.UserDB{User: e})
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}
	s := service.Init(db).GetOrCreateInstanceForUser(&e)

	actualService, ok := s.(*service.Service)
	if !ok {
		t.Fatal("wrong service implementation")
	}

	return actualService, &e
}

func TestOverviewDynamicWorkinkdays(t *testing.T) {
	s, ePtr := setUp(t)
	o, err := s.CalcOverview(tests.ParseDayTime("2021-01-08 23:59"))
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
	})

	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	o, err = s.CalcOverview(tests.ParseDayTime("2021-01-09 23:59"))
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

	o, err = s.CalcOverview(tests.ParseDayTime("2021-01-12 23:59"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	if o.OvertimeThisDayInMinutes != 0 {
		t.Error("expect 0 but got ", o.OvertimeThisDayInMinutes)
	}

	if o.ActiveTimeThisDayInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisDayInMinutes)
	}

	start = tests.ParseDayTime("2021-01-13 07:04")
	end = tests.ParseDayTime("2021-01-13 10:08")
	_, err = s.AddActivity(pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start:       &start,
			End:         &end,
			Description: "Tests",
		},
		UserID: ePtr.ID,
	})

	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	o, err = s.CalcOverview(tests.ParseDayTime("2021-01-13 23:59"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	if o.OvertimeThisDayInMinutes != -200 {
		t.Error("expect -200 but got ", o.OvertimeThisDayInMinutes)
	}

	if o.ActiveTimeThisDayInMinutes != 184 {
		t.Error("expect 184 but got ", o.ActiveTimeThisDayInMinutes)
	}

	// 384 - 184 = -200
	if o.OvertimeThisWeekInMinutes != -200 {
		t.Error("expect -1472 but got ", o.OvertimeThisWeekInMinutes)
	}

	if o.ActiveTimeThisWeekInMinutes != 184 {
		t.Error("expect 64 but got ", o.ActiveTimeThisWeekInMinutes)
	}

	if o.OvertimeThisMonthInMinutes != -2824 {
		t.Error("expect -2824 but got ", o.OvertimeThisMonthInMinutes)
	}

	if o.ActiveTimeThisMonthInMinutes != 248 {
		t.Error("expect 248 but got ", o.ActiveTimeThisMonthInMinutes)
	}

	if o.OvertimeThisYearInMinutes != -2824 {
		t.Error("expect -2824 but got ", o.OvertimeThisYearInMinutes)
	}

	if o.ActiveTimeThisYearInMinutes != 248 {
		t.Error("expect 248 but got ", o.ActiveTimeThisYearInMinutes)
	}

	if o.ActiveActivity != nil {
		t.Error("expect nil but got ", o.ActiveActivity)
	}

	o, err = s.CalcOverview(tests.ParseDayTime("2021-01-16 23:59"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	// this day must be a working day, otherwise it would not be possible to work the 5 working days
	if o.OvertimeThisDayInMinutes != -384 {
		t.Error("expect -384 but got ", o.OvertimeThisDayInMinutes)
	}

	if o.ActiveTimeThisDayInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisDayInMinutes)
	}
}

func TestOverviewStatic(t *testing.T) {
	s, ePtr := setUp(t)
	ePtr.WorkingDays = "Monday,Tuesday,Wednesday,Thursday,Friday"

	o, err := s.CalcOverview(tests.ParseDayTime("2021-01-08 23:59"))
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
	})

	if err != nil {
		t.Fatal("expect no error but got ", err)
	}

	o, err = s.CalcOverview(tests.ParseDayTime("2021-01-09 23:59"))
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

func TestOverviewStaticWithHolidays(t *testing.T) {
	s, ePtr := setUp(t)
	ePtr.WorkingDays = "Monday,Tuesday,Wednesday,Thursday,Friday"

	o, err := s.CalcOverview(tests.ParseDayTime("2021-01-07 23:59"))
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
	if o.OvertimeThisWeekInMinutes != -1536 {
		t.Error("expect -1536 but got ", o.OvertimeThisWeekInMinutes)
	}

	// 6 workingdays = 6 * 384 = 2304
	if o.OvertimeThisMonthInMinutes != -1920 {
		t.Error("expect -1920 but got ", o.OvertimeThisMonthInMinutes)
	}

	if o.OvertimeThisYearInMinutes != -1920 {
		t.Error("expect -1920 but got ", o.OvertimeThisYearInMinutes)
	}

	if o.ActiveActivity != nil {
		t.Error("expect nil but got ", o.ActiveActivity)
	}

	s.AddHoliday(pkg.Holiday{
		InputHoliday: pkg.InputHoliday{
			Start:       utils.DayStart(tests.ParseDay("2021-01-08")),
			End:         utils.DayEnd(tests.ParseDay("2021-01-08")),
			Description: "Free",
			Type:        pkg.HolidayTypeFree,
		},
	})

	o, err = s.CalcOverview(tests.ParseDayTime("2021-01-08 23:59"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}
	// the 2021-01-09 is a Saturday and not in "Monday,Tuesday,Wednesday,Thursday,Friday"
	if o.OvertimeThisDayInMinutes != 0 {
		t.Error("expect 0 but got ", o.OvertimeThisDayInMinutes)
	}

	if o.ActiveTimeThisDayInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisDayInMinutes)
	}

	// 384 * 5 - 64 = 1856
	if o.OvertimeThisWeekInMinutes != -1536 {
		t.Error("expect -1536 but got ", o.OvertimeThisWeekInMinutes)
	}

	if o.ActiveTimeThisWeekInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisWeekInMinutes)
	}

	if o.OvertimeThisMonthInMinutes != -1920 {
		t.Error("expect -1920 but got ", o.OvertimeThisMonthInMinutes)
	}

	if o.ActiveTimeThisMonthInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisMonthInMinutes)
	}

	if o.OvertimeThisYearInMinutes != -1920 {
		t.Error("expect -1929 but got ", o.OvertimeThisYearInMinutes)
	}

	if o.ActiveTimeThisYearInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisYearInMinutes)
	}

	if o.ActiveActivity != nil {
		t.Error("expect nil but got ", o.ActiveActivity)
	}

	s.AddHoliday(pkg.Holiday{
		InputHoliday: pkg.InputHoliday{
			Start:       utils.DayStart(tests.ParseDay("2021-01-13")),
			End:         utils.DayEnd(tests.ParseDay("2021-01-14")),
			Description: "Free",
			Type:        pkg.HolidayTypeFree,
		},
	})

	o, err = s.CalcOverview(tests.ParseDayTime("2021-01-13 23:59"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}
	// the 2021-01-09 is a Saturday and not in "Monday,Tuesday,Wednesday,Thursday,Friday"
	if o.OvertimeThisDayInMinutes != 0 {
		t.Error("expect 0 but got ", o.OvertimeThisDayInMinutes)
	}
	if o.ActiveTimeThisDayInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisDayInMinutes)
	}

	if o.OvertimeThisWeekInMinutes != -768 {
		t.Error("expect −768 but got ", o.OvertimeThisWeekInMinutes)
	}

	if o.ActiveTimeThisWeekInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisWeekInMinutes)
	}

	// 2 * -384 = −768 - 1920 = −2688
	if o.OvertimeThisMonthInMinutes != -2688 {
		t.Error("expect -1920 but got ", o.OvertimeThisMonthInMinutes)
	}

	if o.ActiveTimeThisMonthInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisMonthInMinutes)
	}

	if o.OvertimeThisYearInMinutes != -2688 {
		t.Error("expect -1929 but got ", o.OvertimeThisYearInMinutes)
	}

	if o.ActiveTimeThisYearInMinutes != 0 {
		t.Error("expect 0 but got ", o.ActiveTimeThisYearInMinutes)
	}

	if o.ActiveActivity != nil {
		t.Error("expect nil but got ", o.ActiveActivity)
	}

	start := tests.ParseDayTime("2021-01-14 12:59")
	end := tests.ParseDayTime("2021-01-14 14:59")
	s.AddActivity(pkg.Activity{InputActivity: pkg.InputActivity{
		Start:       &start,
		End:         &end,
		Description: "Working",
	}})

	o, err = s.CalcOverview(tests.ParseDayTime("2021-01-14 23:59"))
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}
	// the 2021-01-09 is a Saturday and not in "Monday,Tuesday,Wednesday,Thursday,Friday"
	if o.OvertimeThisDayInMinutes != 120 {
		t.Error("expect 0 but got ", o.OvertimeThisDayInMinutes)
	}

	if o.ActiveTimeThisDayInMinutes != 120 {
		t.Error("expect 0 but got ", o.ActiveTimeThisDayInMinutes)
	}

	if o.OvertimeThisWeekInMinutes != -648 {
		t.Error("expect −648 but got ", o.OvertimeThisWeekInMinutes)
	}

	if o.ActiveTimeThisWeekInMinutes != 120 {
		t.Error("expect 120 but got ", o.ActiveTimeThisWeekInMinutes)
	}

	// 2 * -384 = −768 - 1920 + 120 = −2568
	if o.OvertimeThisMonthInMinutes != -2568 {
		t.Error("expect -2568 but got ", o.OvertimeThisMonthInMinutes)
	}

	if o.ActiveTimeThisMonthInMinutes != 120 {
		t.Error("expect 120 but got ", o.ActiveTimeThisMonthInMinutes)
	}

	if o.OvertimeThisYearInMinutes != -2568 {
		t.Error("expect -2568 but got ", o.OvertimeThisYearInMinutes)
	}

	if o.ActiveTimeThisYearInMinutes != 120 {
		t.Error("expect 120 but got ", o.ActiveTimeThisYearInMinutes)
	}

	if o.ActiveActivity != nil {
		t.Error("expect nil but got ", o.ActiveActivity)
	}
}
