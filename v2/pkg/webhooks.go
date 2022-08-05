package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HandleActivity = func(a *Activity) error

type WebhookHandler interface {

	// HandleStartActivity will be called on start events
	HandleStartActivity(handler HandleActivity)

	// HandleEndActivity will be called on end events
	HandleEndActivity(Handler HandleActivity)

	// ListenAndServe listens on the TCP network address addr and then calls Serve with handler
	// to handle requests on incoming connections.
	ListenAndServe(host string) error
}

type webhookHandlerImpl struct {
	start HandleActivity
	end   HandleActivity
}

// NewWebhookHandler returns simple http server, on which callbacks for events can be registered
// and they be automaticly called on their event type.
func NewWebhookHanlder() WebhookHandler {
	return &webhookHandlerImpl{}
}

func (wh *webhookHandlerImpl) HandleStartActivity(handler HandleActivity) {
	wh.start = handler
}

func (wh *webhookHandlerImpl) HandleEndActivity(handler HandleActivity) {
	wh.end = handler
}

func (wh *webhookHandlerImpl) ListenAndServe(addr string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.Method != "POST" {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}

		var body WebhookBody
		var resp []byte
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch body.Event {
		case StartActivityEvent:
			if wh.start != nil {
				resp = wh.handleActivity(body.Payload, wh.start)
			}
		case EndActivityEvent:
			if wh.end != nil {
				resp = wh.handleActivity(body.Payload, wh.end)
			}
		}

		if len(resp) > 0 {
			fmt.Fprintf(w, "%s", resp)
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.WriteHeader(http.StatusNotModified)
		}
	})
	return http.ListenAndServe(addr, nil)
}

// helper to call ActivityHand;er on activty payload
func (wh *webhookHandlerImpl) handleActivity(v interface{}, handler HandleActivity) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var activity Activity
	if err := json.Unmarshal(bytes, &activity); err != nil {
		return nil
	}
	handler(&activity)
	resp, err := json.Marshal(activity)
	if err != nil {
		resp = nil
	}
	return resp
}
