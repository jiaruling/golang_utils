package lib

import (
	"fmt"
	"strings"

	"github.com/jiaruling/golang_utils/logs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 通用DLTag常量定义
const (
	DLTagUndefind      = "_undef"
	DLTagMySqlFailed   = "_com_mysql_failure"
	DLTagRedisFailed   = "_com_redis_failure"
	DLTagMySqlSuccess  = "_com_mysql_success"
	DLTagRedisSuccess  = "_com_redis_success"
	DLTagThriftFailed  = "_com_thrift_failure"
	DLTagThriftSuccess = "_com_thrift_success"
	DLTagHTTPSuccess   = "_com_http_success"
	DLTagHTTPInfo      = "_com_http_info"
	DLTagHTTPFailed    = "_com_http_failure"
	DLTagTCPFailed     = "_com_tcp_failure"
	DLTagRequestIn     = "_com_request_in"
	DLTagRequestOut    = "_com_request_out"
)

const (
	_dlTag          = "dltag"
	_traceId        = "traceid"
	_spanId         = "spanid"
	_childSpanId    = "cspanid"
	_dlTagBizPrefix = "_com_"
	_dlTagBizUndef  = "_com_undef"
	_dlHint         = "hint"
	_dlHintCode     = "hint_code"
)

type Log struct{}

func InitLog(logFileDir, appName string, maxSize, maxBackups, maxAge int, dev bool) {
	data := &logs.Options{
		LogFileDir: logFileDir,
		AppName:    appName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Config:     zap.Config{},
	}
	data.Development = dev
	logs.InitLogger(data)
}

func GetLog() *Log {
	return &Log{}
}

// 链路日志

func (l *Log) Debug(trace *TraceContext, dltag string, m map[string]interface{}) {
	l.write(trace, dltag, m, zapcore.DebugLevel)
}

func (l *Log) Info(trace *TraceContext, dltag string, m map[string]interface{}) {
	l.write(trace, dltag, m, zapcore.InfoLevel)
}

func (l *Log) Warn(trace *TraceContext, dltag string, m map[string]interface{}) {
	l.write(trace, dltag, m, zapcore.WarnLevel)
}

func (l *Log) Error(trace *TraceContext, dltag string, m map[string]interface{}) {
	l.write(trace, dltag, m, zapcore.ErrorLevel)
}

func (l *Log) DPanic(trace *TraceContext, dltag string, m map[string]interface{}) {
	l.write(trace, dltag, m, zapcore.DPanicLevel)
}

func (l *Log) Panic(trace *TraceContext, dltag string, m map[string]interface{}) {
	l.write(trace, dltag, m, zapcore.PanicLevel)
}

func (l *Log) Fatal(trace *TraceContext, dltag string, m map[string]interface{}) {
	l.write(trace, dltag, m, zapcore.FatalLevel)
}

func (l *Log) write(trace *TraceContext, dltag string, m map[string]interface{}, lvl zapcore.Level) {
	if trace == nil {
		trace = NewTrace()
	}
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	m[_dlTag] = checkDLTag(dltag)
	su := logs.GetSugar()
	switch lvl {
	case zapcore.DebugLevel:
		su.Debug(parseParams(m))
	case zapcore.InfoLevel:
		su.Info(parseParams(m))
	case zapcore.WarnLevel:
		su.Warn(parseParams(m), zap.StackSkip("stack", 3))
	case zapcore.ErrorLevel:
		su.Error(parseParams(m), zap.StackSkip("stack", 3))
	case zapcore.DPanicLevel:
		su.DPanic(parseParams(m), zap.StackSkip("stack", 3))
	case zapcore.PanicLevel:
		su.Panic(parseParams(m), zap.StackSkip("stack", 3))
	case zapcore.FatalLevel:
		su.Fatal(parseParams(m), zap.StackSkip("stack", 3))
	}
}

// 校验dltag合法性
func checkDLTag(dltag string) string {
	if strings.HasPrefix(dltag, _dlTagBizPrefix) || dltag == DLTagUndefind {
		return dltag
	}
	return ""
}

// map格式化为string
func parseParams(m map[string]interface{}) string {
	var (
		dltag string = DLTagUndefind
	)
	if v, have := m[_dlTag]; have {
		if val, ok := v.(string); ok {
			dltag = val
		}
	}
	if v, have := m[_traceId]; have {
		if val, ok := v.(string); ok {
			dltag = dltag + "||" + fmt.Sprintf("%v=%+v", _traceId, val)
		}
	}
	if v, have := m[_spanId]; have {
		if val, ok := v.(string); ok {
			dltag = dltag + "||" + fmt.Sprintf("%v=%+v", _spanId, val)
		}
	}
	if v, have := m[_childSpanId]; have {
		if val, ok := v.(string); ok {
			dltag = dltag + "||" + fmt.Sprintf("%v=%+v", _childSpanId, val)
		}
	}
	if v, have := m[_dlHintCode]; have {
		if val, ok := v.(int); ok {
			dltag = dltag + "||" + fmt.Sprintf("%v=%+v", _dlHintCode, val)
		}
	}
	if v, have := m[_dlHint]; have {
		if val, ok := v.(string); ok {
			dltag = dltag + "||" + fmt.Sprintf("%v=%+v", _dlHint, val)
		}
	}
	for _key, _val := range m {
		if _key == _dlTag || _key == _traceId ||
			_key == _spanId || _key == _childSpanId ||
			_key == _dlHintCode || _key == _dlHint {
			continue
		}
		dltag = dltag + "||" + fmt.Sprintf("%v=%+v", _key, _val)
	}
	// dltag = strings.Trim(fmt.Sprintf("%q", dltag), "\"")
	dltag = strings.Trim(dltag, "\"")
	return dltag
}
