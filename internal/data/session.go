package data

import (
	"context"
	"github.com/StartOpsz/pedant/internal/biz"
	"gorm.io/gorm"
)

type sessionDataSource struct {
	data *Data
}

func NewSessionDataSource(data *Data) biz.SessionRepo {
	return &sessionDataSource{
		data: data,
	}
}

func (sessionDataSource *sessionDataSource) CreateSession(ctx context.Context, session biz.Session) error {
	tx := sessionDataSource.data.db.WithContext(ctx).Create(&session)
	return tx.Error
}

func (sessionDataSource *sessionDataSource) ListSession(ctx context.Context, userUuid string) ([]biz.Session, error) {
	var sessions []biz.Session
	tx := sessionDataSource.data.db.WithContext(ctx).
		Where("user_uuid = ?", userUuid).
		Order("create_time").
		Limit(20).
		Find(&sessions)
	
	return sessions, tx.Error
}

func (sessionDataSource *sessionDataSource) DeleteSession(ctx context.Context, uuid, userUuid string) error {
	tx := sessionDataSource.data.db.WithContext(ctx).
		Where("uuid = ? and user_uuid = ?", uuid, userUuid).
		Delete(&biz.Session{})
	return tx.Error
}

func (sessionDataSource *sessionDataSource) ExistsSession(ctx context.Context, uuid, userUuid string) (bool, error) {
	tx := sessionDataSource.data.db.WithContext(ctx).
		Where("uuid = ? and user_uuid = ?", uuid, userUuid).
		First(&biz.Session{})
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, tx.Error
	}
	return true, nil
}

func (sessionDataSource *sessionDataSource) GetSessionContext(ctx context.Context, sessionUuid string) ([]biz.Context, error) {
	var contexts []biz.Context
	tx := sessionDataSource.data.db.WithContext(ctx).
		Where("session_uuid = ?", sessionUuid).
		Order("create_time").Limit(20).
		Find(&contexts)
	return contexts, tx.Error
}

func (sessionDataSource *sessionDataSource) InsertSessionContext(ctx context.Context, c biz.Context) error {
	tx := sessionDataSource.data.db.WithContext(ctx).Create(&c)
	return tx.Error
}
