package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type client struct {
	host       string
	authHeader string
}

func respToJson(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *client) doRequest(method string, resource string, data interface{}) (*http.Response, error) {
	dataBytes := new(bytes.Buffer)
	err := json.NewEncoder(dataBytes).Encode(data)

	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(resource, "/") {
		strings.TrimPrefix(resource, "/")
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%sapi/v1/%s", c.host, resource), bytes.NewBuffer(dataBytes.Bytes()))

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", c.authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode >= 400 {
		return resp, errors.New(resp.Status)
	}

	return resp, err
}

func InitOvertimeClient(host string, authHeader string) OvertimeService {
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "http://" + host
	}
	return &client{
		host:       host,
		authHeader: authHeader,
	}
}

func (c *client) CalcOverview(e Employee, d time.Time) (*Overview, error) {
	resp, err := c.doRequest("GET", "overview", nil)
	if err != nil {
		return nil, err
	}
	var o Overview
	err = respToJson(resp, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (c *client) StartActivity(desc string, employee Employee) (*Activity, error) {
	resp, err := c.doRequest("POST", "activity/"+desc, nil)
	if err != nil {
		return nil, err
	}
	var a Activity
	err = respToJson(resp, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (c *client) AddActivity(activity Activity, employee Employee) (*Activity, error) {
	resp, err := c.doRequest("POST", "activity", activity)
	if err != nil {
		return nil, err
	}
	var a Activity
	err = respToJson(resp, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (c *client) UpdateActivity(activity Activity, employee Employee) (*Activity, error) {
	resp, err := c.doRequest("PUT", fmt.Sprintf("activity/%d", activity.ID), activity)
	if err != nil {
		return nil, err
	}
	var a Activity
	err = respToJson(resp, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (c *client) StopRunningActivity(employee Employee) (*Activity, error) {
	resp, err := c.doRequest("DELETE", "activity", nil)
	if err != nil {
		return nil, err
	}
	var a Activity
	err = respToJson(resp, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (c *client) GetActivity(id uint, employee Employee) (*Activity, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("activity/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var a Activity
	err = respToJson(resp, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func timeFormatForQuery(t time.Time) string {
	return url.QueryEscape(t.Format(time.RFC3339))
}

func (c *client) GetActivities(start time.Time, end time.Time, employee Employee) ([]Activity, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("activity?start=%s&end=%s", timeFormatForQuery(start), timeFormatForQuery(end)), nil)
	if err != nil {
		return nil, err
	}
	var a []Activity
	err = respToJson(resp, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (c *client) DelActivity(id uint, employee Employee) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("activity/%d", id), nil)
	return err
}

func (c *client) AddHoliday(holiday Holiday, employee Employee) (*Holiday, error) {
	resp, err := c.doRequest("POST", "holiday", holiday)
	if err != nil {
		return nil, err
	}
	var h Holiday
	err = respToJson(resp, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (c *client) UpdateHoliday(holiday Holiday, employee Employee) (*Holiday, error) {
	resp, err := c.doRequest("PUT", fmt.Sprintf("holiday/%d", holiday.ID), holiday)
	if err != nil {
		return nil, err
	}
	var h Holiday
	err = respToJson(resp, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (c *client) GetHoliday(id uint, employee Employee) (*Holiday, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("holiday/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var h Holiday
	err = respToJson(resp, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (c *client) GetHolidays(start time.Time, end time.Time, employee Employee) ([]Holiday, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("holiday?start=%s&end=%s", timeFormatForQuery(start), timeFormatForQuery(end)), nil)
	if err != nil {
		return nil, err
	}
	var h []Holiday
	err = respToJson(resp, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (c *client) GetHolidaysByType(start time.Time, end time.Time, hType HolidayType, employee Employee) ([]Holiday, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("holiday?start=%s&end=%s&type=%s", timeFormatForQuery(start), timeFormatForQuery(end), hType), nil)
	if err != nil {
		return nil, err
	}
	var h []Holiday
	err = respToJson(resp, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (c *client) AddWorkDay(workday WorkDay, employee Employee) (*WorkDay, error) {
	resp, err := c.doRequest("POST", "workday", workday)
	if err != nil {
		return nil, err
	}
	var w WorkDay
	err = respToJson(resp, &w)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (c *client) GetWorkDays(start time.Time, end time.Time, employee Employee) ([]WorkDay, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("workday?start=%s&end=%s", timeFormatForQuery(start), timeFormatForQuery(end)), nil)
	if err != nil {
		return nil, err
	}
	var ws []WorkDay
	err = respToJson(resp, &ws)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (c *client) DelHoliday(id uint, employee Employee) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("holiday/%d", id), nil)
	return err
}

func (c *client) SaveEmployee(employee Employee, adminToken string) (*Employee, error) {
	resp, err := c.doRequest("POST", "employee?adminToken="+adminToken, employee)
	if err != nil {
		return nil, err
	}
	var e Employee
	err = respToJson(resp, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (c *client) UpdateAccount(fields map[string]interface{}, employee Employee) (*Employee, error) {
	resp, err := c.doRequest("PATCH", "account", fields)
	if err != nil {
		return nil, err
	}
	var e Employee
	err = respToJson(resp, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (c *client) GetAccount() (*Employee, error) {
	resp, err := c.doRequest("GET", "account", nil)
	if err != nil {
		return nil, err
	}
	var e Employee
	err = respToJson(resp, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (c *client) DeleteEmployee(login string, adminToken string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("employee/%s?adminToken=%s", login, adminToken), nil)
	return err
}

func (c *client) CreateToken(token InputToken, employee Employee) (*Token, error) {
	resp, err := c.doRequest("POST", "token", token)
	if err != nil {
		return nil, err
	}
	var t Token
	err = respToJson(resp, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *client) DeleteToken(tokenID uint, employee Employee) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("token/%d", tokenID), nil)
	return err
}

func (c *client) GetTokens(employee Employee) ([]Token, error) {
	resp, err := c.doRequest("GET", "token", nil)
	if err != nil {
		return nil, err
	}
	var ts []Token
	err = respToJson(resp, &ts)
	if err != nil {
		return nil, err
	}
	return ts, nil
}
