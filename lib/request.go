package lib

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type ContentType string

var (
	requestMethod = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	FORM ContentType = "application/x-www-form-urlencoded"
	JSON ContentType = "application/json"
)

type Client struct{}

func (c *Client) Request(trace *TraceContext, method, url string, body []byte, msTimeout int, header http.Header, contentType ContentType) (data []byte, err error) {
	startTime := time.Now().UnixNano()
	// 链路日志
	log := GetLog()
	if trace == nil {
		trace = NewTrace()
	}

	t := CreateTips()
	if !t.InStr(method, requestMethod) {
		log.Error(trace, DLTagHTTPFailed, map[string]interface{}{
			"url":       url,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    method,
			"hint":      "<request>: 请求方式不合法！",
			"hint_code": 4001,
		})
		return nil, errors.New("")
	}
	if url == "" {
		log.Error(trace, DLTagHTTPFailed, map[string]interface{}{
			"url":       url,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    method,
			"hint":      "<request>: 请求地址不能为空！",
			"hint_code": 4002,
		})
		return nil, errors.New("")
	}
	url = t.SpliceUrl(url)
	req, err := http.NewRequest(method, t.SpliceUrl(url), bytes.NewReader(body))
	if err != nil {
		log.Error(trace, DLTagHTTPFailed, map[string]interface{}{
			"url":       url,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    method,
			"hint":      "<request>: 初始化请求信息失败",
			"hint_code": 4003,
			"error":     err.Error(),
		})
		return nil, errors.New("")
	}
	// 设置请求头
	if len(header) > 0 {
		req.Header = header
	}
	req.Header.Set("Content-Type", string(contentType))
	req = addTrace2Header(req, trace)

	// 设置客户端请求时间为300毫秒
	client := &http.Client{Timeout: time.Duration(msTimeout) * time.Millisecond}
	response, err := client.Do(req)
	if err != nil {
		log.Error(trace, DLTagHTTPFailed, map[string]interface{}{
			"url":       url,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    method,
			"hint":      "<request>: http请求失败",
			"hint_code": 4004,
			"error":     err.Error(),
		})
		return nil, errors.New("")
	}
	defer response.Body.Close()
	// 获取响应数据
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(trace, DLTagHTTPFailed, map[string]interface{}{
			"url":       url,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    method,
			"hint":      "<request>: 获取响应数据失败",
			"hint_code": 4005,
			"error":     err.Error(),
		})
		return nil, errors.New("")
	}
	log.Info(trace, DLTagHTTPSuccess, map[string]interface{}{
		"url":       url,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    method,
		"result":    string(data),
	})
	return data, nil
}

func (c *Client) RequestSimple(trace *TraceContext, method, url string, body []byte) (data []byte, err error) {
	return c.Request(trace, method, url, body, 500, nil, JSON)
}

func addTrace2Header(request *http.Request, trace *TraceContext) *http.Request {
	traceId := trace.TraceId
	cSpanId := NewSpanId()
	if traceId != "" {
		request.Header.Set("didi-header-rid", traceId)
	}
	if cSpanId != "" {
		request.Header.Set("didi-header-spanid", cSpanId)
	}
	trace.CSpanId = cSpanId
	return request
}