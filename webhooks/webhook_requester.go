package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type WebhookRequester[RequestType any, ResponseType any] interface {
	Request(url string, args *RequestType, lastTry bool) (*ResponseType, error)
}

type HttpRequester[RequestType any, ResponseType any] struct {
}

func (r *HttpRequester[RequestType, ResponseType]) Request(url string, args *RequestType, lastTry bool) (*ResponseType, error) {
	requestBodyBytes, err := json.Marshal(*args)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if lastTry {
		req.Header.Set("Webhook-Last-Try", "")
	}

	client := &http.Client{Timeout: time.Second * 200}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 { // only 20X are allowed
		// Read the response body
		bodyText := string(body)
		return nil, fmt.Errorf("webhook fail: %s. status: %d. headers: %v. body: %s", url, resp.StatusCode, resp.Header, strings.ReplaceAll(bodyText, "\n", " "))
	}
	var taskWebhookResponse ResponseType

	if len(body) != 0 {
		if err := json.Unmarshal(body, &taskWebhookResponse); err != nil {
			return nil, fmt.Errorf("webhook returned body not in expected format. return empty body or conform to expected format: %w", err)
		}
	}

	return &taskWebhookResponse, nil
}
