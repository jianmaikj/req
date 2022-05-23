package req

import (
	"encoding/json"
	"github.com/jianmaikj/convert"
	"github.com/valyala/fasthttp"
	"time"
)

type Body []byte
type Args = fasthttp.Args

type BasicAuth struct {
	Username, Password string
}

type Client struct {
	BasicAuth *BasicAuth
	Config
}

type Config struct {
	Params  map[string]interface{}
	Form    map[string]interface{}
	Data    interface{}
	Json    json.RawMessage
	Headers map[string]string
	Timeout time.Duration
}

type Response struct {
	Status int
	Body
}

func NewClient() *Client {
	return &Client{}
}

type Reqs interface {
	POST() interface{}
	GET() interface{}
}

func (c *Client) NewReq(url string, method string) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	params := c.Params
	req.SetRequestURI(url)
	if params != nil {
		for k, v := range params {
			req.URI().QueryArgs().Add(k, convert.Str(v))
		}
	}
	if c.BasicAuth != nil {
		req.URI().SetUsername(c.BasicAuth.Username)
		req.URI().SetPassword(c.BasicAuth.Password)
	}
	return req
}

func (c *Client) AppendHeaders(headers map[string]string) {
	if len(headers) == 0 {
		return
	}
	for k, v := range headers {
		c.Headers[k] = v
	}
	return
}
func (c *Client) AddHeaders(k, v string) {
	if k == "" {
		return
	}
	c.Headers[k] = v
	return
}

func (c *Client) AppendParams(params map[string]interface{}) {
	if len(params) == 0 {
		return
	}
	for k, v := range params {
		c.Params[k] = v
	}
	return
}
func (c *Client) AddParams(k string, v interface{}) {
	if k == "" {
		return
	}
	c.Params[k] = v
	return
}

func (c *Client) AppendForm(form map[string]interface{}) {
	if len(form) == 0 {
		return
	}
	for k, v := range form {
		c.Form[k] = v
	}
	return
}
func (c *Client) AddForm(k string, v interface{}) {
	if k == "" {
		return
	}
	c.Form[k] = v
	return
}

func (c *Client) GET(url string) (res *Response, err error) {
	req := c.NewReq(url, "GET")
	resp := fasthttp.AcquireResponse()
	res, err = c.do(req, resp)
	return
}

func (c *Client) POST(url string) (res *Response, err error) {
	req := c.NewReq(url, "POST")
	resp := fasthttp.AcquireResponse()
	var payload []byte
	form := c.Form
	data := c.Data
	contentType := "application/json"
	if form != nil {
		contentType = "application/x-www-form-urlencoded"
		payload = GetQueryString(form)
	} else if data != nil {
		payload, _ = DefaultConfig.JSONEncoder(data)
	} else {
		payload = c.Json
	}
	req.Header.SetContentType(contentType)
	req.SetBody(payload)
	res, err = c.do(req, resp)
	return
}

func (c *Client) PATCH(url string) (res *Response, err error) {
	req := c.NewReq(url, "PATCH")
	resp := fasthttp.AcquireResponse()
	var payload []byte
	data := c.Data
	contentType := "application/json"
	if data != nil {
		payload, _ = DefaultConfig.JSONEncoder(data)
	} else {
		payload = c.Json
	}
	req.Header.SetContentType(contentType)
	req.SetBody(payload)
	res, err = c.do(req, resp)
	return
}

func (c *Client) PUT(url string) (res *Response, err error) {
	req := c.NewReq(url, "PUT")
	resp := fasthttp.AcquireResponse()
	data := c.Data
	contentType := "application/json"
	var payload []byte
	if data != nil {
		payload, _ = DefaultConfig.JSONEncoder(data)
	} else {
		payload = c.Json
	}
	req.Header.SetContentType(contentType)
	req.SetBody(payload)
	//if DefaultConfig.IsLog {
	//	defer Log(req.URI().String(), payload, res, err)
	//}
	res, err = c.do(req, resp)
	return
}

func (c *Client) do(req *fasthttp.Request, resp *fasthttp.Response) (res *Response, err error) {
	defer func() {
		if DefaultConfig.IsLog {
			Log(req.URI().String(), req.Body(), res, err)
		}
		fasthttp.ReleaseResponse(resp) // 用完需要释放资源
		fasthttp.ReleaseRequest(req)
	}()

	headers := c.Headers
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	if c.Timeout > 0 {
		err = fasthttp.DoTimeout(req, resp, c.Timeout)
	} else {
		err = fasthttp.Do(req, resp)
	}
	if err != nil {
		return
	}
	body := resp.Body()
	status := resp.StatusCode()
	res = &Response{
		status,
		body,
	}
	return
}
