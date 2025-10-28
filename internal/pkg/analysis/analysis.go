package analysis

import (
	"errors"
	"fmt"
	"gochat/internal/pkg/response"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
)

type Validator interface {
	Validate() error
}

// BindParameter -
func BindParameter[k any](c *gin.Context, response *response.SvcRequest) (res k, err error) {
	if c.Request.Method == "GET" {
		if len(c.Request.URL.Query()) == 0 {
			return newIfPointer[k](), nil
		}
		err = c.ShouldBindQuery(&res)
		fmt.Print(err)
	} else if c.Request.Method == "POST" {
		err = c.ShouldBindJSON(&res)
	} else {
		response.JsonError(c, err, "方法不允许")
		return
	}
	if err != nil {
		response.JsonError(c, err, "参数绑定失败")
		return res, err
	}
	// 尝试调用Validate方法（如果存在）
	if validator, ok := any(&res).(Validator); ok {
		if err = validator.Validate(); err != nil {
			response.JsonError(c, err, check(err))
			return res, err
		}
	}
	return res, nil
}

// BindQuery -
func BindQuery[k any](c *gin.Context, response *response.SvcRequest) (res k, err error) {
	if len(c.Request.URL.Query()) == 0 {
		return newIfPointer[k](), nil
	}
	err = c.ShouldBindQuery(&res)
	if err != nil {
		response.JsonError(c, err, "参数绑定失败")
		return res, err
	}
	// 尝试调用Validate方法（如果存在）
	if validator, ok := any(&res).(Validator); ok {
		if err = validator.Validate(); err != nil {
			response.JsonError(c, err, check(err))
			return res, err
		}
	}
	return res, nil
}

func newIfPointer[k any]() k {
	var zero k
	if typ := reflect.TypeOf(zero); typ != nil && typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()).Interface().(k)
	}
	return zero
}

// ParameterError -
type ParameterError interface {
	Error() string
	Field() string
	Reason() string
	Cause() error
	Key() bool
	ErrorName() string
}

// check -
func check(err error) string {
	var pe ParameterError
	if errors.As(err, &pe) {
		return fmt.Sprintf("字段【%s】%s", pe.Field(), message(pe.Reason()))
	}
	return err.Error()
}

// message -
func message(msg string) (data string) {
	if re := regexp.MustCompile(`value must be greater than (\d+)`).FindStringSubmatch(msg); len(re) > 1 {
		data = fmt.Sprintf("值必须大于%s.", re[1])
	}
	return data
}
