package utils

import (
	"math/rand"
	"server/global"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

// SendEmailCode 发送邮箱验证码，支持多种场景
func SendEmailCode(toEmail, code, subject, usage string) error {
	// 创建邮件消息
	msg := gomail.NewMessage()
	msg.SetHeader("From", global.Config.Email.From)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", "您的"+usage+"验证码为: <b>"+code+"</b>，有效期5分钟")

	// 创建SMTP客户端
	dialer := gomail.NewDialer(
		global.Config.Email.Host,
		global.Config.Email.Port,
		global.Config.Email.From,
		global.Config.Email.Secret,
	)
	dialer.SSL = global.Config.Email.IsSsl

	// 发送邮件
	return dialer.DialAndSend(msg)
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min+1)
}

// GenerateEmailCode 生成6位数字邮箱验证码
func GenerateEmailCode() string {
	// 生成6位随机数字
	return strconv.Itoa(RandInt(100000, 999999))
}

// StoreEmailCodeInRedis 将邮箱验证码存储到Redis
func StoreEmailCodeInRedis(email, code string) error {
	// 设置邮箱验证码有效期为5分钟
	expireTime := 5 * time.Minute
	key := "email_code:" + email
	return global.Redis.Set(key, code, expireTime).Err()
}

// VerifyEmailCode 验证邮箱验证码
func VerifyEmailCode(email, code string) bool {
	key := "email_code:" + email
	storedCode, err := global.Redis.Get(key).Result()
	if err != nil {
		return false
	}

	// 验证成功后删除验证码
	global.Redis.Del(key)

	return storedCode == code
}
