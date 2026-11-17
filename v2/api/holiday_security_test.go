package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-overtime/api/v2/api"
	"github.com/your-overtime/api/v2/internal/data"
	"github.com/your-overtime/api/v2/internal/service"
	"github.com/your-overtime/api/v2/pkg"
	"github.com/your-overtime/api/v2/tests"
)

func TestHolidayInvalidTypeAndForeignAccess(t *testing.T) {
	db := tests.SetupDb(t)
	mainSvc := service.Init(db)
	a := api.Init(mainSvc, "1234")
	a.CreateEndpoints()

	newUser := func(login string) (pkg.OvertimeService, *pkg.User) {
		u := pkg.User{Name: login, Surname: "Test", Login: login, Password: "secret", WeekWorkingTimeInMinutes: 1920, NumWorkingDays: 5}
		if err := db.SaveUser(&data.UserDB{User: u}); err != nil {
			t.Fatal(err)
		}
		uDB, err := db.GetUserByLogin(login)
		if err != nil {
			t.Fatal(err)
		}
		u = uDB.User
		return mainSvc.GetOrCreateInstanceForUser(&u), &u
	}
	ownerSvc, owner := newUser("owner")
	otherSvc, _ := newUser("other")

	h, err := ownerSvc.AddHoliday(pkg.Holiday{
		UserID: owner.ID,
		InputHoliday: pkg.InputHoliday{
			Start: tests.ParseDayTime("2021-01-08 00:00"),
			End:   tests.ParseDayTime("2021-01-10 00:00"),
			Type:  pkg.HolidayTypeFree,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	ownerToken, err := ownerSvc.CreateToken(pkg.InputToken{Name: "owner"})
	if err != nil {
		t.Fatal(err)
	}

	otherToken, err := otherSvc.CreateToken(pkg.InputToken{Name: "other"})
	if err != nil {
		t.Fatal(err)
	}

	do := func(method, path, token, body string) int {
		req, _ := http.NewRequest(method, fmt.Sprintf("%s?token=%s", path, token), strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		a.Router.ServeHTTP(w, req)
		return w.Result().StatusCode
	}

	valid := `{"Start":"2021-02-01T00:00:00Z","End":"2021-02-02T00:00:00Z","Type":"sick"}`
	invalid := `{"Start":"2021-02-01T00:00:00Z","End":"2021-02-02T00:00:00Z","Type":"bogus"}`
	path := fmt.Sprintf("/api/v2/holiday/%d", h.ID)

	if code := do("POST", "/api/v2/holiday", ownerToken.Token, invalid); code != http.StatusBadRequest {
		t.Errorf("POST invalid type: expected 400, got %d", code)
	}
	if code := do("PUT", path, ownerToken.Token, invalid); code != http.StatusBadRequest {
		t.Errorf("PUT invalid type: expected 400, got %d", code)
	}
	if code := do("PUT", path, otherToken.Token, valid); code != http.StatusForbidden {
		t.Errorf("PUT foreign holiday: expected 403, got %d", code)
	}
	if code := do("PUT", "/api/v2/holiday/9999", ownerToken.Token, valid); code != http.StatusNotFound {
		t.Errorf("PUT unknown holiday: expected 404, got %d", code)
	}
	if code := do("DELETE", path, otherToken.Token, ""); code != http.StatusForbidden {
		t.Errorf("DELETE foreign holiday: expected 403, got %d", code)
	}

	got, err := ownerSvc.GetHoliday(h.ID)
	if err != nil {
		t.Fatalf("holiday of owner must still exist and belong to owner: %v", err)
	}
	if got.Type != pkg.HolidayTypeFree {
		t.Errorf("holiday must be unchanged, got type %q", got.Type)
	}

	if code := do("PUT", path, ownerToken.Token, valid); code != http.StatusOK {
		t.Errorf("PUT own holiday: expected 200, got %d", code)
	}
	if code := do("DELETE", path, ownerToken.Token, ""); code != http.StatusOK {
		t.Errorf("DELETE own holiday: expected 200, got %d", code)
	}
}
