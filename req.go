package req

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
)

type Body []byte

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
	params := c.Config.Params
	req.SetRequestURI(url)
	if params != nil {
		AddQueryArgs(req, params)
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

func (c *Client) GET(url string) (res Response, err error) {
	req := c.NewReq(url, "GET")
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if DefaultConfig.IsLog {
		defer Log(req.URI().String(), []byte(``), res, err)
	}
	if err = do(req, resp, c.Timeout); err != nil {
		return
	}

	body := resp.Body()
	status := resp.StatusCode()
	res = Response{
		status,
		body,
	}
	return
}

func (c *Client) POST(url string) (res Response, err error) {
	req := c.NewReq(url, "POST")
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源
	form := c.Form
	data := c.Data
	var payload []byte
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
	if DefaultConfig.IsLog {
		defer Log(req.URI().String(), payload, res, err)
	}
	if err = do(req, resp, c.Timeout); err != nil {
		return
	}
	body := resp.Body()
	status := resp.StatusCode()
	res = Response{
		status,
		body,
	}
	return
}

func (c *Client) PATCH(url string) (res Response, err error) {
	req := c.NewReq(url, "PATCH")
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源
	data := c.Data
	var payload []byte
	contentType := "application/json"
	if data != nil {
		payload, _ = DefaultConfig.JSONEncoder(data)
	} else {
		payload = c.Json
	}
	req.Header.SetContentType(contentType)
	req.SetBody(payload)
	if DefaultConfig.IsLog {
		defer Log(req.URI().String(), payload, res, err)
	}
	if err = do(req, resp, c.Timeout); err != nil {
		return
	}
	body := resp.Body()
	status := resp.StatusCode()
	res = Response{
		status,
		body,
	}
	return
}

func (c *Client) PUT(url string) (res Response, err error) {
	req := c.NewReq(url, "PUT")
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源
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
	if DefaultConfig.IsLog {
		defer Log(req.URI().String(), payload, res, err)
	}
	if err = do(req, resp, c.Timeout); err != nil {
		return
	}

	body := resp.Body()
	status := resp.StatusCode()
	res = Response{
		status,
		body,
	}
	return
}

func do(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) (err error) {
	fmt.Println("timeout>>",timeout)
	if timeout > 0 {
		err = fasthttp.DoTimeout(req, resp, timeout)
	} else {
		err = fasthttp.Do(req, resp)
	}
	return
}
