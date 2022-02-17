package req

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {

	client := NewClient()
	client.Config = Config{
		Params: map[string]interface{}{
			"q1": 1,
			"q2": 2,
		},
		Json: map[string]interface{}{
			"data1": 1,
			"data2": 2,
		},
	}
	url := "https://baidu.com"
	res, err := client.GET(url)
	if err != nil {
		return
	}
	fmt.Println(res.Body)
}
