package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-overtime/api/v2/api"
	"github.com/your-overtime/api/v2/internal/data"
	"github.com/your-overtime/api/v2/internal/service"
	"github.com/your-overtime/api/v2/pkg"
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

func setUp(t *testing.T) (*api.API, *service.Service, *pkg.User) {
	db := tests.SetupDb(t)
	err := db.SaveUser(&data.UserDB{User: e})
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}
	s := service.Init(db)
	eDB, err := db.GetUserByLogin(e.Login)
	if err != nil {
		t.Fatal("expect no error but got ", err)
	}
	e = eDB.User
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

	req, err := http.NewRequest("GET", "/api/v2/holiday/1", nil)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Error("expect status 401 but got: ", w.Result().StatusCode)
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

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v2/holiday/%d?token=%s", int(h.ID), token.Token), nil)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}
	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Error("expect status 200 but got: ", w.Result().StatusCode, h)
	}

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v2/holiday?token=%s", token.Token), nil)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}
	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Error("expect status 400 but got: ", w.Result().StatusCode)
	}
}
