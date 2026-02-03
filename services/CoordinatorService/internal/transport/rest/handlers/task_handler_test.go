package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	mocks "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/mocks/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Handler_GetTaskStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(uc *mocks.CoordinatorUsecase)
		expectedStatus int
		expectedResp   CreateTaskResponse
		expectedError  map[string]string
		withParam      bool
	}{
		{
			name: "success",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {
				uc.On("GetTaskStatus", mock.Anything, outPkg.IntToUUID(1).String()).Return(
					&outPkg.Task{
						ID:               outPkg.IntToUUID(1),
						KeyWords:         []string{"Golang", "Developer"},
						Status:           "InProgress",
						TotalSources:     2,
						CompletedSources: 0,
						CreatedAt:        time.Time{},
						UpdatedAt:        time.Time{},
					}, nil,
				)
			},
			expectedStatus: http.StatusOK,
			expectedResp: CreateTaskResponse{
				TaskID:    outPkg.IntToUUID(1),
				Status:    "InProgress",
				CreatedAt: time.Time{}.Format(time.RFC3339),
			},
			expectedError: nil,
			withParam:     true,
		},
		{
			name:           "empty param id",
			mockSetup:      func(uc *mocks.CoordinatorUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   CreateTaskResponse{},
			expectedError: map[string]string{
				"error": outPkg.ErrTaskIDRequired.Error(),
			},
			withParam: false,
		},
		{
			name: "usecase error",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {
				uc.On("GetTaskStatus", mock.Anything, outPkg.IntToUUID(1).String()).Return(
					&outPkg.Task{}, errors.New("uc error"),
				)
			},
			expectedStatus: http.StatusNotFound,
			expectedResp:   CreateTaskResponse{},
			expectedError: map[string]string{
				"error": domain.ErrTaskNotFound.Error(),
			},
			withParam: true,
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
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tasks/1", nil)
			c.Request = req
			if tt.withParam {
				c.Params = gin.Params{
					{Key: "id", Value: outPkg.IntToUUID(1).String()},
				}
			}

			handler.GetTaskStatus(c)
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

func Test_Handler_CreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	id := outPkg.IntToUUID(1)
	tests := []struct {
		name           string
		mockSetup      func(uc *mocks.CoordinatorUsecase)
		inputBody      CreateTaskRequest
		expectedStatus int
		expectedResp   CreateTaskResponse
		expectedError  map[string]string
	}{
		{
			name: "success",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {
				uc.On("CreateTask", mock.Anything, []string{"go", "junior"}, []string{"dou.ua"}).Return(
					&id, nil,
				)
			},
			inputBody: CreateTaskRequest{
				Keywords: []string{"go", "junior"},
				Sources:  []string{"dou.ua"},
			},
			expectedStatus: http.StatusCreated,
			expectedResp: CreateTaskResponse{
				TaskID:    outPkg.IntToUUID(1),
				Status:    "created",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectedError: nil,
		},
		{
			name:      "invalid request body",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {},
			inputBody: CreateTaskRequest{
				Sources: []string{"dou.ua"},
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   CreateTaskResponse{},
			expectedError: map[string]string{
				"error": outPkg.ErrInvalidRequest.Error(),
			},
		},
		{
			name: "usecase error",
			mockSetup: func(uc *mocks.CoordinatorUsecase) {
				uc.On("CreateTask", mock.Anything, []string{"go", "junior"}, []string{"dou.ua"}).Return(
					&id, errors.New("uc error"),
				)
			},
			inputBody: CreateTaskRequest{
				Keywords: []string{"go", "junior"},
				Sources:  []string{"dou.ua"},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   CreateTaskResponse{},
			expectedError: map[string]string{
				"error": domain.ErrFailedToCreateTask.Error(),
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
			jsonData, _ := json.Marshal(tt.inputBody)
			body := bytes.NewBuffer(jsonData)

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/tasks/", body)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.CreateTask(c)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusCreated {
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
