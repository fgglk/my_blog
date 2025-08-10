package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 基础响应结构体
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

// 状态码常量
const (
	ERROR   = 7
	SUCCESS = 0
	NO_AUTH = 401
)

// 通用响应函数
func Result(code int, data interface{}, msg string, statusCode int, c *gin.Context) {
	c.JSON(statusCode, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// 成功响应系列
func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", http.StatusOK, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, http.StatusOK, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", http.StatusOK, c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, http.StatusOK, c)
}

// 错误响应系列
func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", http.StatusBadRequest, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, http.StatusBadRequest, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, http.StatusBadRequest, c)
}

// 认证相关响应
func NoAuth(message string, c *gin.Context) {
	Result(NO_AUTH, gin.H{"reload": true}, message, http.StatusUnauthorized, c)
}

func Forbidden(message string, c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code: ERROR,
		Data: nil,
		Msg:  message,
	})
}