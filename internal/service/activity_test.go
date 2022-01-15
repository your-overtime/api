package service_test

import (
	"testing"
	"time"

	"github.com/your-overtime/api/pkg"
)

func TestActivityDuration(t *testing.T) {
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

	a1, err := service.GetActivity(a.ID)
	if err != nil {
		t.Fatal("expected no error but got", err)
	}
	if a1.ActualDurationInMinutes != 5 {
		t.Error("expected ActualDurationInMinutes to be 5 but got", a1.ActualDurationInMinutes)
	}

	if a1.EventualDurationInMinutes != 5 {
		t.Error("expected EventualDurationInMinutes to be 5 but got", a1.EventualDurationInMinutes)
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
