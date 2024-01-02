package data

import (
	"context"
	"github.com/StartOpsz/pedant/internal/conf"
	"github.com/allegro/bigcache/v3"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData,
	NewLocalCacheDataSource,
	NewSessionDataSource,
	NewMultiModalDataSource,
	NewImageDataSource)

// Data .
type Data struct {
	db         *gorm.DB
	localCache *bigcache.BigCache
	logger     *zap.Logger
}

// NewData .
func NewData(c *conf.Data, logger *zap.Logger) (*Data, func(), error) {
	// mysql
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		logger.Error("打开数据库连接失败", zap.Error(err))
		panic(err)
		//return nil, nil, err
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("返回sql.DB失败", zap.Error(err))
		return nil, nil, err
	}
	
	err = sqlDB.Ping()
	if err != nil {
		logger.Error("ping db 失败", zap.Error(err))
		return nil, nil, err
	}
	
	sqlDB.SetMaxIdleConns(int(c.Database.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(c.Database.MaxOpenConns))
	
	// localCache
	localCacheConfig := bigcache.Config{
		Shards:      8,
		LifeWindow:  20 * time.Minute,
		CleanWindow: 10 * time.Minute,
		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 100 * 10 * 60,
		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 50,
		Verbose:      true,
	}
	
	localCache, err := bigcache.New(context.Background(), localCacheConfig)
	if err != nil {
		return nil, nil, err
	}
	
	//
	d := &Data{
		db:         db.Debug(),
		localCache: localCache,
		logger:     logger,
	}
	
	cleanup := func() {
		err = sqlDB.Close()
		if err != nil {
			logger.Error("关闭MySQL资源失败", zap.Error(err))
		} else {
			logger.Info("关闭MySQL资源成功")
		}
	}
	
	return d, cleanup, nil
}
