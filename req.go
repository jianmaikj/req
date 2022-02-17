package req

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type ResponseBody json.RawMessage

type BasicAuth struct {
	Username, Password string
}
type Client struct {
	BasicAuth *BasicAuth
	Config    Config
}

type Config struct {
	Params  map[string]interface{}
	Form    map[string]interface{}
	Json    interface{}
	Headers map[string]string
	Type    string
}

type Response struct {
	Status int
	Body   ResponseBody
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
	uri := url
	if params != nil {
		queryString := GetQueryString(params)
		uri = uri + "?" + string(queryString)
	}
	req.SetRequestURI(uri)
	if c.BasicAuth != nil {
		req.URI().SetUsername(c.BasicAuth.Username)
		req.URI().SetPassword(c.BasicAuth.Password)
	}
	return req
}

func (c *Client) GET(url string) (res Response, err error) {
	req := c.NewReq(url, "GET")
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if DefaultConfig.IsLog{
		defer Log(req.URI().String(), []byte(``), res, err)
	}
	if err = fasthttp.Do(req, resp); err != nil {
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
	form := c.Config.Form
	jsonData := c.Config.Json
	var payload []byte
	if form != nil {
		payload = GetQueryString(form)
		req.Header.SetContentType("application/x-www-form-urlencoded")
	} else {
		req.Header.SetContentType("application/json")
		payload, _ = DefaultConfig.JSONEncoder(jsonData)
	}
	req.SetBody(payload)
	if err := fasthttp.Do(req, resp); err != nil {
		go Log(req.URI().String(), payload, res, err)
		return Response{}, err
	}

	body := resp.Body()
	status := resp.StatusCode()
	res = Response{
		status,
		body,
	}
	if status != 200 {
		go Log(req.URI().String(), payload, res, nil)
	}
	return
}

func (c *Client) PATCH(url string) (res Response, err error) {
	req := c.NewReq(url, "PATCH")
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源
	form := c.Config.Form
	jsonData := c.Config.Json
	var payload []byte
	if form != nil {
		payload = GetQueryString(form)
		req.Header.SetContentType("application/x-www-form-urlencoded")
	} else {
		req.Header.SetContentType("application/json")
		payload, _ = DefaultConfig.JSONEncoder(jsonData)
	}
	req.SetBody(payload)
	if err := fasthttp.Do(req, resp); err != nil {
		go Log(req.URI().String(), payload, res, err)
		return Response{}, err
	}

	body := resp.Body()
	status := resp.StatusCode()
	res = Response{
		status,
		body,
	}
	if status != 200 {
		go Log(req.URI().String(), payload, res, nil)
	}
	return
}
