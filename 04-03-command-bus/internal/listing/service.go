package listing

import (
	"context"
	mooc "github.com/CodelyTV/go-hexagonal_http_api-course/04-03-command-bus/internal"
)

type CourseService struct {
	courseRepository mooc.CourseRepository
}

func NewCourseService(courseRepository mooc.CourseRepository) CourseService {
	return CourseService{
		courseRepository: courseRepository,
	}
}

func (s CourseService) ListCourses(ctx context.Context) ([]mooc.Course, error) {
	return s.courseRepository.List(ctx)
}
