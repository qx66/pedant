package common

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/startopsz/rule/pkg/response/errCode"
	"google.golang.org/grpc/metadata"
)

// 统一对请求请求参数为 application/json 类型的数据进行 Unmarshal

var validate *validator.Validate

func JsonUnmarshal[r any](c *gin.Context, req r) error {
	rawDataByte, err := c.GetRawData()
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(400, gin.H{"errCode": errCode.ParameterFormatErrCode, "errMsg": errCode.ParameterFormatErrMsg})
		c.Abort()
		return err
	}
	
	err = json.Unmarshal(rawDataByte, req)
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(400, gin.H{"errCode": errCode.ParameterFormatErrCode, "errMsg": errCode.ParameterFormatErrMsg})
		c.Abort()
		return err
	}
	
	validate = validator.New()
	err = validate.Struct(req)
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(400, gin.H{"errCode": errCode.ParameterFormatErrCode, "errMsg": errCode.ParameterFormatErrMsg})
		c.Abort()
		return err
	}
	
	return nil
}

// 将 gin URI 中的参数绑定到 r any 中

func BindUriQuery[r any](c *gin.Context, req r) error {
	err := c.BindQuery(req)
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(400, gin.H{"errCode": errCode.ParameterFormatErrCode, "errMsg": errCode.ParameterFormatErrMsg})
		c.Abort()
		return err
	}
	
	validate = validator.New()
	err = validate.Struct(req)
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(400, gin.H{"errCode": errCode.ParameterFormatErrCode, "errMsg": errCode.ParameterFormatErrMsg})
		c.Abort()
		return err
	}
	return nil
}

// 处理 GRPC 响应 Error 的请求

func ResponseGrpcError(c *gin.Context, err error) {
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(500, gin.H{"errCode": errCode.GRpcCallErrorCode, "errMsg": errCode.GRpcCallErrorMsg})
		c.Abort()
		return
	}
}

// 继承 gin.Context 返回 grpc Context

func GetGrpcCtx(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	md := metadata.Pairs()
	return metadata.NewOutgoingContext(ctx, md)
}
