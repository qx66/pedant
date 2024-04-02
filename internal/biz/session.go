package biz

import (
	"context"
	"github.com/qx66/pedant/internal/biz/common"
	"github.com/qx66/pedant/internal/conf"
	"github.com/qx66/pedant/pkg/baiduCloud"
	"github.com/qx66/pedant/pkg/gemini"
	"github.com/qx66/pedant/pkg/openai"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/startopsz/rule/pkg/response/errCode"
	"go.uber.org/zap"
	"time"
)

type Session struct {
	Uuid       string `json:"uuid,omitempty"`
	UserUuid   string `json:"userUuid,omitempty"`
	Name       string `json:"name,omitempty"` // Subject
	CreateTime int64  `json:"createTime,omitempty"`
}

func (session Session) TableName() string {
	return "session"
}

type Context struct {
	Uuid             string `json:"uuid,omitempty"`
	SessionUuid      string `json:"sessionUuid,omitempty"`
	UserContent      string `json:"userContent,omitempty"`
	AssistantContent string `json:"assistantContent,omitempty"`
	PromptTokens     int    `json:"promptTokens,omitempty"`
	CompletionTokens int    `json:"completionTokens,omitempty"`
	TotalTokens      int    `json:"totalTokens,omitempty"`
	Llm              string `json:"llm,omitempty"`
	CreateTime       int64  `json:"createTime,omitempty"`
}

func (context Context) TableName() string {
	return "session_context"
}

type SessionRepo interface {
	CreateSession(ctx context.Context, session Session) error
	ListSession(ctx context.Context, userUuid string) ([]Session, error)
	DeleteSession(ctx context.Context, uuid, userUuid string) error
	ExistsSession(ctx context.Context, uuid, userUuid string) (bool, error)
	GetSessionContext(ctx context.Context, sessionUuid string) ([]Context, error)
	InsertSessionContext(ctx context.Context, c Context) error
}

type SessionUseCase struct {
	sessionRepo    SessionRepo
	localCacheRepo LocalCacheRepo
	pedant         *conf.Pedant
	llm            *conf.Llm
	logger         *zap.Logger
}

func NewSessionUseCase(sessionRepo SessionRepo, localCacheRepo LocalCacheRepo, pedant *conf.Pedant, llm *conf.Llm, logger *zap.Logger) *SessionUseCase {
	switch pedant.Llm {
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
	
	return &SessionUseCase{
		sessionRepo:    sessionRepo,
		localCacheRepo: localCacheRepo,
		llm:            llm,
		pedant:         pedant,
		logger:         logger,
	}
}

type ListSessionReq struct {
	UserUuid string `json:"userUuid,omitempty" form:"userUuid"  validate:"required"`
}

func (sessionUseCase *SessionUseCase) ListSession(c *gin.Context) {
	req := ListSessionReq{}
	err := common.BindUriQuery(c, &req)
	if err != nil {
		return
	}
	
	//
	sessions, err := sessionUseCase.sessionRepo.ListSession(c.Request.Context(), req.UserUuid)
	if err != nil {
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "sessions": sessions})
}

type CreateSessionReq struct {
	UserUuid string `json:"userUuid,omitempty"  validate:"required"`
	Name     string `json:"name,omitempty"  validate:"required"`
}

func (sessionUseCase *SessionUseCase) CreateSession(c *gin.Context) {
	req := CreateSessionReq{}
	err := common.JsonUnmarshal(c, &req)
	if err != nil {
		return
	}
	
	sessionUuid := uuid.NewString()
	//
	err = sessionUseCase.sessionRepo.CreateSession(c.Request.Context(), Session{
		Uuid:       sessionUuid,
		UserUuid:   req.UserUuid,
		Name:       req.Name,
		CreateTime: time.Now().Unix(),
	})
	
	if err != nil {
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "uuid": sessionUuid})
}

type DelSessionReq struct {
	UserUuid string `json:"userUuid,omitempty" form:"userUuid"  validate:"required"`
	Uuid     string `json:"uuid" form:"uuid"  validate:"required"`
}

func (sessionUseCase *SessionUseCase) DelSession(c *gin.Context) {
	req := DelSessionReq{}
	err := common.BindUriQuery(c, &req)
	if err != nil {
		return
	}
	
	err = sessionUseCase.sessionRepo.DeleteSession(c.Request.Context(), req.Uuid, req.UserUuid)
	if err != nil {
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg})
}

