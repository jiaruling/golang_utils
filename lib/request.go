package lib

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"
)

type ContentType string

var (
	requestMethod             = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	FORM          ContentType = "application/x-www-form-urlencoded"
	JSON          ContentType = "application/json"
	XML           ContentType = "application/xml"
)

type HttpClient struct{}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (c *HttpClient) Request(trace *TraceContext, method, url string, body []byte, msTimeout int, header http.Header, contentType ContentType) (data []byte, err error) {
	startTime := time.Now().UnixNano()
	// 链路日志
	log := GetLog()
	if trace == nil {
		trace = NewTrace()
	}
	log.Info(trace, DLTagHTTPInfo, map[string]interface{}{
		"url":     url,
		"method":  method,
		"request": string(body),
	})
	t := NewTips()
	if !t.InStr(method, requestMethod) {
		return nil, errors.New("请求方式不合法")
	}
	if url == "" {
		return nil, errors.New("请求地址不能为空")
	}
	url = t.SpliceUrl(url)
	req, err := http.NewRequest(method, t.SpliceUrl(url), bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("初始化请求信息失败")
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
		return nil, err
	}
	defer response.Body.Close()
	// 获取响应数据
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Info(trace, DLTagHTTPSuccess, map[string]interface{}{
		"url":       url,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"result":    string(data),
	})
	return data, nil
}

func (c *HttpClient) RequestSimple(trace *TraceContext, method, url string, body []byte) (data []byte, err error) {
	return c.Request(trace, method, url, body, 2000, nil, JSON)
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
