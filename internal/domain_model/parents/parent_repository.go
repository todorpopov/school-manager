package parents

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type IParentRepository interface {
	CreateParent(ctx context.Context, tx pgx.Tx, userId int32) (*Parent, *exceptions.AppError)
	GetParentById(ctx context.Context, tx pgx.Tx, parentId int32) (*Parent, *exceptions.AppError)
	GetParentByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Parent, *exceptions.AppError)
	GetParents(ctx context.Context, tx pgx.Tx) ([]Parent, *exceptions.AppError)
	DeleteParent(ctx context.Context, tx pgx.Tx, parentId int32) *exceptions.AppError
}

type ParentRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewParentRepository(db *persistence.Database, logger *zap.Logger) *ParentRepository {
	return &ParentRepository{db, logger}
}

func (pr *ParentRepository) CreateParent(ctx context.Context, tx pgx.Tx, userId int32) (*Parent, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	sql := `
		INSERT INTO parents (user_id)
		VALUES ($1)
		RETURNING parent_id, user_id;
	`

	var parent Parent
	var err error

	if tx != nil {
		pr.logger.Debug("Creating parent in transaction")
		err = tx.QueryRow(ctx, sql, userId).
			Scan(&parent.ParentId, &parent.UserId)
	} else {
		pr.logger.Debug("Creating parent without transaction")
		err = pr.db.Pool.QueryRow(ctx, sql, userId).
			Scan(&parent.ParentId, &parent.UserId)
	}

	if err != nil {
		pr.logger.Error("Failed to create parent", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &parent, nil
}

func (pr *ParentRepository) GetParentById(ctx context.Context, tx pgx.Tx, parentId int32) (*Parent, *exceptions.AppError) {
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	sql := `
		SELECT p.parent_id, p.user_id, u.first_name, u.last_name, u.email
		FROM parents p
		INNER JOIN users u ON p.user_id = u.user_id
		WHERE p.parent_id = $1;
	`

	var parent Parent
	var err error

	if tx != nil {
		pr.logger.Debug("Getting parent by id in transaction", zap.Int32("parent_id", parentId))
		err = tx.QueryRow(ctx, sql, parentId).
			Scan(&parent.ParentId, &parent.UserId, &parent.FirstName, &parent.LastName, &parent.Email)
	} else {
		pr.logger.Debug("Getting parent by id without transaction", zap.Int32("parent_id", parentId))
		err = pr.db.Pool.QueryRow(ctx, sql, parentId).
			Scan(&parent.ParentId, &parent.UserId, &parent.FirstName, &parent.LastName, &parent.Email)
	}

	if err != nil {
		pr.logger.Error("Failed to get parent by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &parent, nil
}

func (pr *ParentRepository) GetParentByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Parent, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	sql := `
		SELECT p.parent_id, p.user_id, u.first_name, u.last_name, u.email
		FROM parents p
		INNER JOIN users u ON p.user_id = u.user_id
		WHERE p.user_id = $1;
	`

	var parent Parent
	var err error

	if tx != nil {
		pr.logger.Debug("Getting parent by user_id in transaction", zap.Int32("user_id", userId))
		err = tx.QueryRow(ctx, sql, userId).
			Scan(&parent.ParentId, &parent.UserId, &parent.FirstName, &parent.LastName, &parent.Email)
	} else {
		pr.logger.Debug("Getting parent by user_id without transaction", zap.Int32("user_id", userId))
		err = pr.db.Pool.QueryRow(ctx, sql, userId).
			Scan(&parent.ParentId, &parent.UserId, &parent.FirstName, &parent.LastName, &parent.Email)
	}

	if err != nil {
		pr.logger.Error("Failed to get parent by user_id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &parent, nil
}

func (pr *ParentRepository) GetParents(ctx context.Context, tx pgx.Tx) ([]Parent, *exceptions.AppError) {
	sql := `
		SELECT p.parent_id, p.user_id, u.first_name, u.last_name, u.email
		FROM parents p
		INNER JOIN users u ON p.user_id = u.user_id;
	`

	var parents []Parent
	var err error
	var rows pgx.Rows

	if tx != nil {
		pr.logger.Debug("Getting parents in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		pr.logger.Debug("Getting parents without transaction")
		rows, err = pr.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		pr.logger.Error("Failed to get parents", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var parent Parent
		err = rows.Scan(&parent.ParentId, &parent.UserId, &parent.FirstName, &parent.LastName, &parent.Email)
		if err != nil {
			pr.logger.Error("Failed to scan parent row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		parents = append(parents, parent)
	}

	if err = rows.Err(); err != nil {
		pr.logger.Error("Error iterating parents rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return parents, nil
}


func (pr *ParentRepository) DeleteParent(ctx context.Context, tx pgx.Tx, parentId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	sql := "DELETE FROM parents WHERE parent_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		pr.logger.Debug("Deleting parent in transaction")
		cmdTag, err = tx.Exec(ctx, sql, parentId)
	} else {
		pr.logger.Debug("Deleting parent without transaction")
		cmdTag, err = pr.db.Pool.Exec(ctx, sql, parentId)
	}
	if err != nil {
		pr.logger.Error("Failed to delete parent", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		pr.logger.Error("Failed to delete parent - parent not found", zap.Int32("parent_id", parentId))
		return exceptions.NewNotFoundError("Parent not found")
	}

	return nil
}

