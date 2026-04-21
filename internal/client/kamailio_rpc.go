package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kolo/xmlrpc"
)

type KamailioJSONClient struct {
	Endpoint string
	Timeout  time.Duration
}

func (k *KamailioJSONClient) Invoke(method string, params ...interface{}) (interface{}, error) {
	client := resty.New().SetTimeout(k.Timeout)
	var finalParams interface{} = params
	if params == nil {
		finalParams = []interface{}{}
	}

	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  finalParams,
		"id":      time.Now().Unix(),
	}

	var result map[string]interface{}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&result).
		Post(k.Endpoint)

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("kamailio response error: %s", resp.String())
	}
	return result["result"], nil
}

type SupervisorClient struct {
	Endpoint string
}

func (s *SupervisorClient) ControlProcess(action string, name string) (bool, error) {
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", s.Endpoint)
		},
	}

	client, err := xmlrpc.NewClient("http://localhost/RPC2", transport)
	if err != nil {
		return false, err
	}
	var res bool
	method := fmt.Sprintf("supervisor.%sProcess", action)
	err = client.Call(method, []interface{}{name, true}, &res)
	return res, err
}

func (k *KamailioJSONClient) InvokeNativeHTTP(path string, queryParams map[string]string) (string, error) {
        client := resty.New().SetTimeout(k.Timeout)

        // 拼接完整的 URL，例如 http://127.0.0.1:8091/api/startSuccessful
        // 注意：这里需要处理 Endpoint 结尾是否有斜杠的问题
        fullURL := fmt.Sprintf("http://%s%s", "localhost:8091", path)

        resp, err := client.R().
                SetQueryParams(queryParams).
                Post(fullURL)

        if err != nil {
                return "", err
        }
        if resp.IsError() {
                return "", fmt.Errorf("kamailio native api error: %s", resp.Status())
        }

        return resp.String(), nil
}