type ListSessionContextReq struct {
	UserUuid    string `json:"userUuid,omitempty" form:"userUuid"  validate:"required"`
	SessionUuid string `json:"sessionUuid,omitempty" form:"sessionUuid" validate:"required"`
}

func (sessionUseCase *SessionUseCase) ListSessionContext(c *gin.Context) {
	req := ListSessionContextReq{}
	err := common.BindUriQuery(c, &req)
	if err != nil {
		return
	}
	//
	e, err := sessionUseCase.sessionRepo.ExistsSession(c.Request.Context(), req.SessionUuid, req.UserUuid)
	if err != nil {
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	
	if !e {
		c.JSON(404, gin.H{"errCode": errCode.NotFoundCode, "errMsg": errCode.NotFoundMsg})
		return
	}
	
	//
	contexts, err := sessionUseCase.sessionRepo.GetSessionContext(c.Request.Context(), req.SessionUuid)
	if err != nil {
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "contexts": contexts})
	return
}

type CreateSessionContextReq struct {
	UserUuid    string `json:"userUuid,omitempty" form:"userUuid"  validate:"required"`
	SessionUuid string `json:"sessionUuid,omitempty" form:"sessionUuid" validate:"required"`
	Content     string `json:"content,omitempty" form:"content" validate:"required"`
}

