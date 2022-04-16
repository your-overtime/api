package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-overtime/api/api"
	"github.com/your-overtime/api/internal/data"
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

func setUp(t *testing.T) (*api.API, *service.Service, *pkg.User) {
	db := tests.SetupDb(t)
	err := db.SaveUser(&data.UserDB{User: e})
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}
	s := service.Init(db)
	ots := s.GetOrCreateInstanceForUser(&e)

	actualService, ok := ots.(*service.Service)
	if !ok {
		t.Fatal("wrong service implementation")
	}

	api := api.Init(s, "1234")
	api.CreateEndpoints()

	return api, actualService, &e
}

func TestGetHoliday(t *testing.T) {
	api, service, user := setUp(t)

	w := httptest.NewRecorder()

	req, err := http.NewRequest("Get", "/api/v1/holiday/1", nil)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Error("expect status 404 but got: ", w.Result().StatusCode)
	}

	h, err := service.AddHoliday(pkg.Holiday{
		UserID: user.ID,
		InputHoliday: pkg.InputHoliday{
			Start:       tests.ParseDayTime("2021-01-08 00:00"),
			End:         tests.ParseDayTime("2021-01-14 00:00"),
			Description: "Test",
			Type:        pkg.HolidayTypeLegalUnpaidFree,
		},
	})

	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	token, err := service.CreateToken(pkg.InputToken{
		Name: "test",
	})
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	req, err = http.NewRequest("Get", fmt.Sprintf("/api/v1/holiday/%d?token=%s", h.ID, token.Token), nil)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}
	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Error("expect status 404 but got: ", w.Result().StatusCode)
	}

	req, err = http.NewRequest("Get", fmt.Sprintf("/api/v1/holiday?token=%s", token.Token), nil)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}
	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Error("expect status 404 but got: ", w.Result().StatusCode)
	}
}
