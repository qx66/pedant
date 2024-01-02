package biz

import "github.com/google/wire"

type LocalCacheRepo interface {
	SetLocalCache(key string, value []byte) error
	GetLocalCache(key string) ([]byte, error)
}

var ProviderSet = wire.NewSet(NewSessionUseCase, NewMultiModalUseCase, NewImageUseCase)

type LLM string

const (
	GoogleLLM     = "gemini"
	OpenAILLM     = "openai"
	BaiduCloudLLM = "ernieBot"
)

const (
	accessTokenKey = "baiduQianFanAccessToken"
)
