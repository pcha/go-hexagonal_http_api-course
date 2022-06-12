package courses

import (
	"errors"
	"net/http"
	"strings"

	mooc "github.com/CodelyTV/go-hexagonal_http_api-course/02-04-domain-validations/internal"
	"github.com/gin-gonic/gin"
)

type createRequest struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Duration string `json:"duration" binding:"required"`
}

type errorResponse struct {
	Error string `json:"error"`
}

const emptyIDMsg = "The field `id` can't be empty"
const invalidIDMsg = "The field `id` isn't a valid UUID"
const emptyNameMsg = "The field `name` can't be empty"
const tooShortNameMsg = "The field `name` must have at least 5 chars"
const emptyDurationMsg = "The field `duration` can't be empty"
const invalidDurationMsg = "The field `duration` has an invalid format, it should be expressed in years, months or weeks"

// CreateHandler returns an HTTP handler for courses creation.
func CreateHandler(courseRepository mooc.CourseRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createRequest
		if err := ctx.BindJSON(&req); err != nil {
			var errMsg string
			switch {
			case strings.Contains(err.Error(), "ID"):
				errMsg = emptyIDMsg
			case strings.Contains(err.Error(), "Name"):
				errMsg = emptyNameMsg
			case strings.Contains(err.Error(), "Duration"):
				errMsg = emptyDurationMsg
			}
			ctx.JSON(http.StatusBadRequest, errorResponse{Error: errMsg})
			return
		}

		course, err := mooc.NewCourse(req.ID, req.Name, req.Duration)
		if err != nil {
			var errMsg string
			switch {
			case errors.Is(err, mooc.ErrInvalidCourseID):
				errMsg = invalidIDMsg
			case errors.Is(err, mooc.ErrEmptyCourseName):
				errMsg = emptyNameMsg
			case errors.Is(err, mooc.ErrTooShortCourseName):
				errMsg = tooShortNameMsg
			case errors.Is(err, mooc.ErrEmptyDuration):
				errMsg = emptyDurationMsg
			case errors.Is(err, mooc.ErrInvalidFormatDuration):
				errMsg = invalidDurationMsg
			default:
				errMsg = err.Error()
			}
			ctx.JSON(http.StatusBadRequest, errorResponse{Error: errMsg})
			return
		}

		if err := courseRepository.Save(ctx, course); err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Status(http.StatusCreated)
	}
}
