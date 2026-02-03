package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	mocks "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/mocks/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Handler_GetVacancies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(uc *mocks.CoordinatorUsecase)
		expectedStatus int
		url            string
		expectedResp   SearchVacanciesResponse
		expectedError  map[string]string
	}{
		{
			name: "success",
			url:  "/api/v1/vacancies?query=Golang+Developer&location=Kyiv&requirements=Experience+with+REST+APIs&limit=10&offset=0",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {
				uc.On("GetVacancies", mock.Anything, outPkg.VacancyFilter{
					Query:        "Golang Developer",
					Location:     "Kyiv",
					Requirements: "Experience with REST APIs",
					Limit:        10,
					Offset:       0,
				}).Return(
					[]*outPkg.Vacancy{
						{
							ID:           outPkg.IntToUUID(1),
							Title:        "Golang Developer",
							Company:      "Tech Corp",
							Description:  "",
							Link:         "https://techcorp.com/careers/123",
							About:        "Join our team as a Golang Developer.",
							Requirements: "Experience with REST APIs",
						},
					},
					int64(1),
					nil)
			},
			expectedStatus: http.StatusOK,
			expectedResp: SearchVacanciesResponse{
				Items: []VacancyResponse{
					{
						ID:           outPkg.IntToUUID(1),
						Title:        "Golang Developer",
						Company:      "Tech Corp",
						Link:         "https://techcorp.com/careers/123",
						About:        "Join our team as a Golang Developer.",
						Requirements: "Experience with REST APIs",
					},
				},
				Total: int64(1),
			},
			expectedError: nil,
		},
		{
			name: "invalid query param",
			url:  "/api/v1/vacancies?query=Golang+Developer&location=Kyiv&requirements=Experience+with+REST+APIs&limit=10&offset=abc",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {

			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   SearchVacanciesResponse{},
			expectedError: map[string]string{
				"error": domain.ErrInvalidRequestBody.Error(),
			},
		},
		{
			name: "uc error",
			url:  "/api/v1/vacancies?query=Golang+Developer&location=Kyiv&requirements=Experience+with+REST+APIs&limit=10&offset=0",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {
				uc.On("GetVacancies", mock.Anything, outPkg.VacancyFilter{
					Query:        "Golang Developer",
					Location:     "Kyiv",
					Requirements: "Experience with REST APIs",
					Limit:        10,
					Offset:       0,
				}).Return(
					nil,
					int64(1),
					errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   SearchVacanciesResponse{},
			expectedError: map[string]string{
				"error": domain.ErrFailedToGetVacancies.Error(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := mocks.NewCoordinatorUsecase(t)
			tt.mockSetup(uc)

			handler := NewHandler(uc)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			c.Request = req

			handler.GetVacancies(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				expectedJSON, err := json.Marshal(tt.expectedResp)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedJSON), w.Body.String())
			} else {
				expectedJSON, err := json.Marshal(tt.expectedError)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedJSON), w.Body.String())
			}
		})
	}
}
