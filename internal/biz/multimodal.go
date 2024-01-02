package biz

import (
	"context"
	"encoding/json"
	"github.com/StartOpsz/pedant/internal/biz/common"
	"github.com/StartOpsz/pedant/internal/conf"
	"github.com/StartOpsz/pedant/pkg/gemini"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/startopsz/rule/pkg/response/errCode"
	"go.uber.org/zap"
	"time"
)

type MultiModal struct {
	Uuid             string `json:"uuid,omitempty"`
	UserUuid         string `json:"userUuid,omitempty"`
	UserContent      string `json:"userContent,omitempty"`
	Images           []byte `json:"images,omitempty"`
	AssistantContent string `json:"assistantContent,omitempty"`
	Llm              string `json:"llm,omitempty"`
	CreateTime       int64  `json:"createTime,omitempty"`
}

func (multiModal MultiModal) TableName() string {
	return "multi_modal"
}

type MultiModalRepo interface {
	CreateMultiModal(ctx context.Context, multiModal MultiModal) error
	ListMultiModal(ctx context.Context, userUuid string) ([]MultiModal, error)
}

type MultiModalUseCase struct {
	multiModalRepo MultiModalRepo
	localCacheRepo LocalCacheRepo
	pedant         *conf.Pedant
	llm            *conf.Llm
	logger         *zap.Logger
}

func NewMultiModalUseCase(multiModalRepo MultiModalRepo, localCacheRepo LocalCacheRepo, pedant *conf.Pedant, llm *conf.Llm, logger *zap.Logger) *MultiModalUseCase {
	return &MultiModalUseCase{
		multiModalRepo: multiModalRepo,
		localCacheRepo: localCacheRepo,
		pedant:         pedant,
		llm:            llm,
		logger:         logger,
	}
}

type GetMultiModalReq struct {
	UserUuid string `json:"userUuid,omitempty" form:"userUuid"  validate:"required"`
}

func (multiModalUseCase *MultiModalUseCase) Get(c *gin.Context) {
	var req GetMultiModalReq
	err := common.BindUriQuery(c, &req)
	if err != nil {
		return
	}
	
	//
	contents, err := multiModalUseCase.multiModalRepo.ListMultiModal(c.Request.Context(), req.UserUuid)
	if err != nil {
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "content": contents})
}

type CreateMultiModalReq struct {
	UserUuid string   `json:"userUuid,omitempty" form:"userUuid"  validate:"required"`
	Content  string   `json:"content,omitempty" form:"content" validate:"required"`
	Images   []string `json:"images,omitempty" form:"images" validate:"required"`
}

func (multiModalUseCase *MultiModalUseCase) Create(c *gin.Context) {
	var req CreateMultiModalReq
	err := common.JsonUnmarshal(c, &req)
	if err != nil {
		return
	}
	
	//
	imageByte, err := json.Marshal(req.Images)
	if err != nil {
		multiModalUseCase.logger.Error("Json序列化Images失败", zap.Error(err))
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	
	//
	if len(req.Images) > 6 {
		c.JSON(500, gin.H{"errCode": errCode.ParameterFormatErrCode, "errMsg": errCode.ParameterFormatErrMsg})
		return
	}
	
	//
	switch multiModalUseCase.pedant.Llm {
	case OpenAILLM:
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": "UnSupport LLM"})
		return
	
	case GoogleLLM:
		apiKey := multiModalUseCase.llm.Gemini.ApiKey
		k := gemini.ApiKey(apiKey)
		
		resp, err := k.MultiModal(req.Content, req.Images...)
		if err != nil {
			multiModalUseCase.logger.Error("请求Google Gemini Api失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		err = multiModalUseCase.multiModalRepo.CreateMultiModal(c.Request.Context(), MultiModal{
			Uuid:             uuid.NewString(),
			UserUuid:         req.UserUuid,
			UserContent:      req.Content,
			Images:           imageByte,
			AssistantContent: resp.Candidates[0].Content.Parts[0].Text,
			Llm:              multiModalUseCase.pedant.Llm,
			CreateTime:       time.Now().Unix(),
		})
		if err != nil {
			multiModalUseCase.logger.Error("插入数据库失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		//c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "resp": resp})
		c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "content": resp.Candidates[0].Content.Parts[0].Text})
		return
	
	case BaiduCloudLLM:
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": "UnSupport LLM"})
		return
	
	default:
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": "UnSupport LLM"})
		return
	}
	
}
