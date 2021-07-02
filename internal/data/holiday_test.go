package data_test

import (
	"testing"
	"time"

	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/tests"
)

func createHoliday(start string) pkg.Holiday {
	startTime, _ := time.Parse("2006-01-02", start)
	return pkg.Holiday{
		UserID:       1,
		Start:        startTime,
		End:          startTime,
		Description:  "Test",
		LegalHoliday: false,
	}
}

func createHolidayWithEnd(start, end string) pkg.Holiday {
	startTime, _ := time.Parse("2006-01-02", start)
	endTime, _ := time.Parse("2006-01-02", end)
	return pkg.Holiday{
		UserID:       1,
		Start:        startTime,
		End:          endTime,
		Description:  "Test",
		LegalHoliday: false,
	}
}
func TestSaveHoliday(t *testing.T) {
	db := tests.SetupDb(t)
	insert := createHoliday("2021-07-01")
	if err := db.SaveHoliday(&insert); err != nil {
		t.Fatal(err)
	}

	if insert.ID == 0 {
		t.Fatal("exptected auto incremented id after insert")
	}

	if insert.ID != 1 {
		t.Fatal("expected first inserted holiday has ID 1")
	}
}

func TestGetHoliday(t *testing.T) {
	db := tests.SetupDb(t)
	expected := createHoliday("2021-07-02")
	db.SaveHoliday(&expected)

	actual, err := db.GetHoliday(1)
	if err != nil {
		t.Fatal(err)
	}

	if !tests.Equals(t, actual, expected) {
		t.Fatalf("%v does not equal %v", actual, expected)
	}

	expected2 := createHoliday("2021-07-03")
	db.SaveHoliday(&expected2)

	actual2, err := db.GetHoliday(2)
	if err != nil {
		t.Fatal(err)
	}
	if !tests.Equals(t, actual2, expected2) {
		t.Fatalf("%v does not equal %v", actual2, expected2)
	}
	if tests.Equals(t, actual2, expected) {
		t.Fatalf("%v shoud not equal %v", actual2, expected)
	}

}

func TestGetHolidayBetweenStartAndEnd(t *testing.T) {
	start, _ := time.Parse("2006-01-02", "2021-06-01")
	end, _ := time.Parse("2006-01-02", "2021-06-07")

	// h1 ok			----sxxx|xxxxxxxxxxxxxxx|xxxxe-----
	// h2 ok			--------|----sxxxxxxxxxx|xxxe------
	// h3 ok			-----sxx|xxxxxxxxxe-----|----------
	// h4 ok			--------|-----sxxxxe----|----------
	// h5 not ok		--sxe---|---------------|----------
	// h6 not ok		--------|---------------|--sxxxxxe-
	//
	// query range: 	--------sxxxxxxxxxxxxxxxe----------

	h1 := createHolidayWithEnd("2021-05-20", "2021-06-09")
	h2 := createHolidayWithEnd("2021-06-04", "2021-06-08")
	h3 := createHolidayWithEnd("2021-05-30", "2021-06-04")
	h4 := createHolidayWithEnd("2021-06-02", "2021-06-06")
	h5 := createHolidayWithEnd("2021-05-20", "2021-05-30")
	h6 := createHolidayWithEnd("2021-06-09", "2021-06-30")

	// insert
	db := tests.SetupDb(t)
	db.SaveHoliday(&h1)
	db.SaveHoliday(&h2)
	db.SaveHoliday(&h3)
	db.SaveHoliday(&h4)
	db.SaveHoliday(&h5)
	db.SaveHoliday(&h6)

	list, err := db.GetHolidaysBetweenStartAndEnd(start, end, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 4 {
		t.Fatalf("expected result len is 4 not %v", len(list))
	}

	if !(tests.Equals(t, list[0], h1) && tests.Equals(t, list[1], h2) && tests.Equals(t, list[2], h3) && tests.Equals(t, list[3], h4)) {
		t.Fatalf("nope")
	}

	// * start is on start
	// * start is on end
	// * end is on start
	h7 := createHoliday("2021-06-01")
	h8 := createHoliday("2021-06-07")
	h9 := createHolidayWithEnd("2021-05-01", "2021-06-01")
	db.SaveHoliday(&h7)
	db.SaveHoliday(&h8)
	db.SaveHoliday(&h9)

	list, err = db.GetHolidaysBetweenStartAndEnd(start, end, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 7 {
		t.Fatalf("expect result len is 7 not %v", len(list))
	}
	if !tests.Equals(t, list[4], h7) && tests.Equals(t, list[5], h8) && tests.Equals(t, list[6], h9) {
		t.Fatal("nope")
	}

}
