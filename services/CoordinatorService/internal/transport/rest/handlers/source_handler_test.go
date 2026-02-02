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

func Test_Handler_GetAvailableSources(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(uc *mocks.CoordinatorUsecase)
		expectedStatus int
		expectedResp   ListSourcesResponse
		expectedError  map[string]string
	}{
		{
			name: "success",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {
				uc.On("GetAvailableSources", mock.Anything).Return(
					[]outPkg.Source{
						{ID: outPkg.IntToUUID(1), Name: "GitHub", URL: "https://github.com", IsActive: true},
						{ID: outPkg.IntToUUID(2), Name: "LinkedIn", URL: "https://linkedin.com", IsActive: true},
					}, int64(2), nil,
				)
			},
			expectedStatus: http.StatusOK,
			expectedResp: ListSourcesResponse{
				Sources: []SourceResponse{
					{ID: outPkg.IntToUUID(1), Name: "GitHub"},
					{ID: outPkg.IntToUUID(2), Name: "LinkedIn"},
				},
				Total: int64(2),
			},
			expectedError: nil,
		},
		{
			name: "usecase error",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {
				uc.On("GetAvailableSources", mock.Anything).Return(
					nil, int64(0), errors.New("db error"),
				)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   ListSourcesResponse{},
			expectedError: map[string]string{
				"error": domain.ErrFailedToGetSources.Error(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ucMock := mocks.NewCoordinatorUsecase(t)
			tt.mockSetup(ucMock)

			handler := NewHandler(ucMock)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/sources", nil)
			c.Request = req

			handler.GetAvailableSources(c)
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
