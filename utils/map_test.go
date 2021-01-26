package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

var oldStr = `
{
	"minerBalance":"1000FIL",
	"gpuInfo":{
			"0":{
					"name":"0",
					"use":"30"
				},
			"1":{
					"name":"1",
					"use":"20"
				}
	}
}
`
var newStr = `
{
	"minerBalance":"1000FIL",
	"gpuInfo":{
			"0":{
					"name":"0",
					"use":"20"
				},
			"1":{
					"name":"1",
					"use":"20"
				}
	}
}
`

func TestDiffMap(t *testing.T) {
	var new, old map[string]interface{}
	err := json.Unmarshal([]byte(newStr), &new)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(new)
	err = json.Unmarshal([]byte(oldStr), &old)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(old)
	diffMap := DeepDiffMap(old, new)
	fmt.Println(diffMap)
	bytes, err := json.Marshal(diffMap)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(bytes))

}

type Info struct {
	Name string      `json:"name"`
	List interface{} `json:"list"`
}

func TestMap(t *testing.T) {
	var res Info
	err := json.Unmarshal([]byte(oldStr), &res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res)
}
