package biz

import (
	"context"
	"encoding/json"
	"github.com/qx66/pedant/internal/biz/common"
	"github.com/qx66/pedant/internal/conf"
	"github.com/qx66/pedant/pkg/baiduCloud"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/startopsz/rule/pkg/response/errCode"
	"go.uber.org/zap"
	"time"
)

type Image struct {
	Uuid           string `json:"uuid,omitempty"`
	UserUuid       string `json:"userUuid,omitempty"`
	Prompt         string `json:"prompt,omitempty"`
	NegativePrompt string `json:"negativePrompt,omitempty"` // 否定的prompt
	Images         []byte `json:"images,omitempty"`
	PromptTokens   int    `json:"promptTokens,omitempty"`
	TotalTokens    int    `json:"totalTokens,omitempty"`
	CreateTime     int64  `json:"createTime,omitempty"`
}

func (image Image) TableName() string {
	return "image"
}

type ImageRepo interface {
	CreateImage(ctx context.Context, image Image) error
	ListImage(ctx context.Context, userUuid string) ([]Image, error)
}

type ImageUseCase struct {
	imageRepo      ImageRepo
	localCacheRepo LocalCacheRepo
	pedant         *conf.Pedant
	llm            *conf.Llm
	logger         *zap.Logger
}

func NewImageUseCase(imageRepo ImageRepo, localCacheRepo LocalCacheRepo, pedant *conf.Pedant, llm *conf.Llm, logger *zap.Logger) *ImageUseCase {
	switch pedant.ImageLlm {
	case OpenAILLM:
		if llm.Openai.ApiKey == "" {
			panic("配置使用openai大模型语言，但未配置apikey")
		}
	case GoogleLLM:
		if llm.Gemini.ApiKey == "" {
			panic("配置使用google大模型语言，但未配置apikey")
		}
	case BaiduCloudLLM:
		if llm.Qianfan.ApiKey == "" || llm.Qianfan.SecretKey == "" {
			panic("配置使用百度云大模型语言，但未配置apikey/secretKey")
		}
	default:
		panic("配置使用未知的大模型语言")
	}
	
	return &ImageUseCase{
		imageRepo:      imageRepo,
		localCacheRepo: localCacheRepo,
		pedant:         pedant,
		llm:            llm,
		logger:         logger,
	}
}

type GetImageReq struct {
	UserUuid string `json:"userUuid,omitempty" form:"userUuid"  validate:"required"`
}

func (imageUseCase *ImageUseCase) Get(c *gin.Context) {
	var req GetImageReq
	err := common.BindUriQuery(c, &req)
	if err != nil {
		return
	}
	
	images, err := imageUseCase.imageRepo.ListImage(c.Request.Context(), req.UserUuid)
	if err != nil {
		imageUseCase.logger.Error("查询数据库失败", zap.Error(err))
		c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "images": images})
}

type GenerateImageReq struct {
	UserUuid       string `json:"userUuid,omitempty" form:"userUuid"  validate:"required"`
	Prompt         string `json:"prompt,omitempty" form:"prompt" validate:"required"`
	NegativePrompt string `json:"negativePrompt" form:"negativePrompt"`
	Count          int    `json:"count,omitempty" form:"count"`
}

func (imageUseCase *ImageUseCase) Create(c *gin.Context) {
	var req GenerateImageReq
	err := common.JsonUnmarshal(c, &req)
	if err != nil {
		return
	}
	
	//
	switch imageUseCase.pedant.ImageLlm {
	case OpenAILLM:
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": "UnSupport LLM"})
		return
	
	case GoogleLLM:
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": "UnSupport LLM"})
		return
	
	case BaiduCloudLLM:
		//sessionUseCase.llm.Qianfan
		token, err := imageUseCase.GetAccessToken(c.Request.Context())
		if err != nil {
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		generateImageReq := baiduCloud.StableDiffusionXLReq{
			Prompt:         req.Prompt,
			NegativePrompt: req.NegativePrompt,
			Size:           baiduCloud.StableDiffusionSize1024x1024,
			N:              req.Count,
			SamplerIndex:   baiduCloud.StableDiffusionSamplerIndexEuler,
		}
		resp, err := baiduCloud.GenerateStableDiffusionXLImage(token, generateImageReq)
		
		if err != nil {
			imageUseCase.logger.Error("请求百度云API失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		imagesByte, err := json.Marshal(resp.Data)
		if err != nil {
			imageUseCase.logger.Error("Json序列化结果失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		err = imageUseCase.imageRepo.CreateImage(c.Request.Context(), Image{
			Uuid:           uuid.NewString(),
			UserUuid:       req.UserUuid,
			Prompt:         req.Prompt,
			NegativePrompt: req.NegativePrompt,
			Images:         imagesByte,
			PromptTokens:   resp.Usage.PromptTokens,
			TotalTokens:    resp.Usage.TotalTokens,
			CreateTime:     time.Now().Unix(),
		})
		if err != nil {
			imageUseCase.logger.Error("插入数据到数据库失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "content": resp})
		return
	
	default:
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": "UnSupport LLM"})
		return
	}
}

func (imageUseCase *ImageUseCase) GetAccessToken(ctx context.Context) (string, error) {
	accessTokenByte, err := imageUseCase.localCacheRepo.GetLocalCache(accessTokenKey)
	
	if string(accessTokenByte) != "" && err == nil {
		return string(accessTokenByte), nil
	}
	
	// 获取失败，则通过 API 重新获取Token
	if err != nil {
		imageUseCase.logger.Error("从LocalCache中获取AccessToken失败", zap.Error(err))
	}
	
	// 通过 API 获取Token
	accessToken, err := baiduCloud.GetQianFanAccessToken(imageUseCase.llm.Qianfan.ApiKey, imageUseCase.llm.Qianfan.SecretKey)
	if err != nil {
		imageUseCase.logger.Error("调用百度千帆API获取AccessToken失败", zap.Error(err))
		return "", err
	}
	
	err = imageUseCase.localCacheRepo.SetLocalCache(accessTokenKey, []byte(accessToken.AccessToken))
	if err != nil {
		imageUseCase.logger.Error("设置LocalCache的AccessToken失败", zap.Error(err))
	}
	
	return accessToken.AccessToken, nil
}

func (imageUseCase *ImageUseCase) GetAccessTokenByLocalCache(ctx context.Context) ([]byte, error) {
	return imageUseCase.localCacheRepo.GetLocalCache(accessTokenKey)
}
