package shellParsing

import (
	"testing"
)

type User struct {
	Name string
}

func Test(t *testing.T) {
	param := make(map[string]map[string]interface{})
	var res []User
	res = append(res, User{Name: "test"})

}
