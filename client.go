package grpc_captcha_validation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func clientRequest(address, method string, parameter map[string]string, headers map[string]string, response any) error {
	client := http.DefaultClient

	req, err := http.NewRequest(method, address, nil)
	if err != nil {
		return ERR_FAILED_CREATE_REQUEST
	}

	query := url.Values{}
	for k, v := range parameter {
		query[k] = []string{v}
	}

	req.URL.RawQuery = query.Encode()

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("secret: failed http do request, got error %s", err.Error())
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("secret: failed to read body, got error %s", err.Error())
	}

	if err := json.Unmarshal(b, response); err != nil {
		return fmt.Errorf("secret: failed to unmarshal data to response, got error %s", err.Error())
	}

	return nil
}
