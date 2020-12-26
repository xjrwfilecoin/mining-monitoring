package utils

import (
	"github.com/satori/go.uuid"
)

// 获取唯一的Id
func GetUUID() string {
	id := uuid.NewV4()
	return id.String()
}
