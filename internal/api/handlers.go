package api

import (
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Path[len("/send/"):]
	if email == "" {
		http.Error(w, "missing email", http.StatusBadRequest)
		return
	}

	emailReq, err := parseSendEmailRequest(r)
	if err != nil {
		http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("email queued for " + email + ":\n" +
		"From: " + emailReq.From + "\n" +
		"Subject: " + emailReq.Subject + "\n" +
		"Body: " + emailReq.Body + "\n"))
}
