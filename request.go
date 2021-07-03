package agorago

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Request struct {
	appid       string
	secret      string
	credentials string
	client      *http.Client
}

type Option func(*Request)

func SetKeyAndSecret(appid, secret string) Option {
	return func(request *Request) {
		request.appid = appid
		request.secret = secret
	}
}

func SetClient(client *http.Client) Option {
	return func(request *Request) {
		request.client = client
	}
}

func SetCredentials(credentials string) Option {
	return func(request *Request) {
		request.credentials = credentials
	}
}

func NewRequest(opts ...Option) *Request {
	req := &Request{}
	for _, opt := range opts {
		opt(req)
	}
	if len(req.appid) == 0 || len(req.secret) == 0 {
		panic("key or secret not empty !")
	}
	if req.client == nil {
		req.client = http.DefaultClient
	}
	if len(req.credentials) == 0 {
		req.credentials = base64.StdEncoding.EncodeToString([]byte(req.appid + ":" + req.secret))
	}
	return req
}

// 发起请求
func (self *Request) Do(uri, method string,
	body interface{}, r func(req *http.Request),
	resp func(resp *http.Response) error, ret interface{}) error {

	var payload *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(raw)
	}
	request, err := http.NewRequest(method, uri, payload)
	if err != nil {
		return err
	}
	// 这里定义请求参数
	if r != nil {
		r(request)
	}
	// 增加 Authorization header
	request.Header.Add("Authorization", "Basic "+self.credentials)
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	// 发送 HTTP 请求
	res, err := self.client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 这里http返回结果自定义处理逻辑
	if resp != nil {
		return resp(res)
	}
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusOK ||
		res.StatusCode == http.StatusPartialContent {
		if ret != nil {
			err = json.Unmarshal(raw, ret)
			if err != nil {
				return fmt.Errorf("decode:%s uri:%s", string(raw), uri)
			}
		}
		return nil
	}
	return fmt.Errorf("http:%d msg:\\%s uri:%s", res.StatusCode, string(raw), uri)
}