func (sessionUseCase *SessionUseCase) CreateSessionContext(c *gin.Context) {
	req := CreateSessionContextReq{}
	err := common.JsonUnmarshal(c, &req)
	if err != nil {
		return
	}
	
	//
	e, err := sessionUseCase.sessionRepo.ExistsSession(c.Request.Context(), req.SessionUuid, req.UserUuid)
	if err != nil {
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	
	if !e {
		c.JSON(404, gin.H{"errCode": errCode.NotFoundCode, "errMsg": errCode.NotFoundMsg})
		return
	}
	
	//
	contexts, err := sessionUseCase.sessionRepo.GetSessionContext(c.Request.Context(), req.SessionUuid)
	if err != nil {
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
		return
	}
	
	//
	switch sessionUseCase.pedant.Llm {
	case OpenAILLM:
		apiKey := sessionUseCase.llm.Openai.ApiKey
		body := generateOpenAiContext(contexts)
		resp, err := openai.SendChat(body, apiKey)
		if err != nil {
			sessionUseCase.logger.Error("请求OpenAi API失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		err = sessionUseCase.sessionRepo.InsertSessionContext(c.Request.Context(), Context{
			Uuid:             uuid.NewString(),
			SessionUuid:      req.SessionUuid,
			UserContent:      req.Content,
			AssistantContent: resp.Choices[0].Message.Content,
			PromptTokens:     int(resp.Usage.PromptTokens),
			CompletionTokens: int(resp.Usage.CompletionTokens),
			TotalTokens:      int(resp.Usage.TotalTokens),
			Llm:              OpenAILLM,
			CreateTime:       time.Now().Unix(),
		})
		if err != nil {
			sessionUseCase.logger.Error("插入数据库失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "content": resp.Choices[0].Message.Content})
		return
	
	case GoogleLLM:
		apiKey := sessionUseCase.llm.Gemini.ApiKey
		k := gemini.ApiKey(apiKey)
		
		his := generateGeminiContext(contexts)
		resp, err := k.Chat(his, req.Content)
		if err != nil {
			sessionUseCase.logger.Error("请求Google Gemini Api失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		err = sessionUseCase.sessionRepo.InsertSessionContext(c.Request.Context(), Context{
			Uuid:             uuid.NewString(),
			SessionUuid:      req.SessionUuid,
			UserContent:      req.Content,
			AssistantContent: resp.Candidates[0].Content.Parts[0].Text,
			//PromptTokens:     resp.Usage.PromptTokens,
			//CompletionTokens: resp.Usage.CompletionTokens,
			//TotalTokens:      resp.Usage.TotalTokens,
			Llm:        GoogleLLM,
			CreateTime: time.Now().Unix(),
		})
		if err != nil {
			sessionUseCase.logger.Error("插入数据库失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "content": resp.Candidates[0].Content.Parts[0].Text})
		return
	
	case BaiduCloudLLM:
		//sessionUseCase.llm.Qianfan
		token, err := sessionUseCase.GetAccessToken(c.Request.Context())
		if err != nil {
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		body := generateQianFanContext(contexts, req.Content)
		
		resp, err := baiduCloud.SendERNIEBotTurbo(token, body)
		if err != nil {
			sessionUseCase.logger.Error("请求百度云API失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		err = sessionUseCase.sessionRepo.InsertSessionContext(c.Request.Context(), Context{
			Uuid:             uuid.NewString(),
			SessionUuid:      req.SessionUuid,
			UserContent:      req.Content,
			AssistantContent: resp.Result,
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
			Llm:              BaiduCloudLLM,
			CreateTime:       time.Now().Unix(),
		})
		if err != nil {
			sessionUseCase.logger.Error("插入数据库失败", zap.Error(err))
			c.JSON(200, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": errCode.BizOpErrorMsg})
			return
		}
		
		c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "content": resp.Result})
		return
	
	default:
		c.JSON(500, gin.H{"errCode": errCode.BizOpErrorCode, "errMsg": "UnSupport LLM"})
		return
	}
}

func (sessionUseCase *SessionUseCase) GetAccessToken(ctx context.Context) (string, error) {
	accessTokenByte, err := sessionUseCase.localCacheRepo.GetLocalCache(accessTokenKey)
	
	if string(accessTokenByte) != "" && err == nil {
		return string(accessTokenByte), nil
	}
	
	// 获取失败，则通过 API 重新获取Token
	if err != nil {
		sessionUseCase.logger.Error("从LocalCache中获取AccessToken失败", zap.Error(err))
	}
	
	// 通过 API 获取Token
	accessToken, err := baiduCloud.GetQianFanAccessToken(sessionUseCase.llm.Qianfan.ApiKey, sessionUseCase.llm.Qianfan.SecretKey)
	if err != nil {
		sessionUseCase.logger.Error("调用百度千帆API获取AccessToken失败", zap.Error(err))
		return "", err
	}
	
	err = sessionUseCase.localCacheRepo.SetLocalCache(accessTokenKey, []byte(accessToken.AccessToken))
	if err != nil {
		sessionUseCase.logger.Error("设置LocalCache的AccessToken失败", zap.Error(err))
	}
	
	return accessToken.AccessToken, nil
}

func (sessionUseCase *SessionUseCase) GetAccessTokenByLocalCache(ctx context.Context) ([]byte, error) {
	return sessionUseCase.localCacheRepo.GetLocalCache(accessTokenKey)
}

func generateOpenAiContext(contexts []Context) openai.GptTurbo0301 {
	var messages []openai.GptTurbo0301Message
	messages = append(messages, openai.GptTurbo0301Message{
		Role:    openai.ChatRoleSystem,
		Content: "你是一个聪明的小助理",
	})
	
	for _, c := range contexts {
		if c.Llm == OpenAILLM {
			messages = append(messages, openai.GptTurbo0301Message{
				Role:    openai.ChatRoleUser,
				Content: c.UserContent,
			})
			
			messages = append(messages, openai.GptTurbo0301Message{
				Role:    openai.ChatRoleAssistant,
				Content: c.AssistantContent,
			})
		}
	}
	
	return openai.GptTurbo0301{
		Model:    openai.ChatModuleGpt35Turbo,
		Messages: messages,
	}
}

func generateGeminiContext(contexts []Context) []gemini.Content {
	var cs []gemini.Content
	
	for _, c := range contexts {
		if c.Llm == GoogleLLM {
			cs = append(cs, gemini.Content{
				Role: gemini.ChatRoleUser,
				Parts: []interface{}{
					gemini.ContentText{
						Text: c.UserContent,
					},
				},
			})
			
			cs = append(cs, gemini.Content{
				Role: gemini.ChatRoleModel,
				Parts: []interface{}{
					gemini.ContentText{
						Text: c.AssistantContent,
					},
				},
			})
		}
	}
	
	return cs
}

func generateQianFanContext(contexts []Context, text string) baiduCloud.ERNIEBotTurboReq {
	var messages []baiduCloud.ERNIEBotTurboMessage
	
	for _, c := range contexts {
		if c.Llm == BaiduCloudLLM {
			messages = append(messages, baiduCloud.ERNIEBotTurboMessage{
				Role:    baiduCloud.ChatRoleUser,
				Content: c.UserContent,
			})
			
			messages = append(messages, baiduCloud.ERNIEBotTurboMessage{
				Role:    baiduCloud.ChatRoleAssistant,
				Content: c.AssistantContent,
			})
		}
	}
	
	messages = append(messages, baiduCloud.ERNIEBotTurboMessage{
		Role:    baiduCloud.ChatRoleUser,
		Content: text,
	})
	
	return baiduCloud.ERNIEBotTurboReq{
		Messages: messages,
		Stream:   false,
	}
}
