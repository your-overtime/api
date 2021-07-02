package data_test

import (
	"testing"
	"time"

	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/tests"
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

func TestGetActivty(t *testing.T) {
	db := tests.SetupDb(t)
	expected := createActivity("2021-07-02 08:56")
	db.SaveActivity(&expected)

	actual, err := db.GetActivity(1)
	if err != nil {
		t.Fatal(err)
	}
	if !tests.Equals(t, actual, expected) {
		t.Fatalf("%v does not equal %v", actual, expected)
	}
}

func TestGetActivtyBetweenStartAndEnd(t *testing.T) {
	start, _ := time.Parse("2006-01-02", "2021-06-28")
	end, _ := time.Parse("2006-01-02", "2021-07-02")

	a1 := createActivityWithEnd("2021-06-28 08:00", "2021-06-28 16:00")

	db := tests.SetupDb(t)
	db.SaveActivity(&a1)

	list, err := db.GetActivitiesBetweenStartAndEnd(start, end, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("expected result len is 1 not %v", len(list))
	}
}

func createActivity(start string) pkg.Activity {
	startTime, _ := time.Parse("2006-01-02 15:04", start)
	return pkg.Activity{
		Start:       &startTime,
		Description: "testing",
		UserID:      1,
	}
}

func createActivityWithEnd(start, end string) pkg.Activity {
	startTime, _ := time.Parse("2006-01-02 15:04", start)
	endTime, _ := time.Parse("2006-01-02 15:04", end)
	return pkg.Activity{
		Start:       &startTime,
		End:         &endTime,
		Description: "testing",
		UserID:      1,
	}
}
