package utils

import (
    "strconv"
    "fmt"
)

// StringToUint 将字符串转换为uint类型
// 如果转换失败，返回错误信息
func StringToUint(s string) (uint, error) {
    id, err := strconv.ParseUint(s, 10, 64)
    if err != nil {
        return 0, fmt.Errorf("无效的ID: %s, 错误: %w", s, err)
    }
    return uint(id), nil
}