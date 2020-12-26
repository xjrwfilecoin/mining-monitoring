package utils

import (
	"fmt"
	"testing"
)

func TestInviteCode(t *testing.T) {
	code := GenInviteCode()
	fmt.Println(code)
}
