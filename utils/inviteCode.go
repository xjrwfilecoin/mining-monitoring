package utils

import (
	"time"
)

var seeds = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

var saltFlag int64 = 1575445924991042800

// 申请邀请码,转换为32进制的格式
func GenInviteCode() string {
	uid := time.Now().UnixNano() - saltFlag
	maxNum := int64(len(seeds))
	code := ""
	for uid != 0 {
		mod := uid % maxNum
		uid = uid / maxNum
		code = string(seeds[mod]) + code
	}
	for len(code) < 6 {
		code = "0" + code
	}
	return code
}
