package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-overtime/api/v2/pkg"
)

func TestUpdateAccount(t *testing.T) {
	api, service, user := setUp(t)

	w := httptest.NewRecorder()

	fields := map[string]interface{}{
		"WeekWorkingTimeInMinutes": uint(1800),
	}
	dataBytes := new(bytes.Buffer)
	err := json.NewEncoder(dataBytes).Encode(fields)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	token, err := service.CreateToken(pkg.InputToken{
		Name: "test",
	})
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/api/v1/account", bytes.NewBuffer(dataBytes.Bytes()))
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.Token))

	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Error("expect status 200 but got: ", w.Result().StatusCode)
	}

	req, err = http.NewRequest(http.MethodGet, "/api/v1/account", nil)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.Token))

	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	respUser := pkg.User{}

	err = json.NewDecoder(w.Result().Body).Decode(&respUser)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	if respUser.WeekWorkingTimeInMinutes != 1800 {
		t.Error("WeekWorkingTimeInMinutes not changed")
	}

	if respUser.ID != user.ID {
		t.Error("ids not match")
	}
}
