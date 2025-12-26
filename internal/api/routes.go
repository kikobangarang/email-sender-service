package api

import (
	"net/http"

	"github.com/kikobangarang/email-sender-service/internal/email"
)

func RegisterHandlers(mux *http.ServeMux, svc *email.Service) {
	mux.HandleFunc("/health", healthHandler)
	// expects /send/{email}
	mux.HandleFunc("/email/send/", func(w http.ResponseWriter, r *http.Request) {
		sendEmailHandler(w, r, svc)
	})
	mux.HandleFunc("/email/", func(w http.ResponseWriter, r *http.Request) {
		GetJobByIDHandler(w, r, svc)
	})
}
