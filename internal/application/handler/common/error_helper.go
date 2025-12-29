package common

import (
	"gochat/internal/pkg/code_msg"
)

// HandleError 统一处理错误，返回业务错误码
// 如果err为nil，返回0（成功）
// 如果err不为nil，返回ServerError
func HandleError(err error) (code_msg.BusinessCode, error) {
	if err != nil {
		return code_msg.ServerError, err
	}
	return 0, nil
}

// ValidateRequired 验证必填字段
func ValidateRequired(value interface{}, fieldName string) (code_msg.BusinessCode, error) {
	if value == nil {
		return code_msg.BadRequest, nil
	}
	
	switch v := value.(type) {
	case string:
		if v == "" {
			return code_msg.BadRequest, nil
		}
	case int64:
		if v == 0 {
			return code_msg.BadRequest, nil
		}
	case uint32:
		if v == 0 {
			return code_msg.BadRequest, nil
		}
	}
	
	return 0, nil
}


