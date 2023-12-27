package baiduCloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/startopsz/rule/pkg/http"
)

type AccessToken struct {
	RefreshToken  string `json:"refresh_token,omitempty"`
	ExpiresIn     int    `json:"expires_in,omitempty"`
	SessionKey    string `json:"session_key,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
	Scope         string `json:"scope,omitempty"`
	SessionSecret string `json:"session_secret,omitempty"`
}

// 获取百度千帆调用接口的 AccessToken

func GetQianFanAccessToken(apiKey, secretKey string) (AccessToken, error) {
	var accessToken AccessToken
	baseUrl := "https://aip.baidubce.com/oauth/2.0/token"
	url := fmt.Sprintf("%s?grant_type=client_credentials&client_id=%s&client_secret=%s", baseUrl, apiKey, secretKey)
	
	req := http.Req{
		Method:  http.Post,
		Url:     url,
		Timeout: 5,
	}
	
	resp, err := req.Do()
	if err != nil {
		return accessToken, err
	}
	
	if resp.StatusCode != 200 {
		return accessToken, errors.New(
			fmt.Sprintf("http status: %d, message: %s", resp.StatusCode, string(resp.Body)),
		)
	}
	
	err = json.Unmarshal(resp.Body, &accessToken)
	if err != nil {
		return accessToken, err
	}
	
	return accessToken, nil
}

type ERNIEBotTurboReq struct {
	Messages []ERNIEBotTurboMessage `json:"messages,omitempty"`
	Stream   bool                   `json:"stream,omitempty"` // 默认false
	//Temperature  float32                `json:"temperature,omitempty"`   // 默认0.95，范围 (0, 1.0]，不能为0, 建议top_p和temperature不要同时更改
	//TopP         float32                `json:"top_p,omitempty"`         // 默认0.8，取值范围 [0, 1.0], 建议top_p和temperature不要同时更改
	//PenaltyScore bool                   `json:"penalty_score,omitempty"` // 通过对已生成的token增加惩罚，减少重复生成的现象。 默认1.0，取值范围：[1.0, 2.0]
	//UserId       string                 `json:"user_id,omitempty"`
}

type ERNIEBotTurboMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ERNIEBotTurboResponse struct {
	Id               string                     `json:"id,omitempty"`                 // 本轮对话的id
	Object           string                     `json:"object,omitempty"`             // 回包类型, chat.completion：多轮对话返回
	Created          int64                      `json:"created,omitempty"`            // 时间戳
	SentenceId       int64                      `json:"sentence_id,omitempty"`        // 表示当前子句的序号。只有在流式接口模式下会返回该字段
	IsEnd            bool                       `json:"is_end,omitempty"`             // 表示当前子句是否是最后一句。只有在流式接口模式下会返回该字段
	IsTruncated      bool                       `json:"is_truncated,omitempty"`       // 当前生成的结果是否被截断
	Result           string                     `json:"result,omitempty"`             // 对话返回结果
	NeedClearHistory bool                       `json:"need_clear_history,omitempty"` // 表示用户输入是否存在安全，是否关闭当前会话，清理历史会话信息。 true：是，表示用户输入存在安全风险，建议关闭当前会话，清理历史会话信息。 false：否，表示用户输入无安全风险
	Usage            ERNIEBotTurboResponseUsage `json:"usage,omitempty"`              // token统计信息，token数 = 汉字数+单词数*1.3 （仅为估算逻辑）
	ErrorCode        int                        `json:"error_code,omitempty"`         // 错误码
	ErrorMsg         string                     `json:"error_msg,omitempty"`          // 错误描述信息，帮助理解和解决发生的错误
}

type ERNIEBotTurboResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`     // 问题tokens数
	CompletionTokens int `json:"completion_tokens,omitempty"` // 回答tokens数
	TotalTokens      int `json:"total_tokens,omitempty"`      // tokens总数
}

func SendERNIEBotTurbo(accessToken string, body ERNIEBotTurboReq) (ERNIEBotTurboResponse, error) {
	var ernieBotResp ERNIEBotTurboResponse
	baseUrl := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant"
	url := fmt.Sprintf("%s?access_token=%s", baseUrl, accessToken)
	
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return ernieBotResp, err
	}
	
	//
	req := http.Req{
		Method:  http.Post,
		Url:     url,
		Body:    bodyByte,
		Timeout: 120,
	}
	
	resp, err := req.Do()
	if err != nil {
		return ernieBotResp, err
	}
	//
	
	if resp.StatusCode != 200 {
		return ernieBotResp, errors.New(
			fmt.Sprintf("http status: %d, message: %s", resp.StatusCode, string(resp.Body)),
		)
	}
	
	//
	err = json.Unmarshal(resp.Body, &ernieBotResp)
	if err != nil {
		return ernieBotResp, err
	}
	
	if ernieBotResp.ErrorCode != 0 {
		return ernieBotResp, errors.New(ernieBotResp.ErrorMsg)
	}
	
	//
	return ernieBotResp, nil
}
