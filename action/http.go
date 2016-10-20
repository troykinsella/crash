package action

import (
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
)

type Http struct {
	config *ActionConfig
}

func (h *Http) Run() (*Result, error) {
	url := h.config.Params.GetString("url")
	if url == "" {
		return nil, fmt.Errorf("url parameter required")
	}

	method := h.config.Params.GetString("method")
	if method == "" {
		method = "GET"
	} else {
		method = strings.ToUpper(method)
	}

	client := &http.Client{}

	h.config.Log.Debugf("%s %s", method, url)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	result, err := h.genResult(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *Http) genResult(resp *http.Response) (*Result, error) {
	data := make(map[string]interface{})
	data["status-code"] = resp.StatusCode
	data["headers"] = resp.Header

	h.config.Log.Infof("%s %s -> %d", resp.Request.Method, resp.Request.URL.String(), resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data["raw-body"] = bodyBytes
	data["body"] = string(bodyBytes)

	return &Result{
		Data:    data,
	}, nil
}

func NewHttp(config *ActionConfig) *Http {
	return &Http{
		config: config,
	}
}
