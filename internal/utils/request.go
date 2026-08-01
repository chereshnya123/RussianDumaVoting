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
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 YaBrowser/26.4.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Minute}
	resp, err := client.Do(req)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("API request error: %w", err)
	}

	return resp, nil
}
