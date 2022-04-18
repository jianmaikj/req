package req

import (
	"github.com/jianmaikj/convert"
	"github.com/valyala/fasthttp"
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
