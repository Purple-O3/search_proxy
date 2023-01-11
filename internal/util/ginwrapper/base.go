package ginwrapper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/iancoleman/strcase"
)

var (
	OnlySupportGetAndPostErr = errors.New("bindAndValidate method only support GET and POST request")
)

type Base struct{}

// bindAndValidate  绑定和验证请求数据
// 请求类型: POST、GET
// 绑定参数: json、form、uri
// 参数验证规则: https://github.com/go-playground/validator
func (base *Base) BindAndValidate(c *gin.Context, reqStruct interface{}) error {
	switch c.Request.Method {
	case http.MethodGet:
		if err := c.ShouldBind(reqStruct); err != nil {
			return err
		}
	case http.MethodPost:
		// 解决body数据为空字符串或空json对象异常
		dataBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		if len(dataBytes) == 0 {
			break
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(dataBytes))

		if err := c.BindJSON(reqStruct); err != nil {
			return err
		}
	default:
		return OnlySupportGetAndPostErr
	}

	// 绑定uri参数
	if err := c.ShouldBindUri(reqStruct); err != nil {
		return err
	}
	// 绑定header参数
	if err := c.ShouldBindHeader(reqStruct); err != nil {
		return err
	}

	// validator验证
	validate := validator.New()
	err := validate.Struct(reqStruct)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		// 如果有多个错误，只会返回第一个错误类型
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("%s can't pass %s rule", strcase.ToLowerCamel(err.Field()), err.Tag())
		}
		return err
	}
	return nil
}

func (base *Base) SucMsg(ctx *gin.Context, data ...interface{}) {
	res := make(map[string]interface{})
	res["code"] = 0
	res["message"] = "success"
	if len(data) > 0 {
		res["data"] = data[0]
	} else {
		res["data"] = make(map[string]interface{})
	}
	ctx.JSON(http.StatusOK, res)
}

func (base *Base) ErrMsg(ctx *gin.Context, err error, data ...interface{}) {
	res := make(map[string]interface{})
	res["code"] = -1
	res["message"] = err.Error()
	if len(data) > 0 {
		res["data"] = data[0]
	} else {
		res["data"] = make(map[string]interface{})
	}
	ctx.JSON(http.StatusOK, res)
}

func (base *Base) ProxyMsg(ctx *gin.Context, retByte []byte) {
	ctx.Data(http.StatusOK, "application/json", retByte)
}
