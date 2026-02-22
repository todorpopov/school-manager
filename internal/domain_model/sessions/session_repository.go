package sessions

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type ISessionRepository interface {
	CreateOrRenewSession(ctx context.Context, tx pgx.Tx, userId int32) (*Session, *exceptions.AppError)
}

type SessionRepository struct {
	sessionExpiration time.Duration
	db                *persistence.Database
	logger            *zap.Logger
}

func NewSessionRepository(db *persistence.Database, sessionExpiration time.Duration, logger *zap.Logger) *SessionRepository {
	return &SessionRepository{db: db, sessionExpiration: sessionExpiration, logger: logger}
}

func (sr *SessionRepository) CreateOrRenewSession(ctx context.Context, tx pgx.Tx, userId int32) (*Session, *exceptions.AppError) {
	sql := "INSERT INTO sessions (user_id, expires_at) " +
		"VALUES ($1, now() + $2) ON CONFLICT (user_id) DO UPDATE SET expires_at = now() + $2 " +
		"RETURNING session_id, user_id, expires_at;"

	var session Session
	var err error
	if tx != nil {
		sr.logger.Debug("Creating or getting session in transaction")
		err = tx.QueryRow(ctx, sql, userId, sr.sessionExpiration).Scan(
			&session.SessionId,
			&session.UserId,
			&session.ExpiresAt,
		)
	} else {
		sr.logger.Debug("Creating or getting session without transaction")
		err = sr.db.Pool.QueryRow(ctx, sql, userId, sr.sessionExpiration).Scan(
			&session.SessionId,
			&session.UserId,
			&session.ExpiresAt,
		)
	}

	if err != nil {
		sr.logger.Error("Failed to create or get session", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	return &session, nil
}
