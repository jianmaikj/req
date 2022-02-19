package req

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	//修改默认配置
	DefaultConfig.IsLog = true

	client := NewClient()
	client.Config = Config{
		Params: map[string]interface{}{
			"q1": 1,
			"q2": 2,
		},
		Data: map[string]interface{}{
			"data1": 1,
			"data2": 2,
		},
		Timeout: 30 * time.Second,
	}

	//client.Data = map[string]interface{}{
	//	"data1": 1,
	//	"data2": 2,
	//}
	url := "https://baidu.com"
	res, err := client.POST(url)
	if err != nil {
		return
	}
	fmt.Println(res.String())
}
