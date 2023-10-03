package req

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	//修改默认配置
	DefaultConfig.IsLog = true

	client := NewClient()
	url := "https://www.baidu.com/test?i=9"
	res, err := client.POST(url, &Config{
		Params: map[string]interface{}{
			"q1": 1,
			"q2": 2,
		},
		Data: map[string]interface{}{
			"data1": 1,
			"data2": 2,
		},
		Timeout: 30,
	}).Do()
	if err != nil {
		return
	}
	fmt.Println(res.String())
}