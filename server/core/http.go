package core

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func SendRequest(method, url string, param map[string]string, header map[string]string) ([]byte, error) {
	var request *http.Request
	if http.MethodGet == method {
		if param != nil {
			var s string
			for k, v := range param {
				s += k + "=" + v + "&"
			}
			s, _ = strings.CutSuffix(s, "&")
			url += "?" + s
		}
		request, _ = http.NewRequest(method, url, nil)
	} else {
		marshal, _ := json.Marshal(param)
		request, _ = http.NewRequest(method, url, bytes.NewReader(marshal))
	}
	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}
	return sendHttpRequest(request)
}

func sendHttpRequest(request *http.Request) ([]byte, error) {
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}
