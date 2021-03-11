package sciensano

import (
	"net/http"
)

// APIClient queries different Sciensano APIs
type APIClient struct {
	client http.Client
}

const baseURL = "https://epistat.sciensano.be/Data/"
