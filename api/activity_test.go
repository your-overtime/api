package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-overtime/api/pkg"
)

func TestStartActivity(t *testing.T) {
	api, service, _ := setUp(t)

	token, err := service.CreateToken(pkg.InputToken{
		Name: "test",
	})
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	w := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/api/v1/activity/start?desc=Test&token="+token.Token, nil)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatal("expect status 201 but got: ", w.Result().StatusCode)
	}

	defer w.Result().Body.Close()
	a := pkg.Activity{}

	err = json.NewDecoder(w.Result().Body).Decode(&a)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	if a.Description != "Test" {
		t.Error("expect Test got ", a.Description)
	}

	w = httptest.NewRecorder()
	payload, _ := json.Marshal(map[string]string{
		"desc": "Test2",
	})
	req, err = http.NewRequest("POST", "/api/v1/activity/start?token="+token.Token, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatal("expect status 201 but got: ", w.Result().StatusCode)
	}

	defer w.Result().Body.Close()
	a = pkg.Activity{}

	err = json.NewDecoder(w.Result().Body).Decode(&a)
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	if a.Description != "Test2" {
		t.Error("expect Test got ", a.Description)
	}
}
