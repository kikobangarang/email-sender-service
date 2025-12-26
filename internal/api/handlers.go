package api

import (
	"net/http"
	"strconv"

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

	emailAddress := r.URL.Path[len("/email/send/"):]
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

func GetJobByIDHandler(w http.ResponseWriter, r *http.Request, svc *email.Service) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idstr := r.URL.Path[len("/email/"):]

	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "invalid job ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	job, err := svc.GetJobByID(id)
	if err != nil {
		http.Error(w, "error fetching job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		"To: " + job.To + "\n" +
			"Subject: " + job.Subject + "\n" +
			"Body: " + job.Body + "\n" +
			"Status: " + string(job.Status) + "\n" +
			"Attempts: " + strconv.Itoa(job.Attempts) + "\n" +
			"Created At: " + job.CreatedAt.String() + "\n",
	))
}
