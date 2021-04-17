package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	defer resp.Body.Close()

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

func (c *client) CalcCurrentOverview(e Employee) (*Overview, error) {
	resp, err := c.doRequest("GET", "overview/current", nil)
	if err != nil {
		return nil, err
	}
	var o Overview
	err = respToJson(resp, o)
	return &o, err
}

func (c *client) CalcOverviewForThisYear(e Employee) (*Overview, error) {
	resp, err := c.doRequest("GET", "overview", nil)
	if err != nil {
		return nil, err
	}
	var o Overview
	err = respToJson(resp, o)
	return &o, err
}

func (c *client) StartActivity(desc string, employee Employee) (*Activity, error) {
	resp, err := c.doRequest("POST", "activity/"+desc, nil)
	if err != nil {
		return nil, err
	}
	var a Activity
	err = respToJson(resp, a)
	return &a, err
}

func (c *client) AddActivity(activity Activity, employee Employee) (*Activity, error) {
	resp, err := c.doRequest("POST", "activity", activity)
	if err != nil {
		return nil, err
	}
	var a Activity
	err = respToJson(resp, a)
	return &a, err
}

func (c *client) StopRunningActivity(employee Employee) (*Activity, error) {
	resp, err := c.doRequest("DELETE", "activity", nil)
	if err != nil {
		return nil, err
	}
	var a Activity
	err = respToJson(resp, a)
	return &a, err
}

func (c *client) GetActivity(id uint, employee Employee) (*Activity, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("activity/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var a Activity
	err = respToJson(resp, a)
	return &a, err
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
	err = respToJson(resp, h)
	return &h, err
}

func (c *client) GetHollyday(id uint, employee Employee) (*Hollyday, error) {
	resp, err := c.doRequest("POST", fmt.Sprintf("hollyday/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var h Hollyday
	err = respToJson(resp, h)
	return &h, err
}

func (c *client) DelHollyday(id uint, employee Employee) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("hollyday/%d", id), nil)
	return err
}
