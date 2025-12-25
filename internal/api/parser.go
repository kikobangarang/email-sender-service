package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type SendEmailRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func parseSendEmailRequest(r *http.Request) (*SendEmailRequest, error) {
	ct := r.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(ct, "application/json"):
		var req SendEmailRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			return nil, err
		}
		if req.Subject == "" || req.Body == "" {
			return nil, errors.New("subject and body are required")
		}
		return &req, nil

	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		return &SendEmailRequest{
			Subject: r.FormValue("subject"),
			Body:    r.FormValue("body"),
		}, nil

	default:
		return nil, errors.New("unsupported content type")
	}
}
