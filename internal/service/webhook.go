package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
)

func (s *Service) CreateWebhook(webhook pkg.WebhookInput) (*pkg.Webhook, error) {
	_, err := url.ParseRequestURI(webhook.TargetURL)
	if err != nil {
		return nil, err
	}

	hook := pkg.Webhook{
		WebhookInput: webhook,
		UserID:       s.user.ID,
	}

	hookDB, err := s.db.SaveWebhook(data.WebhookDB{Webhook: hook})
	if err != nil {
		return nil, err
	}

	return &hookDB.Webhook, nil
}

func (s *Service) GetWebhooks() ([]pkg.Webhook, error) {
	hookDBs, err := s.db.GetWebhooksByUserID(s.user.ID)
	if err != nil {
		return nil, err
	}

	return castWebhookDBToPkgArray(hookDBs), nil
}

//-- webhook handlers --//

func (s *Service) startActivityHook(a *pkg.Activity) (*pkg.Activity, bool) {

	modified, errs, mayBeModified := s.hookHandler(a.UserID, pkg.StartActivityEvent, a)
	if errs != nil {
		log.Debug(errs)
		return a, false
	}

	return modified.(*pkg.Activity), mayBeModified
}

func (s *Service) endActivityHook(a *pkg.Activity) *pkg.Activity {
	modifed, errs, _ := s.hookHandler(a.UserID, pkg.EndActivityEvent, a)
	if errs != nil {
		return a
	}
	return modifed.(*pkg.Activity)
}

func (s *Service) hookHandler(userID uint, eventName pkg.WebhookEvent, payload interface{}) (interface{}, []error, bool) {
	hooks, err := s.db.GetWebhooksByUserID(userID)
	if err != nil {
		return nil, []error{err}, false
	}
	mayBeModified := false
	modifyErrors := []error{}
	for _, hook := range hooks {
		if hook.ReadOnly {
			go s.callHook(hook.WebhookInput, eventName, payload)
		} else {
			resp, err := s.callHook(hook.WebhookInput, eventName, payload)
			if err == nil {
				payload = resp
				mayBeModified = true
			} else {
				modifyErrors = append(modifyErrors, err)
			}
		}
	}
	if len(modifyErrors) == 0 {
		modifyErrors = nil
	}
	return payload, modifyErrors, mayBeModified
}

func (s *Service) callHook(hook pkg.WebhookInput, eventName pkg.WebhookEvent, payload interface{}) (interface{}, error) {
	body, err := json.Marshal(pkg.WebhookBody{
		Event:   eventName,
		Payload: payload,
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

func castWebhookDBToPkgArray(in []data.WebhookDB) []pkg.Webhook {
	con := make([]pkg.Webhook, len(in))
	for i := range in {
		con[i] = in[i].Webhook
	}

	return con
}
