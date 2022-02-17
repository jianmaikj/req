package req

import (
	"encoding/json"
	"errors"
)

type GlobalConfig struct {
	IsLog       bool
	LogHandler  func(uri string, payload []byte, res Response, err error)
	NotOkError  error
	JSONEncoder func(v interface{}) ([]byte, error)
	JSONDecoder func(data []byte, v interface{}) error
}

var DefaultConfig = &GlobalConfig{
	IsLog:       false,
	LogHandler:  nil,
	NotOkError:  errors.New("request not success"),
	JSONEncoder: json.Marshal,
	JSONDecoder: json.Unmarshal,
}
