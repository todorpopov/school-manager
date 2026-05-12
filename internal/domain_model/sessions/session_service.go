package sessions

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type ISessionService interface {
	CreateOrRenewSession(ctx context.Context, tx pgx.Tx, userId int32) (*Session, *exceptions.AppError)
	GetActiveSessionById(ctx context.Context, tx pgx.Tx, sessionId string) (*Session, *exceptions.AppError)
	SetSessionRole(ctx context.Context, tx pgx.Tx, sessionId string, role string) (*Session, *exceptions.AppError)
}

type SessionService struct {
	sessionRepo ISessionRepository
}

func NewSessionService(sessionRepo ISessionRepository) *SessionService {
	return &SessionService{sessionRepo}
}

func (ss *SessionService) CreateOrRenewSession(ctx context.Context, tx pgx.Tx, userId int32) (*Session, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	return ss.sessionRepo.CreateOrRenewSession(ctx, tx, userId)
}

func (ss *SessionService) GetActiveSessionById(ctx context.Context, tx pgx.Tx, sessionId string) (*Session, *exceptions.AppError) {
	return ss.sessionRepo.GetActiveSessionById(ctx, tx, sessionId)
}

func (ss *SessionService) SetSessionRole(ctx context.Context, tx pgx.Tx, sessionId string, role string) (*Session, *exceptions.AppError) {
	if msg := domain_model.ValidateRoleName(role); msg != "" {
		return nil, exceptions.NewValidationError("Invalid role", map[string]string{"role": msg})
	}

	return ss.sessionRepo.SetSessionRole(ctx, tx, sessionId, role)
}
