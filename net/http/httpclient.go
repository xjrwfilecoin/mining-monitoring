package httpsvr

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

/**
向lotus节点获取数据的疯转 client
*/

var client *http.Client

func NewHttpClient() *http.Client {
	if client == nil {
		client = http.DefaultClient
	}
	return client
}
func NewRequest(host, data, token string) (string, error) {
	request, err := http.NewRequest("POST", host, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)
	res, err := NewHttpClient().Do(request)
	if err != nil {
		return "", err
	}
	if res.StatusCode < 200 || res.StatusCode > 300 {
		return "", errors.New("response fail ")
	}
	if res == nil {
		return "", errors.New("response is nil")
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}




