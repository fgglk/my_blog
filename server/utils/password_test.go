package utils

import (
    "testing"
)

func TestBcryptCheck(t *testing.T) {
    hash := "$2a$10$7bY/QH.C35HNEjIc9RBgPulckgCHhY/Jw3g4WQiPZI49HwuHbUba6"
    password := "12345678"
    
    if !BcryptCheck(password, hash) {
        t.Errorf("密码验证失败，哈希值与明文不匹配")
    } else {
        t.Log("密码验证成功")
    }
}