package gorequester

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type RequestInput struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
	Retries int // optional retry count
}

// Send performs the HTTP request with optional retries
func Send(input RequestInput) ([]byte, error) {
	if input.Method == "" {
		input.Method = "GET"
	}
	if input.Timeout == 0 {
		input.Timeout = 10 * time.Second
	}
	if input.Retries == 0 {
		input.Retries = 1
	}

	var lastErr error
	for i := 0; i < input.Retries; i++ {
		client := &http.Client{Timeout: input.Timeout}
		req, err := http.NewRequest(input.Method, input.URL, bytes.NewReader(input.Body))
		if err != nil {
			return nil, err
		}

		for k, v := range input.Headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 500 * time.Millisecond) // exponential-ish backoff
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, errors.New("request failed with status: " + resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}
	return nil, lastErr
}

// SendJSON simplifies sending a JSON POST/PUT
func SendJSON(url, method string, payload interface{}, headers map[string]string) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"

	return Send(RequestInput{
		URL:     url,
		Method:  method,
		Headers: headers,
		Body:    body,
		Retries: 2,
	})
}

// ParseJSON is a utility to decode JSON response to a struct
func ParseJSON(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}

// Convenience method for GET
func GetData(url string, headers map[string]string) ([]byte, error) {
	return Send(RequestInput{
		URL:     url,
		Method:  "GET",
		Headers: headers,
		Retries: 1,
	})
}

// Convenience method for POST
func PostData(url string, body []byte, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return Send(RequestInput{
		URL:     url,
		Method:  "POST",
		Headers: headers,
		Body:    body,
		Retries: 1,
	})
}

// Convenience method for PUT
func PutData(url string, body []byte, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return Send(RequestInput{
		URL:     url,
		Method:  "PUT",
		Headers: headers,
		Body:    body,
		Retries: 1,
	})
}

// Convenience method for DELETE
func DeleteData(url string, headers map[string]string) ([]byte, error) {
	return Send(RequestInput{
		URL:     url,
		Method:  "DELETE",
		Headers: headers,
		Retries: 1,
	})
}

// Convenience method for PATCH
func PatchData(url string, body []byte, headers map[string]string) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return Send(RequestInput{
		URL:     url,
		Method:  "PATCH",
		Headers: headers,
		Body:    body,
		Retries: 1,
	})
}

// Convenience method for OPTIONS
func Options(url string, headers map[string]string) ([]byte, error) {
	return Send(RequestInput{
		URL:     url,
		Method:  "OPTIONS",
		Headers: headers,
		Retries: 1,
	})
}
