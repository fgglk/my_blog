package utils

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"

	"server/global"
)

var (
	captchaInstance *base64Captcha.Captcha
	once            sync.Once
)

// 初始化验证码驱动
func initCaptcha() {
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)

	// 使用默认内存存储
	store := base64Captcha.DefaultMemStore

	captchaInstance = base64Captcha.NewCaptcha(driver, store)
}

// GenerateCaptcha 生成验证码
func GenerateCaptcha() (string, string, error) {
	once.Do(initCaptcha)

	// 生成验证码
	id, b64s, code, err := captchaInstance.Generate()
	if err != nil {
		global.ZapLog.Error("生成验证码失败", zap.Error(err))
		return "", "", err
	}

	// 从内存存储获取验证码文本
	code = captchaInstance.Store.Get(id, false)
	if code == "" {
		return "", "", fmt.Errorf("验证码不存在或已过期")
	}
	if err != nil {
		global.ZapLog.Error("获取验证码文本失败", zap.Error(err))
		return "", "", err
	}

	// 将验证码文本存储到Redis，设置5分钟过期
	redisKey := fmt.Sprintf("captcha:%s", id)
	expireTime := time.Duration(global.Config.Captcha.Expiration) * time.Minute
	if err := global.Redis.Set(redisKey, code, expireTime).Err(); err != nil {
		global.ZapLog.Error("存储验证码到Redis失败", zap.Error(err))
		return "", "", err
	}

	return id, b64s, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id string, code string) bool {
	if id == "" || code == "" {
		return false
	}

	// 从Redis获取验证码
	redisKey := fmt.Sprintf("captcha:%s", id)
	storedCode, err := global.Redis.Get(redisKey).Result()
	if err != nil {
		global.ZapLog.Error("从Redis获取验证码失败", zap.Error(err))
		return false
	}

	// 记录验证信息用于调试
	global.ZapLog.Info("验证码验证",
		zap.String("id", id),
		zap.String("stored", storedCode),
		zap.String("input", code),
	)

	// 不区分大小写比较验证码
	if strings.EqualFold(storedCode, code) {
		// 验证成功后删除验证码
		global.Redis.Del(redisKey)
		return true
	}

	return false
}
