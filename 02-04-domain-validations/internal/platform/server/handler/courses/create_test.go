package courses

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CodelyTV/go-hexagonal_http_api-course/02-04-domain-validations/internal/platform/storage/storagemocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Create(t *testing.T) {
	courseRepository := new(storagemocks.CourseRepository)
	courseRepository.On("Save", mock.Anything, mock.AnythingOfType("mooc.Course")).Return(nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/courses", CreateHandler(courseRepository))

	t.Run("given an invalid request it returns 400", func(t *testing.T) {
		type testCase struct {
			createCourseReq createRequest
			expectedMessage string
		}

		cases := map[string]testCase{
			"No ID": {
				createCourseReq: createRequest{
					Name:     "Demo Curse",
					Duration: "10 months",
				},
				expectedMessage: emptyIDMsg,
			},
			"No UIID ID": {
				createCourseReq: createRequest{
					ID:       "123",
					Name:     "Demo Course",
					Duration: "10 months",
				},
				expectedMessage: invalidIDMsg,
			},
			"No Name": {
				createCourseReq: createRequest{
					ID:       "adb77d46-ffe7-44a8-8520-49eca6dc2ed8",
					Duration: "10 months",
				},
				expectedMessage: emptyNameMsg,
			},
			"Too short Name": {
				createCourseReq: createRequest{
					ID:       "adb77d46-ffe7-44a8-8520-49eca6dc2ed8",
					Name:     "asdf",
					Duration: "10 months",
				},
				expectedMessage: tooShortNameMsg,
			},
			"No Duration": {
				createCourseReq: createRequest{
					ID:   "adb77d46-ffe7-44a8-8520-49eca6dc2ed8",
					Name: "Demo Course",
				},
				expectedMessage: emptyDurationMsg,
			},
			"Invalid Duration": {
				createCourseReq: createRequest{
					ID:       "adb77d46-ffe7-44a8-8520-49eca6dc2ed8",
					Name:     "Demo Course",
					Duration: "12",
				},
				expectedMessage: invalidDurationMsg,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				b, err := json.Marshal(tc.createCourseReq)
				require.NoError(t, err)

				req, err := http.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(b))
				require.NoError(t, err)

				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)

				res := rec.Result()
				defer res.Body.Close()

				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				errResp := errorResponse{}
				err = json.Unmarshal(body, &errResp)
				require.NoError(t, err)

				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
				assert.Equal(t, tc.expectedMessage, errResp.Error)
			})
		}
	})

	t.Run("given a valid request it returns 201", func(t *testing.T) {
		createCourseReq := createRequest{
			ID:       "8a1c5cdc-ba57-445a-994d-aa412d23723f",
			Name:     "Demo Course",
			Duration: "10 months",
		}

		b, err := json.Marshal(createCourseReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}
