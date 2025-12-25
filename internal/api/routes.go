package api

import "net/http"

func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/send/", sendEmailHandler) // expects /send/{email}
}
