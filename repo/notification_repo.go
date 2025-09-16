package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"notification-service/entity"
)

type NotificationRepo struct {
	db *sqlx.DB
}

func NewNotificationRepo(db *sqlx.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

type Page struct {
	Offset int
	Limit  int
}

func (r *NotificationRepo) FindAllByPage(ctx context.Context, page Page) ([]model.Notification, error) {
	const query = `
		SELECT id, created_at, modified_at, expiration_date,
		       message, error, user_uid, message_type,
		       link, status, subject, created_by
		FROM   notifications
		ORDER  BY created_at DESC
		LIMIT  $1 OFFSET $2`

	var nn []model.Notification
	if err := r.db.SelectContext(ctx, &nn, query, page.Limit, page.Offset); err != nil {
		return nil, fmt.Errorf("select notifications: %w", err)
	}
	return nn, nil
}

func (r *NotificationRepo) Save(ctx context.Context, n *model.Notification) error {
	const stmt = `
		INSERT INTO notifications(id, created_at, modified_at, expiration_date,
		                          message, error, user_uid, message_type,
		                          link, status, subject, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, stmt,
		n.ID, n.CreatedAt, n.ModifiedAt, n.ExpirationDate,
		n.Message, n.Error, n.UserUID, n.MessageType,
		n.Link, n.Status, n.Subject, n.CreatedBy)
	return err
}