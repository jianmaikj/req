package req

import (
	"encoding/json"
	"fmt"
	"github.com/jianmaikj/convert"
	"github.com/valyala/fasthttp"
	"io"
	"reflect"
	"time"
)

type Body []byte
type Args = fasthttp.Args

// var client = &fasthttp.Client{}
type BasicAuth struct {
	Username, Password string
}

type Client struct {
	*fasthttp.Client
	BasicAuth *BasicAuth
}

type Config struct {
	Params      map[string]interface{}
	Form        map[string]interface{}
	Data        interface{}
	Json        json.RawMessage
	Headers     map[string]string
	Timeout     time.Duration //second
	MaxRedirect int
	GetHeaders  bool
	BodyWriteTo io.Writer
}

type Req struct {
	*fasthttp.Request
	client *Client
	config *Config
}

type Response struct {
	Status int
	Body
	Header *fasthttp.ResponseHeader
}

func NewClient() *Client {
	return &Client{&fasthttp.Client{}, nil}
}

type Reqs interface {
	POST() interface{}
	GET() interface{}
}

func (c *Client) NewReq(url string, method string, config *Config) *Req {
	req := fasthttp.AcquireRequest()

	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	r := &Req{
		Request: req,
	}
	if config != nil {
		if config.Headers != nil {
			r.AppendHeaders(config.Headers)
		}
		params := config.Params
		if params != nil {
			r.AppendParams(params)
		}
		if config.Timeout > 0 {
			req.SetTimeout(config.Timeout * time.Second)
		}
	} else {
		config = &Config{}
	}

	if c.BasicAuth != nil {
		req.URI().SetUsername(c.BasicAuth.Username)
		req.URI().SetPassword(c.BasicAuth.Password)
	}
	r.config = config
	r.client = c
	return r
}

func (r *Req) AppendHeaders(headers map[string]string) *Req {
	if len(headers) == 0 {
		return r
	}
	for k, v := range headers {
		r.Header.Add(k, v)
	}
	return r
}
func (r *Req) AddHeader(k, v string) *Req {
	if k == "" {
		return r
	}
	r.Header.Add(k, v)
	return r
}

func (r *Req) AppendParams(params map[string]interface{}) *Req {
	if len(params) == 0 {
		return r
	}
	for k, v := range params {
		value := reflect.ValueOf(v)
		if IsValuePtr(value) {
			if value.IsNil() {
				continue
			} else {
				value = reflect.Indirect(value)
			}
		}
		if IsValueSlice(value) || IsValueArray(value) {
			for i := 0; i < value.Len(); i++ {
				ele := value.Index(i)
				r.URI().QueryArgs().Add(k, convert.ReflectValue2Str(ele))
			}
		} else {
			r.URI().QueryArgs().Add(k, convert.ReflectValue2Str(value))
		}
	}
	//r.AppendParams(params)
	return r
}
func (r *Req) AddParam(k string, v interface{}) *Req {
	if k == "" {
		return r
	}
	r.URI().QueryArgs().Add(k, convert.Str(reflect.ValueOf(v)))
	return r
}

func (r *Req) AppendForms(form map[string]interface{}) *Req {
	if len(form) == 0 {
		return r
	}
	for k, v := range form {
		r.PostArgs().Add(k, convert.Str(reflect.ValueOf(v)))
	}
	return r
}
func (r *Req) AddForm(k string, v interface{}) *Req {
	if k == "" {
		return r
	}
	r.PostArgs().Add(k, convert.Str(reflect.ValueOf(v)))
	return r
}

func (c *Client) GET(url string, config ...*Config) (req *Req) {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	}
	req = c.NewReq(url, "GET", cfg)
	return
}

func (c *Client) POST(url string, config *Config) (r *Req) {
	r = c.NewReq(url, "POST", config)
	var payload []byte
	form := config.Form
	data := config.Data
	contentType := "application/json"
	if form != nil {
		contentType = "application/x-www-form-urlencoded"
		payload = GetQueryString(form)
	} else if data != nil {
		payload, _ = DefaultConfig.JSONEncoder(data)
	} else {
		payload = config.Json
	}
	r.Header.SetContentType(contentType)
	r.SetBody(payload)
	return
}

func (c *Client) PATCH(url string, config *Config) (r *Req) {
	r = c.NewReq(url, "PATCH", config)
	var payload []byte
	data := config.Data
	contentType := "application/json"
	if data != nil {
		payload, _ = DefaultConfig.JSONEncoder(data)
	} else {
		payload = config.Json
	}
	r.Header.SetContentType(contentType)
	r.SetBody(payload)
	return
}

func (c *Client) PUT(url string, config *Config) (r *Req) {
	r = c.NewReq(url, "PUT", config)
	data := config.Data
	contentType := "application/json"
	var payload []byte
	if data != nil {
		payload, _ = DefaultConfig.JSONEncoder(data)
	} else {
		payload = config.Json
	}
	r.Header.SetContentType(contentType)
	r.SetBody(payload)
	return
}

func (r *Req) Do() (res *Response, err error) {
	resp := fasthttp.AcquireResponse()
	req := r.Request
	defer func() {
		if DefaultConfig.IsLog {
			Log(req.URI().String(), req.Body(), res, err)
		}
		fasthttp.ReleaseResponse(resp) // 用完需要释放资源
		fasthttp.ReleaseRequest(req)
	}()
	c := r.client
	config := r.config
	if config.MaxRedirect != 0 {
		err = c.DoRedirects(req, resp, config.MaxRedirect)
	} else {
		err = c.Do(req, resp)
	}
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	if config.BodyWriteTo != nil {
		fmt.Println("c.BodyWriteTo:", config.BodyWriteTo)
		err = resp.BodyWriteTo(config.BodyWriteTo)
		if err != nil {
			fmt.Println("BodyWriteTo err:", err)
			return
		}
	}
	body := resp.Body()
	status := resp.StatusCode()

	var respHeaders *fasthttp.ResponseHeader
	if config.GetHeaders {
		respHeaders = &fasthttp.ResponseHeader{}
		resp.Header.CopyTo(respHeaders)
	}
	res = &Response{
		status,
		body,
		respHeaders,
	}
	return
}