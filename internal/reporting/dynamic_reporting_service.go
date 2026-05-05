package reporting

import (
	"context"
	"fmt"
	"strings"

	"github.com/todorpopov/school-manager/internal/domain_model/schools"
	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/domain_model/teachers"
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

func (drs *DynamicReportingService) queryForAbsences(ctx context.Context, reportQuery *ReportQueryRequest) (DynamicReportingResponse, *exceptions.AppError) {
	sql := `
		SELECT 
			s.school_id, s.school_name, s.school_address,
			t.teacher_id, u_t.first_name as teacher_first_name, u_t.last_name as teacher_last_name, u_t.email as teacher_email,
			subj.subject_id, subj.subject_name,
			COUNT(a.absence_id) as total_absences,
			SUM(CASE WHEN a.is_excused = true THEN 1 ELSE 0 END) as excused_absences,
			SUM(CASE WHEN a.is_excused = false THEN 1 ELSE 0 END) as non_excused_absences
		FROM absences a
		INNER JOIN curricula cur ON a.curriculum_id = cur.curriculum_id
		INNER JOIN classes cl ON cur.class_id = cl.class_id
		INNER JOIN schools s ON cl.school_id = s.school_id
		INNER JOIN teachers t ON cur.teacher_id = t.teacher_id
		INNER JOIN users u_t ON t.user_id = u_t.user_id
		INNER JOIN subjects subj ON cur.subject_id = subj.subject_id
		WHERE 1=1
	`

	var args []interface{}
	argIndex := 1

	if len(reportQuery.SchoolIds) > 0 {
		sql += fmt.Sprintf(" AND s.school_id = ANY($%d)", argIndex)
		args = append(args, reportQuery.SchoolIds)
		argIndex++
	}

	if len(reportQuery.TeacherIds) > 0 {
		sql += fmt.Sprintf(" AND t.teacher_id = ANY($%d)", argIndex)
		args = append(args, reportQuery.TeacherIds)
		argIndex++
	}

	if len(reportQuery.SubjectIds) > 0 {
		sql += fmt.Sprintf(" AND subj.subject_id = ANY($%d)", argIndex)
		args = append(args, reportQuery.SubjectIds)
		argIndex++
	}

	sql += `
		GROUP BY s.school_id, s.school_name, s.school_address, 
		         t.teacher_id, u_t.first_name, u_t.last_name, u_t.email,
		         subj.subject_id, subj.subject_name
		ORDER BY s.school_name, u_t.last_name, u_t.first_name, subj.subject_name
	`

	rows, err := drs.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		drs.logger.Error("Failed to query absences", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	schoolMap := make(map[int32]*SchoolAbsences)
	teacherMap := make(map[int32]map[int32]*TeacherSubjectsAbsences)

	for rows.Next() {
		var (
			schoolId, teacherId, subjectId                                                          int32
			schoolName, schoolAddress, teacherFirstName, teacherLastName, teacherEmail, subjectName string
			totalAbsences, excusedAbsences, nonExcusedAbsences                                      int32
		)

		err := rows.Scan(
			&schoolId, &schoolName, &schoolAddress,
			&teacherId, &teacherFirstName, &teacherLastName, &teacherEmail,
			&subjectId, &subjectName,
			&totalAbsences, &excusedAbsences, &nonExcusedAbsences,
		)
		if err != nil {
			drs.logger.Error("Failed to scan absence row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}

		if _, exists := schoolMap[schoolId]; !exists {
			schoolMap[schoolId] = &SchoolAbsences{
				School: schools.School{
					SchoolId:      schoolId,
					SchoolName:    schoolName,
					SchoolAddress: schoolAddress,
				},
				Teachers: []TeacherSubjectsAbsences{},
			}
			teacherMap[schoolId] = make(map[int32]*TeacherSubjectsAbsences)
		}

		if _, exists := teacherMap[schoolId][teacherId]; !exists {
			teacherMap[schoolId][teacherId] = &TeacherSubjectsAbsences{
				Teacher: teachers.Teacher{
					TeacherId: teacherId,
					UserId:    0, // Not needed for report
					FirstName: teacherFirstName,
					LastName:  teacherLastName,
					Email:     teacherEmail,
				},
				Subjects: []SubjectAbsenceSummary{},
			}
		}

		teacherMap[schoolId][teacherId].Subjects = append(
			teacherMap[schoolId][teacherId].Subjects,
			SubjectAbsenceSummary{
				Subject: subjects.Subject{
					SubjectId:   subjectId,
					SubjectName: subjectName,
				},
				AbsenceSummary: AbsenceSummary{
					NumberOfAbsences:   totalAbsences,
					NumberOfExcuses:    excusedAbsences,
					NumberOfNonExcuses: nonExcusedAbsences,
				},
			},
		)
	}

	if err = rows.Err(); err != nil {
		drs.logger.Error("Error iterating absence rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	report := AbsencesReport{}
	for schoolId, schoolData := range schoolMap {
		for _, teacherData := range teacherMap[schoolId] {
			schoolData.Teachers = append(schoolData.Teachers, *teacherData)
		}
		report = append(report, *schoolData)
	}

	return report, nil
}

func (drs *DynamicReportingService) queryForGrades(ctx context.Context, reportQuery *ReportQueryRequest) (DynamicReportingResponse, *exceptions.AppError) {
	sql := `
		SELECT 
			s.school_id, s.school_name, s.school_address,
			t.teacher_id, u_t.first_name as teacher_first_name, u_t.last_name as teacher_last_name, u_t.email as teacher_email,
			subj.subject_id, subj.subject_name,
			COUNT(g.grade_id) as total_grades,
			SUM(CASE WHEN g.grade_value >= 5.50 THEN 1 ELSE 0 END) as sixes,
			SUM(CASE WHEN g.grade_value >= 4.50 AND g.grade_value < 5.50 THEN 1 ELSE 0 END) as fives,
			SUM(CASE WHEN g.grade_value >= 3.50 AND g.grade_value < 4.50 THEN 1 ELSE 0 END) as fours,
			SUM(CASE WHEN g.grade_value >= 3.00 AND g.grade_value < 3.50 THEN 1 ELSE 0 END) as threes,
			SUM(CASE WHEN g.grade_value >= 2.00 AND g.grade_value < 3.00 THEN 1 ELSE 0 END) as twos
		FROM grades g
		INNER JOIN curricula cur ON g.curriculum_id = cur.curriculum_id
		INNER JOIN classes cl ON cur.class_id = cl.class_id
		INNER JOIN schools s ON cl.school_id = s.school_id
		INNER JOIN teachers t ON cur.teacher_id = t.teacher_id
		INNER JOIN users u_t ON t.user_id = u_t.user_id
		INNER JOIN subjects subj ON cur.subject_id = subj.subject_id
		WHERE 1=1
	`

	var args []interface{}
	argIndex := 1

	if len(reportQuery.SchoolIds) > 0 {
		sql += fmt.Sprintf(" AND s.school_id = ANY($%d)", argIndex)
		args = append(args, reportQuery.SchoolIds)
		argIndex++
	}

	if len(reportQuery.TeacherIds) > 0 {
		sql += fmt.Sprintf(" AND t.teacher_id = ANY($%d)", argIndex)
		args = append(args, reportQuery.TeacherIds)
		argIndex++
	}

	if len(reportQuery.SubjectIds) > 0 {
		sql += fmt.Sprintf(" AND subj.subject_id = ANY($%d)", argIndex)
		args = append(args, reportQuery.SubjectIds)
		argIndex++
	}

	sql += `
		GROUP BY s.school_id, s.school_name, s.school_address, 
		         t.teacher_id, u_t.first_name, u_t.last_name, u_t.email,
		         subj.subject_id, subj.subject_name
		ORDER BY s.school_name, u_t.last_name, u_t.first_name, subj.subject_name
	`

	rows, err := drs.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		drs.logger.Error("Failed to query grades", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	schoolMap := make(map[int32]*SchoolGrades)
	teacherMap := make(map[int32]map[int32]*TeacherSubjectsGrades)

	for rows.Next() {
		var (
			schoolId, teacherId, subjectId                                                          int32
			schoolName, schoolAddress, teacherFirstName, teacherLastName, teacherEmail, subjectName string
			totalGrades, sixes, fives, fours, threes, twos                                          int32
		)

		err := rows.Scan(
			&schoolId, &schoolName, &schoolAddress,
			&teacherId, &teacherFirstName, &teacherLastName, &teacherEmail,
			&subjectId, &subjectName,
			&totalGrades, &sixes, &fives, &fours, &threes, &twos,
		)
		if err != nil {
			drs.logger.Error("Failed to scan grade row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}

		if _, exists := schoolMap[schoolId]; !exists {
			schoolMap[schoolId] = &SchoolGrades{
				School: schools.School{
					SchoolId:      schoolId,
					SchoolName:    schoolName,
					SchoolAddress: schoolAddress,
				},
				Teachers: []TeacherSubjectsGrades{},
			}
			teacherMap[schoolId] = make(map[int32]*TeacherSubjectsGrades)
		}

		if _, exists := teacherMap[schoolId][teacherId]; !exists {
			teacherMap[schoolId][teacherId] = &TeacherSubjectsGrades{
				Teacher: teachers.Teacher{
					TeacherId: teacherId,
					UserId:    0, // Not needed for report
					FirstName: teacherFirstName,
					LastName:  teacherLastName,
					Email:     teacherEmail,
				},
				Subjects: []SubjectGradeSummary{},
			}
		}

		teacherMap[schoolId][teacherId].Subjects = append(
			teacherMap[schoolId][teacherId].Subjects,
			SubjectGradeSummary{
				Subject: subjects.Subject{
					SubjectId:   subjectId,
					SubjectName: subjectName,
				},
				GradeSummary: GradeSummary{
					NumberOfGrades: totalGrades,
					GradeDistribution: GradeDistribution{
						NumberOfSixes:  sixes,
						NumberOfFives:  fives,
						NumberOfFours:  fours,
						NumberOfThrees: threes,
						NumberOfTwos:   twos,
					},
				},
			},
		)
	}

	if err = rows.Err(); err != nil {
		drs.logger.Error("Error iterating grade rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	report := GradesReport{}
	for schoolId, schoolData := range schoolMap {
		for _, teacherData := range teacherMap[schoolId] {
			schoolData.Teachers = append(schoolData.Teachers, *teacherData)
		}
		report = append(report, *schoolData)
	}

	return report, nil
}
