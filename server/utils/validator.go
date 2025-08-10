package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// TranslateValidationError 将validator验证错误转换为可读消息
func TranslateValidationError(errs validator.ValidationErrors) string {
	var errMsg []string
	for _, e := range errs {
		// 根据不同的验证标签生成对应的错误消息
		switch e.Tag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("%s不能为空", getFieldName(e.Field())))
		case "min":
			// 将%d改为%s，因为e.Param()是字符串类型
			errMsg = append(errMsg, fmt.Sprintf("%s不能少于%s个字符", getFieldName(e.Field()), e.Param()))
		case "max":
			// 同样修复max标签的问题
			errMsg = append(errMsg, fmt.Sprintf("%s不能超过%s个字符", getFieldName(e.Field()), e.Param()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("%s验证失败: %s", getFieldName(e.Field()), e.Tag()))
		}
	}
	return strings.Join(errMsg, ", ")
}

// getFieldName 转换字段名为中文（可根据实际需求扩展）
func getFieldName(field string) string {
	// 这里可以根据实际结构体字段添加更多映射
	fieldMap := map[string]string{
		"Title":       "文章标题",
		"Content":     "文章内容",
		"CategoryID":  "分类ID",
		"AuthorID":    "作者ID",
		"UserID":      "用户ID", // 添加UserID的中文映射
	}
	if name, ok := fieldMap[field]; ok {
		return name
	}
	return field
}

// ValidateStruct 通用结构体验证函数
func ValidateStruct(data interface{}) string {
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return TranslateValidationError(validationErrors)
		}
		return "参数验证失败"
	}
	return ""
}