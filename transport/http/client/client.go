package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpSend struct {
	proxy       string
	link        string
	method_type string
	header      map[string]string
	body        string
}

func NewHttpSend(link string) HttpSend {
	return HttpSend{
		link:        link,
		method_type: http.MethodGet,
	}
}

func (h *HttpSend) SetBody(body string) {
	h.body = body
}

func (h *HttpSend) SetProxy(proxy string) {
	h.proxy = proxy
}

func (h *HttpSend) SetHeader(header map[string]string) {
	h.header = header
}
func (h *HttpSend) SetMethod(method string) {
	h.method_type = method
}

func (h *HttpSend) Do() ([]byte, error) {
	return h.send()
}

func (h *HttpSend) send() ([]byte, error) {
	var (
		req       *http.Request
		resp      *http.Response
		client    http.Client
		send_data string
		err       error
	)

	if len(h.body) > 0 {
		send_data = h.body
	}

	client = http.Client{
		Timeout: 60 * time.Second,
	}
	if len(h.proxy) > 0 {
		uri, err := url.Parse(h.proxy)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
		}
	}

	req, err = http.NewRequest(h.method_type, h.link, strings.NewReader(send_data))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	for k, v := range h.header {
		if strings.ToLower(k) == "host" {
			req.Host = v
		} else {
			req.Header.Add(k, v)
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong url: %s; but status code: %d", h.link, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
