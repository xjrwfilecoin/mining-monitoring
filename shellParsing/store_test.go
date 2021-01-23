package shellParsing

import (
	"fmt"
	"testing"
)

func Test(t *testing.T){
	data := CmdData{Data:IoInfo{
		DiskR:   "",
		WriteIO: "",
	}}
	fmt.Println(data)
}
