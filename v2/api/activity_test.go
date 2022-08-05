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

func TestStartActivity(t *testing.T) {
	api, service, _ := setUp(t)

	token, err := service.CreateToken(pkg.InputToken{
		Name: "test",
	})
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	w := httptest.NewRecorder()
	payload, _ := json.Marshal(map[string]string{
		"Description": "Test2",
	})
	req, err := http.NewRequest("POST", "/api/v1/activity?token="+token.Token, bytes.NewBuffer(payload))
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
		fmt.Println(err)
		t.Fatal("expect no error but got: ", err)
	}

	if a.Description != "Test2" {
		t.Error("expect Test got ", a.Description)
	}

	w = httptest.NewRecorder()
	payload, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest("POST", "/api/v1/activity?token="+token.Token, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal("expect no error but got: ", err)
	}

	api.Router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatal("expect status 400 but got: ", w.Result().StatusCode)
	}
}
