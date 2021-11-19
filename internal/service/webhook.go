package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
)

func (s *Service) CreateWebhook(webhook pkg.Webhook) (*pkg.Webhook, error) {
	return nil, errors.New("not implemented")
}

func (s *Service) GetWebhooks(employee pkg.Employee) ([]pkg.Webhook, error) {
	return nil, errors.New("not implemented")
}

func (s *Service) startActivityHook(a pkg.Activity) pkg.Activity {

	modified, errs := s.hookHandler(a.UserID, "start_activity", a)
	if errs != nil {
		log.Debug(errs)
		return a
	}

	return modified.(pkg.Activity)
}

func (s *Service) hookHandler(userID uint, eventName string, payload interface{}) (interface{}, []error) {
	hooks, err := s.db.GetWebhooksByUserID(userID)
	if err != nil {
		return nil, []error{err}
	}
	modifyErrors := []error{}
	for _, hook := range hooks {
		if hook.ReadOnly {
			go s.callHook(hook, eventName, payload)
		}
		resp, err := s.callHook(hook, eventName, payload)
		if err == nil {
			payload = resp
		} else {
			modifyErrors = append(modifyErrors, err)
		}
	}
	return payload, modifyErrors
}

func (s *Service) callHook(hook pkg.Webhook, eventName string, payload interface{}) (interface{}, error) {
	body, err := json.Marshal(map[string]interface{}{
		"event":   eventName,
		"payload": payload,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", hook.TargetURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(hook.HeaderKey) > 0 && len(hook.HeaderValue) > 0 {
		req.Header.Set(hook.HeaderKey, hook.HeaderValue)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 304 {
		return payload, nil
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("webhook response %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(payload); err != nil {
		return nil, err
	}
	return payload, nil
}
