package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	token, err := service.CreateToken(pkg.InputToken{
		Name: "test",
	})
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/api/v2/holiday/1?token="+token.Token, nil)
	req.Header.Add("accept", "application/json")
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Error("expect status 404 but got: ", w.Result().StatusCode)
	}

	iph := pkg.InputHoliday{
		Start:       tests.ParseDayTime("2021-01-08 00:00"),
		End:         tests.ParseDayTime("2021-01-14 23:59"),
		Description: "Test",
		Type:        pkg.HolidayTypeLegalUnpaidFree,
	}
	h, err := service.AddHoliday(pkg.Holiday{
		UserID:       user.ID,
		InputHoliday: iph,
	})

	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v2/holiday/%d?token=%s", h.ID, token.Token), nil)
	req.Header.Add("accept", "application/json")
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}
	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Error("expect status 200 but got: ", w.Result().StatusCode)
	}

	respHoliday := pkg.Holiday{}
	defer w.Result().Body.Close()
	err = json.NewDecoder(w.Result().Body).Decode(&respHoliday)

	if err != nil {
		t.Error("expect nil but got ", err)
	}

	if h.End.String() != iph.End.String() {
		t.Errorf("expect %s but got %s", h.End.String(), iph.End.String())
	}

	if h.Start.String() != iph.Start.String() {
		t.Errorf("expect %s but got %s", h.Start.String(), iph.Start.String())
	}

	if respHoliday.End.String() != iph.End.String() {
		t.Errorf("expect %s but got %s", respHoliday.End.String(), iph.End.String())
	}

	if respHoliday.Start.String() != iph.Start.String() {
		t.Errorf("expect %s but got %s", respHoliday.Start.String(), iph.Start.String())
	}

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v2/holiday?token=%s&start=%s&end=%s", token.Token, h.Start.Format(time.RFC3339Nano), h.End.Format(time.RFC3339Nano)), nil)
	req.Header.Add("accept", "application/json")
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}
	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Error("expect status 200 but got: ", w.Result().StatusCode)
	}

	respHolidays := []pkg.Holiday{}
	defer w.Result().Body.Close()
	err = json.NewDecoder(w.Result().Body).Decode(&respHolidays)
	if err != nil {
		t.Error("expect nil but got ", err)
	}

	if len(respHolidays) != 1 {
		t.Error("expect 1 but got: ", len(respHolidays))
	}
}
