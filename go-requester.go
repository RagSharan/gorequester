package gorequester

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var defaultClient = resty.New()

func init() {
	defaultClient.SetTimeout(10 * time.Second)
	defaultClient.SetHeader("User-Agent", "gorequester/1.0")
}

// NewClient allows creating a custom resty client
func NewClient(timeout time.Duration, defaultHeaders map[string]string) *resty.Client {
	client := resty.New()
	client.SetTimeout(timeout)
	for k, v := range defaultHeaders {
		client.SetHeader(k, v)
	}
	return client
}

func Get(url string, headers map[string]string) (*resty.Response, error) {
	request := defaultClient.R()
	for k, v := range headers {
		request.SetHeader(k, v)
	}
	return request.Get(url)
}

func Post(url string, body []byte, headers map[string]string) (*resty.Response, error) {
	request := defaultClient.R().SetBody(body)
	for k, v := range headers {
		request.SetHeader(k, v)
	}
	return request.Post(url)
}

func Put(url string, body []byte, headers map[string]string) (*resty.Response, error) {
	request := defaultClient.R().SetBody(body)
	for k, v := range headers {
		request.SetHeader(k, v)
	}
	return request.Put(url)
}

func Patch(url string, body []byte, headers map[string]string) (*resty.Response, error) {
	request := defaultClient.R().SetBody(body)
	for k, v := range headers {
		request.SetHeader(k, v)
	}
	return request.Patch(url)
}

func Delete(url string, headers map[string]string) (*resty.Response, error) {
	request := defaultClient.R()
	for k, v := range headers {
		request.SetHeader(k, v)
	}
	return request.Delete(url)
}

func SendJSON(url string, method string, jsonBody interface{}, headers map[string]string) (*resty.Response, error) {
	request := defaultClient.R().SetHeader("Content-Type", "application/json").SetBody(jsonBody)
	for k, v := range headers {
		request.SetHeader(k, v)
	}

	switch method {
	case http.MethodPost:
		return request.Post(url)
	case http.MethodPut:
		return request.Put(url)
	case http.MethodPatch:
		return request.Patch(url)
	case http.MethodDelete:
		return request.Delete(url)
	case http.MethodGet:
		return request.Get(url)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}
