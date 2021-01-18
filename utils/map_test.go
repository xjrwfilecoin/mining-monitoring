package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

var test01 = `
{
	"name":"test01",
	"list":[
								   { "hostName":"worker01",
									"id":"67081a0d",
									"sector":"34",
									"state":"running",
									"task":"PC1",
									"time":"6m53.8s",
									"worker":"84af77dc"
								},
								{
									"hostName":"worker01",
									"id":"05d6450f",
									"sector":"29",
									"state":"running",
									"task":"PC1",
									"time":"2h37m1.3s",
									"worker":"84af77dc"
								}
	]
}
`

type Info struct {
	Name string      `json:"name"`
	List interface{} `json:"list"`
}

func TestMap(t *testing.T) {
	var res Info
	err := json.Unmarshal([]byte(test01), &res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res)
}

var newSrc = `
	{
		"name":"test01",
	    "age":10,
		"info":{
			"score":11,
			"title":"测试"
		}
	}
`

var oldSrc = `
	{
		"name":"test01",
	    "age":11,
		"info":{
			"score":10,
			"title":"测试"
		}
	}
`

func TestDiffMap(t *testing.T) {
	var new, old map[string]interface{}
	err := json.Unmarshal([]byte(newSrc), &new)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(new)
	err = json.Unmarshal([]byte(oldSrc), &old)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(old)
	diffMap := DiffMap(old, new)
	fmt.Println(diffMap)

}
