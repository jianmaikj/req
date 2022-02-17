package req

func (body ResponseBody) Map() (map[string]interface{}, error) {
	var tempMap map[string]interface{}
	if err := DefaultConfig.JSONDecoder(body, &tempMap); err != nil {
		return nil, err
	}
	return tempMap, nil
}
func (body ResponseBody) Struct(s interface{}) error {
	if err := DefaultConfig.JSONDecoder(body, &s); err != nil {
		return err
	}
	return nil
}

func (body ResponseBody) String() string {
	return string(body)
}

func (res Response)IsOk() bool {
	return res.Status==200
}

func (res Response)Error(msg ...string) (err error) {
	if !res.IsOk(){
		err = DefaultConfig.NotOkError
	}
	return
}