package req

func (body Body) Map() (map[string]interface{}, error) {
	var tempMap map[string]interface{}
	if err := DefaultConfig.JSONDecoder(body, &tempMap); err != nil {
		return nil, err
	}
	return tempMap, nil
}
func (body Body) Struct(s interface{}) error {
	if err := DefaultConfig.JSONDecoder(body, &s); err != nil {
		return err
	}
	return nil
}

func (body Body) String() string {
	return string(body)
}

func (res Response) IsOk() bool {
	status := res.Status
	return status == 200 || (status > 200 && status < 300)
}

func (res Response) Error(msg ...string) (err error) {
	if !res.IsOk() {
		err = DefaultConfig.NotOkError
	}
	return
}
