package scientia

import (
	"net/http"
	"time"
)

var token string

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

const (
	baseURL = "https://api-materials.doc.ic.ac.uk/"
)
