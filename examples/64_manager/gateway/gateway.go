package gateway

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gitlab.com/mandalore/go-app/app"
)

// Gateway ...
type Gateway struct {
	host           string
	client         *http.Client
	defaultHeaders map[string]string
}

// NewGateway ...
func NewGateway(host string) *Gateway {
	gateway := &Gateway{
		host:           host,
		client:         &http.Client{},
		defaultHeaders: make(map[string]string),
	}

	return gateway
}

// AddDefaultHeader ...
func (gateway *Gateway) AddDefaultHeader(key, value string) {
	gateway.defaultHeaders[key] = value
}

// Request ...
func (gateway *Gateway) Request(method string, endpoint string, headers map[string]string, body io.Reader) (int, []byte, error) {
	url := gateway.host + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, nil, app.NewApplicationError(app.ErrorUnexpected, fmt.Sprintf("error creating request [url:%s]", err.Error()), err)
	}

	for key, value := range gateway.defaultHeaders {
		req.Header.Add(key, value)
	}
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	response, err := gateway.client.Do(req)
	if err != nil {
		return 0, nil, app.NewApplicationError(app.ErrorUnexpected, fmt.Sprintf("error running request [url:%s]", err.Error()), err)
	}
	defer response.Body.Close()

	output, err := ioutil.ReadAll(response.Body)

	return response.StatusCode, output, err
}
