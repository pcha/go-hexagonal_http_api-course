package mysql

import (
	"context"
	"errors"
	"testing"

	mooc "github.com/CodelyTV/go-hexagonal_http_api-course/03-02-repository-test/internal"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CourseRepository_Save_RepositoryError(t *testing.T) {
	courseID, courseName, courseDuration := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test Course", "10 months"
	course, err := mooc.NewCourse(courseID, courseName, courseDuration)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO courses (id, name, duration) VALUES (?, ?, ?)").
		WithArgs(courseID, courseName, courseDuration).
		WillReturnError(errors.New("something-failed"))

	repo := NewCourseRepository(db)

	err = repo.Save(context.Background(), course)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

func Test_CourseRepository_Save_Succeed(t *testing.T) {
	courseID, courseName, courseDuration := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test Course", "10 months"
	course, err := mooc.NewCourse(courseID, courseName, courseDuration)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO courses (id, name, duration) VALUES (?, ?, ?)").
		WithArgs(courseID, courseName, courseDuration).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewCourseRepository(db)

	err = repo.Save(context.Background(), course)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func Test_List_RepositoryError(t *testing.T) {
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectQuery("SELECT courses.id, courses.name, courses.duration FROM courses").
		WillReturnError(errors.New("error duerying database"))

	repo := NewCourseRepository(db)

	_, err = repo.List(context.Background())

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

func Test_List_RepositorySucceed(t *testing.T) {
	const course1ID = "adb77d46-ffe7-44a8-8520-49eca6dc2ed8"
	const course1Name = "Course 1"
	const course1Dur = "10 months"
	course1, err := mooc.NewCourse(course1ID, course1Name, course1Dur)
	require.NoError(t, err)
	const course2ID = "9d7c9779-d12e-4faf-81b8-824457227a02"
	const course2Name = "Course 2"
	const course2Dur = "10 months"
	course2, err := mooc.NewCourse(course2ID, course2Name, course2Dur)
	require.NoError(t, err)
	const course3ID = "802ce374-aeb4-4a35-89d5-eb12e55299a1"
	const course3Name = "Course 3"
	const course3Dur = "10 months"
	course3, err := mooc.NewCourse(course3ID, course3Name, course3Dur)
	require.NoError(t, err)
	wantedCourses := []mooc.Course{
		course1,
		course2,
		course3,
	}

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectQuery("SELECT courses.id, courses.name, courses.duration FROM courses").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "duration"}).
			AddRow(course1ID, course1Name, course1Dur).
			AddRow(course2ID, course2Name, course2Dur).
			AddRow(course3ID, course3Name, course3Dur))

	repo := NewCourseRepository(db)

	courses, err := repo.List(context.Background())

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Equal(t, wantedCourses, courses)
}
