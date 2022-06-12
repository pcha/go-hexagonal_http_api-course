package listing

import (
	"context"
	"errors"
	"github.com/CodelyTV/go-hexagonal_http_api-course/04-03-command-bus/kit/query"
)

const CourseQueryType query.Type = "query.listing.course"

type CourseQuery struct {
}

func NewsCourseQuery() CourseQuery {
	return CourseQuery{}
}

func (q CourseQuery) Type() query.Type {
	return CourseQueryType
}

type CourseQueryHandler struct {
	service CourseService
}

func NewCourseQueryHandler(service CourseService) CourseQueryHandler {
	return CourseQueryHandler{service: service}
}

func (h CourseQueryHandler) Handle(ctx context.Context, qry query.Query) (query.Result, error) {
	_, ok := qry.(CourseQuery)
	if !ok {
		return nil, errors.New("unexpected query")
	}

	return h.service.ListCourses(ctx)
}
