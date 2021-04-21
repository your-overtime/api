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

	"git.goasum.de/jasper/go-utils/pkg/string_utils"
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
		string_utils.TrimPrefix(resource, "/")
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

func (c *client) CalcOverview(e Employee) (*Overview, error) {
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

func (c *client) AddHollyday(hollyday Hollyday, employee Employee) (*Hollyday, error) {
	resp, err := c.doRequest("POST", "hollyday", hollyday)
	if err != nil {
		return nil, err
	}
	var h Hollyday
	err = respToJson(resp, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (c *client) GetHollyday(id uint, employee Employee) (*Hollyday, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("hollyday/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var h Hollyday
	err = respToJson(resp, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (c *client) GetHollydays(start time.Time, end time.Time, employee Employee) ([]Hollyday, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("hollyday?start=%s&end=%s", timeFormatForQuery(start), timeFormatForQuery(end)), nil)
	if err != nil {
		return nil, err
	}
	var h []Hollyday
	err = respToJson(resp, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (c *client) DelHollyday(id uint, employee Employee) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("hollyday/%d", id), nil)
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
