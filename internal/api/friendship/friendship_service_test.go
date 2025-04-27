package friendship_test

import (
	"errors"
	"socialAPI/internal/api/friendship"
	"socialAPI/internal/mocks"
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	errExample      = errors.New("example error")
	incorrectStatus = "incorrect status"
	friendshipID    = uint(1)
	senderID        = uint(1)
	receiverID      = uint(2)
	amity           = r.Friendship{SenderID: senderID, ReceiverID: receiverID, Status: r.StatusPending}
	users           = []*r.FriendWithID{
		{
			Friend:       &r.User{ID: receiverID},
			FriendshipID: friendshipID,
		},
	}
)

type friendshipServiceMocks struct {
	friendshipRepo *mocks.FriendshipRepository
	friendshipSrv  friendship.FriendshipService
}

func setupFriendshipService() friendshipServiceMocks {
	friedshipRepo := new(mocks.FriendshipRepository)
	logger := zap.NewNop().Sugar()

	srv := friendship.NewFriendshipService(friedshipRepo, logger)

	return friendshipServiceMocks{
		friendshipRepo: friedshipRepo,
		friendshipSrv:  srv,
	}
}

func TestFriendshipService_SendFriendRequest(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m friendshipServiceMocks)
		wantErr    bool
		errMessage string
	}{
		{
			name:       "cannot send friend request to yourself",
			setup:      func(m friendshipServiceMocks) {},
			errMessage: "Cannot send friend request to yourself",
			wantErr:    true,
		},
		{
			name: "error checking if friendship exists",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("Exists", senderID, receiverID).Return(false, errExample)
			},
			errMessage: shared.InternalError.Error(),
			wantErr:    true,
		},
		{
			name: "error checking if friendship exists",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("Exists", senderID, receiverID).Return(true, nil)
			},
			errMessage: "friendship exists",
			wantErr:    true,
		},
		{
			name: "error checking if friendship exists",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("Exists", senderID, receiverID).Return(false, nil)
				m.friendshipRepo.On("SendRequest", &amity).Return(errExample)
			},
			errMessage: shared.InternalError.Error(),
			wantErr:    true,
		},
		{
			name: "friend request successfully sent",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("Exists", senderID, receiverID).Return(false, nil)
				m.friendshipRepo.On("SendRequest", &amity).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := setupFriendshipService()
			tt.setup(m)

			id1 := senderID
			id2 := receiverID
			if tt.name == "cannot send friend request to yourself" {
				id2 = senderID
			}

			err := m.friendshipSrv.SendFriendRequest(id1, id2)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}
			m.friendshipRepo.AssertExpectations(t)
		})
	}
}
func TestFriendshipService_GetAllFriends(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m friendshipServiceMocks)
		wantErr    bool
		wantUsers  bool
		errMessage string
		status     string
	}{
		{
			name:       "incorrect status",
			setup:      func(m friendshipServiceMocks) {},
			errMessage: "incorrect status",
			wantErr:    true,
			status:     incorrectStatus,
			wantUsers:  false,
		},
		{
			name: "error getting list of friends",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("GetAllFriends", senderID, mock.Anything).Return(nil, errExample)
			},
			errMessage: shared.InternalError.Error(),
			wantErr:    true,
			status:     string(r.StatusPending),
			wantUsers:  false,
		},
		{
			name: "successfully retrieved list of pending guys",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("GetAllFriends", senderID, mock.Anything).Return(users, nil)
			},
			wantErr:   false,
			status:    string(r.StatusPending),
			wantUsers: true,
		},
		{
			name: "successfully retrieved list of rejected guys",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("GetAllFriends", senderID, mock.Anything).Return(users, nil)
			},
			wantErr:   false,
			status:    string(r.StatusRejected),
			wantUsers: true,
		},
		{
			name: "successfully retrieved list of friends",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("GetAllFriends", senderID, mock.Anything).Return(users, nil)
			},
			wantErr:   false,
			status:    string(r.StatusFriendship),
			wantUsers: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := setupFriendshipService()
			tt.setup(m)

			users, err := m.friendshipSrv.GetAllFriends(senderID, tt.status)

			if tt.wantUsers {
				assert.NotNil(t, users)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}

			m.friendshipRepo.AssertExpectations(t)
		})
	}
}
func TestFriendshipService_PatchFriendship(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m friendshipServiceMocks)
		wantErr    bool
		errMessage string
	}{
		{
			name: "friend request not found or not yours",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("SetStatus", friendshipID, r.StatusRejected).Return(gorm.ErrRecordNotFound)
			},
			errMessage: "Friend request not found or not yours",
			wantErr:    true,
		},
		{
			name: "error updating friendship status",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("SetStatus", friendshipID, r.StatusRejected).Return(errExample)
			},
			errMessage: shared.InternalError.Error(),
			wantErr:    true,
		},
		{
			name: "friendship status successfully updated",
			setup: func(m friendshipServiceMocks) {
				m.friendshipRepo.On("SetStatus", friendshipID, r.StatusRejected).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := setupFriendshipService()
			tt.setup(m)

			err := m.friendshipSrv.PatchFriendship(friendshipID, friendship.ChangeStatusRequest{r.StatusRejected})

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}

			m.friendshipRepo.AssertExpectations(t)
		})
	}
}
