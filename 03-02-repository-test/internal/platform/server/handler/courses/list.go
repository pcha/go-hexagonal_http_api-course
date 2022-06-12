package courses

import (
	mooc "github.com/CodelyTV/go-hexagonal_http_api-course/03-02-repository-test/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type List []ListItem

func ListHandler(repository mooc.CourseRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		courses, err := mooc.ListCourses(ctx, repository)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		var list List
		for _, c := range courses {
			list = append(list, ListItem{
				ID:       c.ID().String(),
				Name:     c.Name().String(),
				Duration: c.Duration().String(),
			})
		}

		ctx.JSON(http.StatusOK, list)
	}
}
