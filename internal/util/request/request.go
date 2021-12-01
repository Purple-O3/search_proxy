package request

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//TODO:长连接池
func Get(ctx context.Context, url string, timeout int) ([]byte, error) {
	return request(ctx, "GET", url, "", []byte{}, timeout)
}

func Post(ctx context.Context, url string, contentType string, body []byte, timeout int) ([]byte, error) {
	return request(ctx, "POST", url, contentType, body, timeout)
}

func request(ctx context.Context, method string, url string, contentType string, body []byte, timeout int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()
	client := &http.Client{}
	var bodyReader io.Reader
	if len(body) != 0 {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bodyReader)
	if err != nil {
		return nil, err
	}
	if method != "GET" {
		if contentType == "" {
			contentType = "application/json"
		}
		req.Header.Set("Content-Type", contentType)
	}
	trackid, ok := ctx.Value("trackid").(uint64)
	if ok {
		req.Header.Set("X-Trackid", strconv.FormatUint(trackid, 10))
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respByte, nil
}
