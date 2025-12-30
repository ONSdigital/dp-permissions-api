package sdk

import (
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/v3/request"
)

type Headers struct {
	ServiceAuthToken string
}

// Add adds any provided headers to the request
func (h *Headers) Add(req *http.Request) {
	if h.ServiceAuthToken != "" {
		dprequest.AddServiceTokenHeader(req, h.ServiceAuthToken)
	}
}
