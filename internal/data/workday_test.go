package data_test

import (
	"testing"

	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/tests"
)

func TestGetWorkDay(t *testing.T) {

	db := tests.SetupDb(t)

	none, err := db.GetWorkDay(tests.ParseDay("2021-07-05"), 1)
	if none != nil {
		t.Fatalf("expected no workdays but got %v", none)
	}
	if err == nil {
		t.Fatalf("expected error to be nil but got %v", err)
	}
	wd := data.WorkDayDB{
		WorkDay: pkg.WorkDay{
			InputWorkDay: pkg.InputWorkDay{
				Day:        tests.ParseDay("2021-07-05"),
				Overtime:   12,
				ActiveTime: 120,
				UserID:     1,
			},
		},
	}
	db.Conn.Create(&wd)

	actual, err := db.GetWorkDay(tests.ParseDay("2021-07-05"), 1)
	if err != nil {
		t.Fatalf("exptect error to be nil, but got %v", err)
	}
	if actual == nil {
		t.Fatalf("expected actual to be %v bot got nil", wd)
	}
	if !tests.Equals(t, actual, wd) {
		t.Fatalf("expected %v equals %v", actual, wd)
	}
}

func TestDeleteWorkDay(t *testing.T) {
	db := tests.SetupDb(t)
	wd := data.WorkDayDB{
		WorkDay: pkg.WorkDay{
			InputWorkDay: pkg.InputWorkDay{
				Day:        tests.ParseDay("2021-07-05"),
				Overtime:   12,
				ActiveTime: 120,
				UserID:     1,
			},
		},
	}
	db.Conn.Create(&wd)
	doesNotExist := db.DeleteWorkDay(tests.ParseDay("2021-07-04"), 1)
	if doesNotExist != nil {
		t.Fatalf("expected no error but got %v", doesNotExist)
	}

	err := db.DeleteWorkDay(tests.ParseDay("2021-07-05"), 1)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}

func TestGetWorkDaysBetweenStartAndEnd(t *testing.T) {
	db := tests.SetupDb(t)
	wd := data.WorkDayDB{
		WorkDay: pkg.WorkDay{
			InputWorkDay: pkg.InputWorkDay{
				Day:        tests.ParseDay("2021-07-05"),
				Overtime:   12,
				ActiveTime: 120,
				UserID:     1,
			},
		},
	}
	db.Conn.Create(&wd)

	workDayList, err := db.GetWorkDaysBetweenStartAndEnd(
		tests.ParseDay("2021-07-05"), tests.ParseDay("2021-07-05"), 1,
	)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	if len(workDayList) != 1 {
		t.Fatalf("expected workday list has len 1 but got %v", len(workDayList))
	}
	if !tests.Equals(t, workDayList[0], wd) {
		t.Fatalf("expected %v equals %v", workDayList[0], wd)
	}

	workDayList, err = db.GetWorkDaysBetweenStartAndEnd(
		tests.ParseDayTime("2021-07-05 08:00"),
		tests.ParseDayTime("2021-07-05 16:00"),
		1,
	)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	if len(workDayList) != 1 {
		t.Fatalf("expected workday list has len 1 but got %v", len(workDayList))
	}
	if !tests.Equals(t, workDayList[0], wd) {
		t.Fatalf("expected %v equals %v", workDayList[0], wd)
	}

	workDayList, err = db.GetWorkDaysBetweenStartAndEnd(
		tests.ParseDay("2021-06-30"),
		tests.ParseDay("2021-07-04"), 1,
	)

	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	if len(workDayList) != 0 {
		t.Fatalf("expected workday list has len 0 but got %v", len(workDayList))
	}
}
