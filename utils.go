package req

import (
	"github.com/jianmaikj/convert"
	"github.com/valyala/fasthttp"
	"reflect"
)

// Log uri:请求的完整路径,payload:请求负载,res:响应,err:请求错误,ErrNoFreeConns is returned if all DefaultMaxConnsPerHost connections to the requested host are busy.
func Log(uri string, payload []byte, res *Response, err error) {
	if res == nil {
		DefaultConfig.LogHandler(uri, payload, Response{}, err)
	} else {
		DefaultConfig.LogHandler(uri, payload, *res, err)
	}
}

func GetQueryString(params map[string]interface{}) []byte {
	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)
	for k, v := range params {
		args.Set(k, convert.Str(v))
	}
	str := args.QueryString()
	return str
}

func IsNil(v interface{}) bool {
	vi := reflect.ValueOf(v)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func IsValuePtr(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr
}

func IsValueNil(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	return false
}

func IsValueSlice(v reflect.Value) bool {
	return v.Kind() == reflect.Slice
}

func IsValueArray(v reflect.Value) bool {
	return v.Kind() == reflect.Array
}