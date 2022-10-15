package scientia

import (
	"net/http"
	"time"
)

var token string

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}
