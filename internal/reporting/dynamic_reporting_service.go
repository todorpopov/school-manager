package reporting

import (
	"context"
	"fmt"
	"strings"

	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type IDynamicReportingService interface {
	GenerateReport(ctx context.Context, reportQuery *ReportQueryRequest) (DynamicReportingResponse, *exceptions.AppError)
}

type DynamicReportingService struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewDynamicReportingService(db *persistence.Database, logger *zap.Logger) *DynamicReportingService {
	return &DynamicReportingService{db: db, logger: logger}
}

func (drs *DynamicReportingService) validateReportQuery(ctx context.Context, reportQuery *ReportQueryRequest) *exceptions.AppError {
	if err := ValidateReportQueryRequest(reportQuery); err != nil {
		return err
	}

	if len(reportQuery.TeacherIds) == 0 && len(reportQuery.SubjectIds) == 0 {
		return nil
	}

	return drs.validateFilters(ctx, reportQuery.SchoolIds, reportQuery.TeacherIds, reportQuery.SubjectIds)
}

func (drs *DynamicReportingService) validateFilters(
	ctx context.Context,
	schoolIds []int32,
	teacherIds []int32,
	subjectIds []int32,
) *exceptions.AppError {
	var sql string
	var args []interface{}

	if len(teacherIds) > 0 && len(subjectIds) > 0 {
		// Both teachers and subjects provided
		if len(schoolIds) > 0 {
			sql = `
				SELECT DISTINCT cl.school_id
				FROM curricula c
				INNER JOIN classes cl ON c.class_id = cl.class_id
				WHERE cl.school_id = ANY($1)
					AND c.teacher_id = ANY($2)
					AND c.subject_id = ANY($3)
			`
			args = []interface{}{schoolIds, teacherIds, subjectIds}
		} else {
			// No school filter - validate all schools
			sql = `
				SELECT DISTINCT cl.school_id
				FROM curricula c
				INNER JOIN classes cl ON c.class_id = cl.class_id
				WHERE c.teacher_id = ANY($1)
					AND c.subject_id = ANY($2)
			`
			args = []interface{}{teacherIds, subjectIds}
		}
	} else if len(teacherIds) > 0 {
		// Only teachers provided
		if len(schoolIds) > 0 {
			sql = `
				SELECT DISTINCT cl.school_id
				FROM curricula c
				INNER JOIN classes cl ON c.class_id = cl.class_id
				WHERE cl.school_id = ANY($1)
					AND c.teacher_id = ANY($2)
			`
			args = []interface{}{schoolIds, teacherIds}
		} else {
			// No school filter - validate all schools
			sql = `
				SELECT DISTINCT cl.school_id
				FROM curricula c
				INNER JOIN classes cl ON c.class_id = cl.class_id
				WHERE c.teacher_id = ANY($1)
			`
			args = []interface{}{teacherIds}
		}
	} else {
		// Only subjects provided
		if len(schoolIds) > 0 {
			sql = `
				SELECT DISTINCT cl.school_id
				FROM curricula c
				INNER JOIN classes cl ON c.class_id = cl.class_id
				WHERE cl.school_id = ANY($1)
					AND c.subject_id = ANY($2)
			`
			args = []interface{}{schoolIds, subjectIds}
		} else {
			// No school filter - validate all schools
			sql = `
				SELECT DISTINCT cl.school_id
				FROM curricula c
				INNER JOIN classes cl ON c.class_id = cl.class_id
				WHERE c.subject_id = ANY($1)
			`
			args = []interface{}{subjectIds}
		}
	}

	rows, err := drs.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		drs.logger.Error("Failed to validate school filters", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	validSchoolIds := make(map[int32]bool)
	for rows.Next() {
		var schoolId int32
		if err := rows.Scan(&schoolId); err != nil {
			drs.logger.Error("Failed to scan school_id", zap.Error(err))
			return exceptions.PgErrorToAppError(err)
		}
		validSchoolIds[schoolId] = true
	}

	if err = rows.Err(); err != nil {
		drs.logger.Error("Error iterating school validation results", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if len(schoolIds) == 0 {
		if len(validSchoolIds) == 0 {
			messages := map[string]string{}
			if len(teacherIds) > 0 && len(subjectIds) > 0 {
				messages["validation_error"] = "No schools have curricula matching the provided teacher AND subject filters"
			} else if len(teacherIds) > 0 {
				messages["validation_error"] = "No schools have curricula with the provided teachers"
			} else {
				messages["validation_error"] = "No schools have curricula with the provided subjects"
			}
			return exceptions.NewValidationError("Invalid filters", messages)
		}
		return nil
	}

	var invalidSchoolIds []int32
	for _, schoolId := range schoolIds {
		if !validSchoolIds[schoolId] {
			invalidSchoolIds = append(invalidSchoolIds, schoolId)
		}
	}

	if len(invalidSchoolIds) > 0 {
		messages := map[string]string{}

		if len(teacherIds) > 0 && len(subjectIds) > 0 {
			messages["validation_error"] = "Some schools have no curricula matching the provided teacher AND subject filters"
		} else if len(teacherIds) > 0 {
			messages["validation_error"] = "Some schools have no curricula with the provided teachers"
		} else {
			messages["validation_error"] = "Some schools have no curricula with the provided subjects"
		}

		invalidSchoolIdsStr := make([]string, len(invalidSchoolIds))
		for i, id := range invalidSchoolIds {
			invalidSchoolIdsStr[i] = fmt.Sprintf("%d", id)
		}
		messages["invalid_school_ids"] = strings.Join(invalidSchoolIdsStr, ", ")

		return exceptions.NewValidationError(
			"Invalid filters for one or more schools",
			messages,
		)
	}

	return nil
}

func (drs *DynamicReportingService) GenerateReport(ctx context.Context, reportQuery *ReportQueryRequest) (DynamicReportingResponse, *exceptions.AppError) {
	if err := drs.validateReportQuery(ctx, reportQuery); err != nil {
		return nil, err
	}

	if reportQuery.ReportType == ReportTypeAbsences {
		return drs.queryForAbsences(ctx, reportQuery)
	}

	return drs.queryForGrades(ctx, reportQuery)
}

func (drs *DynamicReportingService) queryForAbsences(ctx context.Context, reportQuery *ReportQueryRequest) (*AbsencesReport, *exceptions.AppError) {
	return nil, nil
}

func (drs *DynamicReportingService) queryForGrades(ctx context.Context, reportQuery *ReportQueryRequest) (*GradesReport, *exceptions.AppError) {
	return nil, nil
}
