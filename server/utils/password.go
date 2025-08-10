package utils

import "golang.org/x/crypto/bcrypt"

// BcryptHash 使用bcrypt算法对密码进行加密
func BcryptHash(password string) string {
	// 使用默认成本因子(10)进行加密
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes)
}

// BcryptCheck 验证密码与哈希值是否匹配
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
