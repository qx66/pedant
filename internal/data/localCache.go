package data

import "github.com/StartOpsz/pedant/internal/biz"

type localCacheDataSource struct {
	data *Data
}

func NewLocalCacheDataSource(data *Data) biz.LocalCacheRepo {
	return &localCacheDataSource{
		data: data,
	}
}

func (localCacheDataSource *localCacheDataSource) SetLocalCache(key string, value []byte) error {
	return localCacheDataSource.data.localCache.Set(key, value)
}

func (localCacheDataSource *localCacheDataSource) GetLocalCache(key string) ([]byte, error) {
	return localCacheDataSource.data.localCache.Get(key)
}
