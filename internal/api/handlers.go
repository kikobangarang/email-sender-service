package api

import (
	"net/http"

	"github.com/kikobangarang/email-sender-service/internal/email"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request, svc *email.Service) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	emailAddress := r.URL.Path[len("/send/"):]
	if emailAddress == "" {
		http.Error(w, "missing email", http.StatusBadRequest)
		return
	}

	emailReq, err := parseSendEmailRequest(r)
	if err != nil {
		http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := svc.SendEmail(emailAddress, emailReq.Subject, emailReq.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(
		"email queued for " + emailAddress + "\n",
	))
}
