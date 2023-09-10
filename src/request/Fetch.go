package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Fetch(url string, body []byte) (string, []byte, float64, string, error) {
	start := time.Now()
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", nil, 0, "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, 0, "", fmt.Errorf("while reading %s: %v", url, err)
	}

	sces := time.Since(start).Seconds()
	log := fmt.Sprintf("%.2fs %7d %s", sces, len(respBody), url)
	return log, respBody, sces, start.Format(time.DateTime), nil
}
