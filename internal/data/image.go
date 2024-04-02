package data

import (
	"context"
	"github.com/qx66/pedant/internal/biz"
)

type imageDataSource struct {
	data *Data
}

func NewImageDataSource(data *Data) biz.ImageRepo {
	return &imageDataSource{
		data: data,
	}
}

func (imageDataSource *imageDataSource) CreateImage(ctx context.Context, image biz.Image) error {
	tx := imageDataSource.data.db.WithContext(ctx).Create(&image)
	return tx.Error
}

func (imageDataSource *imageDataSource) ListImage(ctx context.Context, userUuid string) ([]biz.Image, error) {
	var images []biz.Image
	tx := imageDataSource.data.db.WithContext(ctx).
		Where("user_uuid = ?", userUuid).
		Order("create_time desc").
		Limit(20).
		Find(&images)
	return images, tx.Error
}
