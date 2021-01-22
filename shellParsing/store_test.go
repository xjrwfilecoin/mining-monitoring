package shellParsing

import (
	"fmt"
	"testing"
)

func Test(t *testing.T){
	data := CmdData{Data:IoInfo{
		ReadIO:  "",
		WriteIO: "",
	}}
	fmt.Println(data)
}
