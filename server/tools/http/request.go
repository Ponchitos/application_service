package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func ExecuteHTTPRequest(ctx context.Context, targetURL, method string, data interface{}, timeoutRequest time.Duration) ([]byte, int, error) {
	bodyRequest, err := convertBodyToByte(data)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, method, targetURL, bodyRequest)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: timeoutRequest,
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer func() {
		response.Body.Close()
		client.CloseIdleConnections()
	}()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}

	return responseBody, response.StatusCode, nil
}

func convertBodyToByte(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	jsonMsg, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonMsg), nil
}
