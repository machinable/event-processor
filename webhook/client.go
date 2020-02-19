package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Client wraps the http client and logger
type Client struct {
	client *http.Client
	logger *logrus.Logger
	mux    sync.Mutex
}

// NewClient returns a pointer to a new `Client` with the provided logger
func NewClient(logger *logrus.Logger) *Client {
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 10 * time.Second,
			},
		},
		logger: logger,
	}
}

// PostHook attempts to POST the provided payload to the configured web hook within the event
func (c *Client) PostHook(event *HookEvent) *HookResult {
	// lock http client
	c.mux.Lock()
	defer c.mux.Unlock()

	result := &HookResult{
		WebHookID: event.Hook.ID,
		ProjectID: event.Hook.ProjectID,
	}

	c.logger.WithFields(logrus.Fields{"hook_id": event.Hook.ID}).Info("sending payload for hook")

	payload := event.Format()
	// marshal payload
	payloadBytes, berr := json.Marshal(payload)
	if berr != nil {
		c.logger.Error(berr)
		result.ErrorMessage = berr.Error()
		return result
	}

	// prepare http request
	req, err := http.NewRequest("POST", event.Hook.HookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.logger.Error(err)
		result.ErrorMessage = err.Error()
		return result
	}

	// add configured headers
	for _, header := range event.Hook.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	// do request
	requestStart := time.Now()
	resp, rerr := c.client.Do(req)
	respTime := time.Now().Sub(requestStart).Seconds() * 1000

	// set response time to result
	result.ResponseTime = respTime

	if rerr != nil {
		c.logger.Error(rerr)
		result.ErrorMessage = rerr.Error()
		return result
	}

	// set response status code to result
	result.StatusCode = resp.StatusCode

	// TODO: Save headers and perhaps response body?

	// log "successful" request
	c.logger.WithFields(logrus.Fields{"result": result}).Info("hook sent")

	return result
}
