package usecase

import (
	"context"
	"testing"
	"time"

	outPkg "github.com/ReilEgor/Vaca/pkg"
	"github.com/ReilEgor/Vaca/services/CoordinatorService/internal/domain"
	mocks "github.com/ReilEgor/Vaca/services/CoordinatorService/internal/mocks/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Usecase_GetTaskStatus(t *testing.T) {
	test := []struct {
		name      string
		mockSetup func(s *mocks.StatusRepository)
		expected  *outPkg.Task
		wantErr   error
	}{
		{
			name: "success",
			mockSetup: func(s *mocks.StatusRepository) {
				s.On("Get", mock.Anything, outPkg.IntToUUID(1).String()).Return(
					map[string]string{"status": "completed"},
				)
			},
			expected: &outPkg.Task{
				ID:        outPkg.IntToUUID(1),
				Status:    "completed",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			wantErr: nil,
		},
		{
			name: "empty status",
			mockSetup: func(s *mocks.StatusRepository) {
				s.On("Get", mock.Anything, outPkg.IntToUUID(1).String()).Return(
					map[string]string{
						"error": "error",
					},
				)
			},
			expected: nil,
			wantErr:  domain.ErrTaskNotFound,
		},
	}
	for _, tt := range test {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sMock := mocks.NewStatusRepository(t)
			rMock := mocks.NewTaskPublisher(t)
			fMock := mocks.NewVacancySearchRepository(t)
			tt.mockSetup(sMock)
			uc := NewCoordinatorUsecase(sMock, rMock, fMock)
			gotTask, err := uc.GetTaskStatus(context.Background(), outPkg.IntToUUID(1).String())

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.expected, gotTask)
		})
	}
}

func Test_Usecase_CreateTask(t *testing.T) {
	t.Parallel()
	id := outPkg.IntToUUID(1)
	keywords := []string{"golang"}
	sources := []string{"dou.ua"}
	searchKey := generateSearchKey(keywords, sources)
	test := []struct {
		name        string
		mockSetup   func(s *mocks.StatusRepository)
		setupBroker func(b *mocks.TaskPublisher)
		keywords    []string
		sources     []string
		expected    *uuid.UUID
		wantErr     error
	}{
		{
			name:     "existing id",
			keywords: keywords,
			sources:  sources,
			mockSetup: func(s *mocks.StatusRepository) {

				s.On("GetIDByHash", mock.Anything, searchKey).Return(
					id.String(), nil,
				)
			},
			setupBroker: func(s *mocks.TaskPublisher) {

			},
			expected: &id,
			wantErr:  nil,
		},
		{
			name:     "not existing id",
			keywords: keywords,
			sources:  sources,
			mockSetup: func(s *mocks.StatusRepository) {
				s.On("GetIDByHash", mock.Anything, searchKey).Return(
					"", nil,
				)
				s.On("Set", mock.Anything, mock.Anything, searchKey, len(sources), time.Minute*2).Return(nil)
			},
			setupBroker: func(s *mocks.TaskPublisher) {
				s.On("PublishTask", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expected: nil,
			wantErr:  nil,
		},
	}
	for _, tt := range test {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sMock := mocks.NewStatusRepository(t)
			rMock := mocks.NewTaskPublisher(t)
			fMock := mocks.NewVacancySearchRepository(t)
			tt.mockSetup(sMock)
			tt.setupBroker(rMock)
			uc := NewCoordinatorUsecase(sMock, rMock, fMock)
			ans, err := uc.CreateTask(context.Background(), tt.keywords, tt.sources)

			assert.Equal(t, tt.wantErr, err)

			if tt.expected != nil {
				assert.Equal(t, tt.expected, ans)
			} else {
				assert.NotNil(t, ans)
			}
		})
	}
}

func Test_Usecase_GetVacancies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		expected        []*outPkg.Vacancy
		searchRepoSetup func(s *mocks.VacancySearchRepository)
		wantErr         error
	}{
		{
			name:     "success",
			expected: []*outPkg.Vacancy{},
			searchRepoSetup: func(s *mocks.VacancySearchRepository) {
				s.On("Search", mock.Anything, outPkg.VacancyFilter{}).Return(
					[]*outPkg.Vacancy{}, nil,
				)
			},
			wantErr: nil,
		},
		{
			name:     "search repository error",
			expected: nil,
			searchRepoSetup: func(s *mocks.VacancySearchRepository) {
				s.On("Search", mock.Anything, outPkg.VacancyFilter{}).Return(
					nil, domain.ErrSearchFailed,
				)
			},
			wantErr: domain.ErrSearchFailed,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sMock := mocks.NewStatusRepository(t)
			rMock := mocks.NewTaskPublisher(t)
			fMock := mocks.NewVacancySearchRepository(t)
			tt.searchRepoSetup(fMock)

			uc := NewCoordinatorUsecase(sMock, rMock, fMock)

			ans, _, err := uc.GetVacancies(context.Background(), outPkg.VacancyFilter{})
			assert.Equal(t, tt.wantErr, err)

			if tt.expected != nil {
				assert.Equal(t, tt.expected, ans)
			} else {
				assert.Nil(t, ans)
			}
		})
	}
}

func Test_Usecase_GetAvailableSources(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(s *mocks.StatusRepository)
		expected  []outPkg.Source
	}{
		{
			name: "success",
			mockSetup: func(s *mocks.StatusRepository) {
				s.On("GetSources", mock.Anything).Return(
					[]outPkg.Source{
						{
							ID:       outPkg.IntToUUID(1),
							Name:     "test",
							URL:      "http://test.com",
							IsActive: true,
						},
					}, nil,
				)
			},
			expected: []outPkg.Source{
				{
					ID:       outPkg.IntToUUID(1),
					Name:     "test",
					URL:      "http://test.com",
					IsActive: true,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sMock := mocks.NewStatusRepository(t)
			rMock := mocks.NewTaskPublisher(t)
			fMock := mocks.NewVacancySearchRepository(t)
			tt.mockSetup(sMock)

			uc := NewCoordinatorUsecase(sMock, rMock, fMock)
			ans, _, err := uc.GetAvailableSources(context.Background())
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, ans)
		})

	}
}
