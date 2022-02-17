# req
A lightweight high-performance request library based on fasthttp

### demo
```go
    // direct request
    client := req.NewClient()
	url := "https://baidu.com"
	res, err := client.GET(url)
	if err != nil {
		return
	}
	fmt.Println(res.Body)

    // request with query or body
	client := req.NewClient()
	client.Config = req.Config{
        Params: map[string]interface{}{
        "q1": 1,
        "q2": 2,
        },
        Json: map[string]interface{}{
        "data1": 1,
        "data2": 2,
        },
    }
	url := "https://api.xxx.com"
	res, err := client.POST(url)
    if err != nil {
        return
    }
    fmt.Println(res.Body)
```