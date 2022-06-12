package courses

import (
	"bytes"
	"encoding/json"
	"errors"
	mooc "github.com/CodelyTV/go-hexagonal_http_api-course/02-04-domain-validations/internal"
	"github.com/CodelyTV/go-hexagonal_http_api-course/02-04-domain-validations/internal/platform/storage/storagemocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListHandler_Ok(t *testing.T) {
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
	coursesList := []mooc.Course{
		course1,
		course2,
		course3,
	}
	expectedList := List{
		ListItem{
			ID:       course1ID,
			Name:     course1Name,
			Duration: course1Dur,
		},
		ListItem{
			ID:       course2ID,
			Name:     course2Name,
			Duration: course2Dur,
		},
		ListItem{
			ID:       course3ID,
			Name:     course3Name,
			Duration: course3Dur,
		},
	}
	repo := new(storagemocks.CourseRepository)
	repo.On("List", mock.Anything).Return(coursesList, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/courses", ListHandler(repo))

	req, err := http.NewRequest(http.MethodGet, "/courses", bytes.NewBuffer([]byte{}))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	list := List{}
	json.Unmarshal(body, &list)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expectedList, list)
}

func TestListHandler_Err(t *testing.T) {
	repoErr := errors.New("error on repository")
	repo := new(storagemocks.CourseRepository)
	repo.On("List", mock.Anything).Return([]mooc.Course{}, repoErr)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/courses", ListHandler(repo))

	req, err := http.NewRequest(http.MethodGet, "/courses", bytes.NewBuffer([]byte{}))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
