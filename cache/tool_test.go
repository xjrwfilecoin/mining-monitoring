package cache

import (
	"fmt"
	"testing"
)

func TestTool(t *testing.T) {
	var res1 []map[string]interface{}
	var res2 []map[string]interface{}
	equal := DeepEqual(res1, res2)
	fmt.Println(equal)
}
