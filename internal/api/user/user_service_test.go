package user_test

import (
	"errors"
	"socialAPI/internal/api/user"
	"socialAPI/internal/mocks"
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

var (
	errExample = errors.New("example error")
	userID     = uint(1)
	users      = []r.User{
		{
			ID: uint(2),
		},
	}
)

type userServiceMocks struct {
	userRepo *mocks.UserRepository
	userSrv  user.UserService
}

func setupUserService() userServiceMocks {
	userRepo := new(mocks.UserRepository)
	logger := zap.NewNop().Sugar()

	srv := user.NewUserService(userRepo, logger)

	return userServiceMocks{
		userRepo: userRepo,
		userSrv:  srv,
	}
}

func TestFriendshipService_SendFriendRequest(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m userServiceMocks)
		wantErr    bool
		wantUsers  bool
		errMessage string
		id         *uint
	}{
		{
			name: "failed to fetch users",
			setup: func(m userServiceMocks) {
				m.userRepo.On("GetAll", mock.AnythingOfType("*uint")).Return(nil, errExample)
			},
			errMessage: shared.InternalError.Error(),
			wantErr:    true,
			wantUsers:  false,
			id:         &userID,
		},
		{
			name: "users successfully fetched",
			setup: func(m userServiceMocks) {
				m.userRepo.On("GetAll", mock.AnythingOfType("*uint")).Return(users, nil)
			},
			wantErr:   false,
			wantUsers: true,
			id:        &userID,
		},
		{
			name: "nil user ID",
			setup: func(m userServiceMocks) {
				m.userRepo.On("GetAll", (*uint)(nil)).Return(users, nil)
			},
			wantErr:   false,
			wantUsers: true,
			id:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := setupUserService()
			tt.setup(m)

			users, err := m.userSrv.GetAllUsers(tt.id)

			if tt.wantUsers {
				assert.NotNil(t, users)
			} else {
				assert.Nil(t, users)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}

			m.userRepo.AssertExpectations(t)
		})
	}
}
