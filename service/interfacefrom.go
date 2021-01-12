package service

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type IValid interface {
	Valid() error
}

type BaseFrom struct {
	Url   string      `json:"url"`
	Data  interface{} `json:"body"`
	MsgId string      `json:"msgId"`
}

func (b *BaseFrom) Valid() error {
	return nil
}

type MinerInfoForm struct {
	MinerId string
}

func (m *MinerInfoForm) Valid() error {
	if len(m.MinerId) == 0 {
		return fmt.Errorf("minerId is empty \n")
	}
	return nil
}

func BindJson(src string, obj IValid) error {

	if len(src) == 0 {
		return fmt.Errorf("data is zero \n")
	}
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return fmt.Errorf("obj is not ptr \n")
	}
	err := json.Unmarshal([]byte(src), obj)
	if err != nil {
		return fmt.Errorf("json unmarshal error %v \n", err)
	}
	return obj.Valid()
}
