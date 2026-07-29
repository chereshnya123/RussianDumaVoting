package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func DoSimpleRequest(apiUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("API request error: %w", err)
	}

	return resp, nil
}
