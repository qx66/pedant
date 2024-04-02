package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/qx66/pedant/internal/biz"
	"github.com/qx66/pedant/internal/conf"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "configPath", "", "-configPath")
}

type app struct {
	sessionUseCase    *biz.SessionUseCase
	multiModalUseCase *biz.MultiModalUseCase
	imageUseCase      *biz.ImageUseCase
}

func newApp(sessionUseCase *biz.SessionUseCase, multiModalUseCase *biz.MultiModalUseCase, imageUseCase *biz.ImageUseCase) *app {
	return &app{
		sessionUseCase:    sessionUseCase,
		multiModalUseCase: multiModalUseCase,
		imageUseCase:      imageUseCase,
	}
}

func main() {
	flag.Parse()
	
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("初始化日志失败")
		return
	}
	
	//
	//
	if configPath == "" {
		logger.Error("configPath 参数为空")
		return
	}
	
	//
	f, err := os.Open(configPath)
	defer f.Close()
	if err != nil {
		logger.Error(
			"加载配置文件失败",
			zap.String("configPath", configPath),
			zap.Error(err),
		)
		return
	}
	
	//
	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		logger.Error(
			"加载配置文件copy内容失败",
			zap.Error(err),
		)
		return
	}
	
	//
	var bootstrap conf.Bootstrap
	err = yaml.Unmarshal(buf.Bytes(), &bootstrap)
	if err != nil {
		logger.Error(
			"序列化配置失败",
			zap.Error(err),
		)
		return
	}
	
	iApp, clean, err := initApp(bootstrap.Data, bootstrap.Pedant, bootstrap.Llm, logger)
	defer clean()
	if err != nil {
		logger.Error("初始化程序失败", zap.Error(err))
	}
	
	route := gin.New()
	
	// session
	route.GET("/chat/session", iApp.sessionUseCase.ListSession)
	route.POST("/chat/session", iApp.sessionUseCase.CreateSession)
	route.DELETE("/chat/session", iApp.sessionUseCase.DelSession)
	
	// session context
	route.GET("/chat/session/context", iApp.sessionUseCase.ListSessionContext)
	route.POST("/chat/session/context", iApp.sessionUseCase.CreateSessionContext)
	
	//
	route.GET("/image", iApp.imageUseCase.Get)
	route.POST("/image", iApp.imageUseCase.Create)
	
	//
	route.GET("/multiModal", iApp.multiModalUseCase.Get)
	route.POST("/multiModal", iApp.multiModalUseCase.Create)
	
	err = route.Run(":20000")
	logger.Error("启动程序失败", zap.Error(err))
}
