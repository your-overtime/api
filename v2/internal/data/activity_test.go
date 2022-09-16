package data_test

import (
	"testing"
	"time"

	"github.com/your-overtime/api/v2/internal/data"
	"github.com/your-overtime/api/v2/pkg"
	"github.com/your-overtime/api/v2/tests"
)

func TestSaveActivty(t *testing.T) {
	db := tests.SetupDb(t)
	insert := createActivity("2021-07-02 08:47")
	if err := db.SaveActivity(&insert); err != nil {
		t.Fatal(err)
	}

	if insert.ID == 0 {
		t.Fatal("expected auto increment id after insert")
	}
	if insert.ID != 1 {
		t.Fatal("expected first inserted activity has ID 1")
	}
}

func TestSaveActivtyWithoutDescriptionErr(t *testing.T) {
	db := tests.SetupDb(t)
	insert := createActivity("2021-07-02 08:06")
	insert.Description = ""
	err := db.SaveActivity(&insert)
	if err != pkg.ErrEmptyDescriptionNotAllowed {
		t.Fatal("empty description should return err")
	}
}

func TestGetActivty(t *testing.T) {
	db := tests.SetupDb(t)
	expected := createActivity("2021-07-02 08:56")
	db.SaveActivity(&expected)

	actual, err := db.GetActivity(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !tests.Equals(t, actual, expected) {
		t.Fatalf("%v does not equal %v", actual, expected)
	}
}

func TestGetRunningActivtyByUserID(t *testing.T) {
	db := tests.SetupDb(t)
	notRunning := createActivityWithEnd("2021-07-05 08:00", "2021-07-05 10:00")
	db.SaveActivity(&notRunning)

	running, err := db.GetRunningActivityByUserID(1)
	if err == nil {
		t.Fatalf("expect %v error but got nil", pkg.ErrNoActivityIsRunning)
	}
	if running != nil {
		t.Fatal("expected running activty is nil")
	}
	if err != nil && err != pkg.ErrNoActivityIsRunning {
		t.Fatalf("expected error to be %v but got %v", pkg.ErrNoActivityIsRunning, err)
	}

	isRunning := createActivity("2021-07-05 11:00")
	db.SaveActivity(&isRunning)

	running, err = db.GetRunningActivityByUserID(1)
	if err != nil {
		t.Fatalf("expect nil error but got %v", err)
	}
	if running == nil {
		t.Fatalf("expected %v activity but got nil", isRunning)
	}
	if !tests.Equals(t, running, isRunning) {
		t.Fatalf("expected %v equals %v", running, isRunning)
	}

}

func TestGetActivtyBetweenStartAndEnd(t *testing.T) {
	start, _ := time.Parse("2006-01-02 15:04", "2021-06-28 00:00")
	end, _ := time.Parse("2006-01-02 15:04", "2021-06-28 23:59")

	a1 := createActivityWithEnd("2021-06-28 08:00", "2021-06-28 16:00")
	a2 := createActivityWithEnd("2021-06-29 08:00", "2021-06-29 16:00")
	a3 := createActivityWithEnd("2021-06-27 08:00", "2021-06-27 16:00")
	db := tests.SetupDb(t)
	if err := db.SaveActivity(&a1); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveActivity(&a2); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveActivity(&a3); err != nil {
		t.Fatal(err)
	}

	list, err := db.GetActivitiesBetweenStartAndEnd(start, end, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("expected result len is 1 not %v", len(list))
	}

	if !tests.Equals(t, list[0], a1) {
		t.Fatalf("expected %v equals %v", list[0], a1)
	}
}

func createActivity(start string) data.ActivityDB {
	startTime, _ := time.Parse("2006-01-02 15:04", start)
	return data.ActivityDB{
		Activity: pkg.Activity{
			InputActivity: pkg.InputActivity{
				Start:       &startTime,
				Description: "testing",
			},
			UserID: 1,
		},
	}
}

func createActivityWithEnd(start, end string) data.ActivityDB {
	startTime, _ := time.Parse("2006-01-02 15:04", start)
	endTime, _ := time.Parse("2006-01-02 15:04", end)
	return data.ActivityDB{
		Activity: pkg.Activity{
			InputActivity: pkg.InputActivity{
				Start:       &startTime,
				End:         &endTime,
				Description: "testing",
			},
			UserID: 1,
		},
	}
}
