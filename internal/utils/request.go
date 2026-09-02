package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DefaultTimeout is the default HTTP client timeout.
const DefaultTimeout = 30 * time.Second

// DoSimpleRequest performs a GET request with a user-agent header and returns the response.
func DoSimpleRequest(apiURL string) (*http.Response, error) {
	_, _ = fmt.Fprintf(os.Stdout, "Do request. URL = %s\n", apiURL)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 YaBrowser/26.4.0.0 Safari/537.36",
	)

	client := &http.Client{Timeout: DefaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil //nolint:nilerr // nil response treated as empty
		}
		return nil, fmt.Errorf("API request error: %w", err)
	}

	return resp, nil
}
