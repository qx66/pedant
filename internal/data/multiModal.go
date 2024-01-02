package data

import (
	"context"
	"github.com/StartOpsz/pedant/internal/biz"
)

type multiModalDataSource struct {
	data *Data
}

func NewMultiModalDataSource(data *Data) biz.MultiModalRepo {
	return &multiModalDataSource{
		data: data,
	}
}

func (multiModalDataSource *multiModalDataSource) CreateMultiModal(ctx context.Context, multiModal biz.MultiModal) error {
	tx := multiModalDataSource.data.db.WithContext(ctx).Create(&multiModal)
	return tx.Error
}

func (multiModalDataSource *multiModalDataSource) ListMultiModal(ctx context.Context, userUuid string) ([]biz.MultiModal, error) {
	var multiModals []biz.MultiModal
	tx := multiModalDataSource.data.db.WithContext(ctx).
		Where("user_uuid = ?", userUuid).
		Order("create_time desc").
		Limit(20).
		Find(&multiModals)
	
	return multiModals, tx.Error
}
