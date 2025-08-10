package utils

import "github.com/gofrs/uuid"

// GenerateUUID 生成UUID字符串
func GenerateUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}
