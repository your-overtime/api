package service_test

import (
	"testing"
	"time"

	"github.com/your-overtime/api/pkg"
)

func TestStartStopActivityDuration(t *testing.T) {
	service, _ := setUp(t)
	start := time.Date(2022, 1, 1, 8, 0, 0, 0, time.UTC)
	end := time.Date(2022, 1, 1, 8, 5, 0, 0, time.UTC)
	_, err := service.StartActivityWithTime("demo", start)
	if err != nil {
		t.Fatal("expected no error but got", err)
	}

	a, err := service.StopRunningActivityWithTime(end)
	if err != nil {
		t.Fatal("expected no error but got", err)
	}
	if a.ActualDurationInMinutes != 5 {
		t.Error("expected ActualDurationInMinutes to be 5 but got", a.ActualDurationInMinutes)
	}

	if a.EventualDurationInMinutes != 5 {
		t.Error("expected EventualDurationInMinutes to be 5 but got", a.EventualDurationInMinutes)
	}
}

func TestStopRunningActivityWhenNewStarting(t *testing.T) {
	service, _ := setUp(t)
	_, err := service.StartActivity("Valid")
	if err != nil {
		t.Fatal("expected error to be nil but got", err)
	}
	o, err := service.CalcOverview(time.Now())
	if err != nil {
		t.Fatal("expected error to be nil but got", err)
	}
	if o.ActiveActivity == nil {
		t.Error("expected active activity not to be nil")
	}
	_, err = service.StartActivity("")
	if err == nil || err != pkg.ErrEmptyDescriptionNotAllowed {
		t.Errorf("expected error not be %s not got %v", pkg.ErrEmptyDescriptionNotAllowed, err)
	}

	o, err = service.CalcOverview(time.Now())
	if err != nil {
		t.Fatal("expected error to be nil but got", err)
	}
	if o.ActiveActivity == nil {
		t.Error("expected active activity not to be nil")
	}
}

func TestCreateActivity(t *testing.T) {
	service, _ := setUp(t)
	start := time.Date(2022, 1, 1, 8, 0, 0, 0, time.UTC)
	end := time.Date(2022, 1, 1, 8, 5, 0, 0, time.UTC)
	in := pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start:       &start,
			End:         &end,
			Description: "test create activity",
		},
	}
	out, err := service.AddActivity(in)
	if err != nil {
		t.Error("expected no error but got", err)
	}
	if out.ID <= 0 {
		t.Error("expected id to be greater than 0 but got", out.ID)
	}
	if out.ActualDurationInMinutes != 5 {
		t.Error("expected ActualDurationInMinutes to be 5 but got", out.ActualDurationInMinutes)
	}

	if out.EventualDurationInMinutes != 5 {
		t.Error("expected EventualDurationInMinutes to be 5 but got", out.EventualDurationInMinutes)
	}

	in = pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start:       &end,
			End:         &start,
			Description: "test create activity",
		},
	}
	_, err = service.AddActivity(in)
	if err == nil || err != pkg.ErrStartMustBeBeforeEnd {
		t.Error("expected error start must before end error but got nil")
	}
}

func TestUpdateActivity(t *testing.T) {
	service, _ := setUp(t)
	start := time.Date(2022, 1, 1, 8, 0, 0, 0, time.UTC)
	end := time.Date(2022, 1, 1, 8, 5, 0, 0, time.UTC)
	in := pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start:       &start,
			End:         &end,
			Description: "test create activity",
		},
	}
	ac, err := service.AddActivity(in)
	if err != nil {
		t.Fatal(err)
	}
	ac.Description = ""
	_, err = service.UpdateActivity(*ac)
	if err == nil || err != pkg.ErrEmptyDescriptionNotAllowed {
		t.Error("expected no empty description error but got ", err)
	}

	ac.Description = "new description"
	out, err := service.UpdateActivity(*ac)
	if err != nil {
		t.Fatal(err)
	}
	if ac.Description != out.Description {
		t.Errorf("expected new description to be %s but got %s", ac.Description, out.Description)
	}

	ac.Start = &end
	ac.End = &start
	_, err = service.UpdateActivity(*ac)
	if err == nil || err != pkg.ErrStartMustBeBeforeEnd {
		t.Error("expected error start must before end error but got", err)
	}

	ac.Start = &start
	newEnd := time.Date(2022, 1, 1, 8, 10, 0, 0, time.UTC)
	ac.End = &newEnd
	out, err = service.UpdateActivity(*ac)
	if err != nil {
		t.Fatal(err)
	}
	if out.ActualDurationInMinutes != 10 {
		t.Error("expected ActualDurationInMinutes to be 10 but got", out.ActualDurationInMinutes)
	}
	if out.EventualDurationInMinutes != 10 {
		t.Error("expected EventualDurationInMinutes to be 10 but got", out.EventualDurationInMinutes)
	}
}

func TestDeleteActivty(t *testing.T) {
	service, _ := setUp(t)
	start := time.Date(2022, 1, 1, 8, 0, 0, 0, time.UTC)
	end := time.Date(2022, 1, 1, 8, 5, 0, 0, time.UTC)
	in := pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start:       &start,
			End:         &end,
			Description: "test create activity",
		},
	}
	ac, err := service.AddActivity(in)
	if err != nil {
		t.Fatal(err)
	}

	err = service.DelActivity(ac.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.GetActivity(ac.ID)
	if err == nil {
		t.Error("expectet error but got nil")
	}
}

func TestCalculateDuration(t *testing.T) {
	service, _ := setUp(t)
	start := time.Date(2022, 1, 1, 8, 0, 0, 0, time.UTC)
	end := time.Date(2022, 1, 1, 8, 5, 0, 0, time.UTC)
	a := pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start: &start,
			End:   &end,
		},
	}
	service.CalculateDuration(&a)
	if a.ActualDurationInMinutes != 5 {
		t.Error("expected ActualDurationInMinutes to be 5 but got", a.ActualDurationInMinutes)
	}

	if a.EventualDurationInMinutes != 5 {
		t.Error("expected EventualDurationInMinutes to be 5 but got", a.EventualDurationInMinutes)
	}

}
