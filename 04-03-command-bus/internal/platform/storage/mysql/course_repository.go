package mysql

import (
	"context"
	"database/sql"
	"fmt"

	mooc "github.com/CodelyTV/go-hexagonal_http_api-course/04-03-command-bus/internal"
	"github.com/huandu/go-sqlbuilder"
)

// CourseRepository is a MySQL mooc.CourseRepository implementation.
type CourseRepository struct {
	db *sql.DB
}

// NewCourseRepository initializes a MySQL-based implementation of mooc.CourseRepository.
func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

// Save implements the mooc.CourseRepository interface.
func (r *CourseRepository) Save(ctx context.Context, course mooc.Course) error {
	courseSQLStruct := sqlbuilder.NewStruct(new(sqlCourse))
	query, args := courseSQLStruct.InsertInto(sqlCourseTable, sqlCourse{
		ID:       course.ID().String(),
		Name:     course.Name().String(),
		Duration: course.Duration().String(),
	}).Build()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist course on database: %v", err)
	}

	return nil
}

func (r *CourseRepository) List(ctx context.Context) ([]mooc.Course, error) {
	courseSQLStruct := sqlbuilder.NewStruct(new(sqlCourse))
	q, _ := courseSQLStruct.SelectFrom(sqlCourseTable).Build()

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("error tryin to query courses on database, %v", err)
	}

	defer rows.Close()
	list := []mooc.Course{}

	for rows.Next() {
		var sqlCour sqlCourse
		rows.Scan(courseSQLStruct.Addr(&sqlCour))

		course, err := mooc.NewCourse(sqlCour.ID, sqlCour.Name, sqlCour.Duration)
		if err != nil {
			return nil, fmt.Errorf("error parsing obtained result: %w", err)
		}

		list = append(list, course)
	}

	return list, nil
}
