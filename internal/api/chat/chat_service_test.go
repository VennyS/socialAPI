package chat_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"socialAPI/internal/api/chat"
	"socialAPI/internal/mocks"
	"socialAPI/internal/shared"
	"socialAPI/internal/storage/repository"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

var (
	userID       = uint(1)
	chatID       = uint(1)
	chatExample  = repository.Chat{ID: uint(1)}
	chatsExample = []*repository.Chat{
		&chatExample,
	}
	createRequestExample = chat.CreateRequest{UserIDs: []uint{1, 2}}
)

type chatServiceMocks struct {
	userRepo   *mocks.UserRepository
	chatRepo   *mocks.ChatRepository
	hub        *mocks.Hub
	wsUpgrader *mocks.Upgrader
	logger     *zap.SugaredLogger
	chatSrv    chat.ChatService
}

func setupChatService() chatServiceMocks {
	userRepo := new(mocks.UserRepository)
	chatRepo := new(mocks.ChatRepository)
	logger := zap.NewNop().Sugar()
	hub := new(mocks.Hub)
	wsUpgrader := new(mocks.Upgrader)

	chatSrv := chat.NewChatService(chatRepo, userRepo, hub, wsUpgrader, logger)

	return chatServiceMocks{
		userRepo:   userRepo,
		chatRepo:   chatRepo,
		hub:        hub,
		wsUpgrader: wsUpgrader,
		logger:     logger,
		chatSrv:    chatSrv,
	}
}

var (
	errExample = errors.New("example error")
)

func TestChatService_GetOne(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m *chatServiceMocks)
		wantErr    bool
		wantChat   bool
		errMessage string
	}{
		{
			name: "failed to check chat existence",
			setup: func(m *chatServiceMocks) {
				m.chatRepo.On("ExistsID", chatID).Return(false, errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
			wantChat:   false,
		},
		{
			name: "chat not found",
			setup: func(m *chatServiceMocks) {
				m.chatRepo.On("ExistsID", chatID).Return(false, nil)
			},
			wantErr:    true,
			errMessage: "chat not found",
			wantChat:   false,
		},
		{
			name: "failed to fetch chat",
			setup: func(m *chatServiceMocks) {
				m.chatRepo.On("ExistsID", chatID).Return(true, nil)
				m.chatRepo.On("GetOne", chatID).Return(nil, errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
			wantChat:   false,
		},
		{
			name: "chat succefully fetched and converted to DTO",
			setup: func(m *chatServiceMocks) {
				m.chatRepo.On("ExistsID", chatID).Return(true, nil)
				m.chatRepo.On("GetOne", chatID).Return(&chatExample, nil)
			},
			wantErr:  false,
			wantChat: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := setupChatService()
			tt.setup(&mocks)

			chatDTO, err := mocks.chatSrv.GetOne(chatID)
			if tt.wantChat {
				assert.NotNil(t, chatDTO)
				assert.NotZero(t, chatDTO.ID)
			} else {
				assert.Nil(t, chatDTO)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}
			mocks.userRepo.AssertExpectations(t)
		})
	}
}

func TestChatService_GetAll(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m *chatServiceMocks)
		wantErr    bool
		wantChats  bool
		errMessage string
	}{
		{
			name: "Failed to fetch chats",
			setup: func(m *chatServiceMocks) {
				m.chatRepo.On("GetAll").Return(nil, errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
			wantChats:  false,
		},
		{
			name: "chats succefully fetched and converted to DTO",
			setup: func(m *chatServiceMocks) {
				m.chatRepo.On("GetAll").Return(chatsExample, nil)
			},
			wantErr:   false,
			wantChats: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := setupChatService()
			tt.setup(&mocks)

			chatsDTO, err := mocks.chatSrv.GetAll()
			if tt.wantChats {
				assert.NotNil(t, chatsDTO)
				assert.NotEmpty(t, chatsDTO)
			} else {
				assert.Nil(t, chatsDTO)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}
			mocks.userRepo.AssertExpectations(t)
		})
	}
}
func TestChatService_Create(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m *chatServiceMocks)
		wantErr    bool
		errMessage string
	}{
		{
			name: "error checking user IDs existence",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(false, errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "some user IDs do not exist",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(false, nil)
			},
			wantErr:    true,
			errMessage: "some user IDs do not exist",
		},
		{
			name: "error checking chat with this userIDs",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(false, errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "chat with the same users already exists",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(true, nil)
			},
			wantErr:    true,
			errMessage: "chat with the same users already exists",
		},
		{
			name: "error creating chat",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(false, nil)
				m.chatRepo.On("Create", createRequestExample.Name, createRequestExample.UserIDs).Return(errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "chat created successfully",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(false, nil)
				m.chatRepo.On("Create", createRequestExample.Name, createRequestExample.UserIDs).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := setupChatService()
			tt.setup(&mocks)

			err := mocks.chatSrv.Create(createRequestExample)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}
			mocks.userRepo.AssertExpectations(t)
		})
	}
}
func TestChatService_Update(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m *chatServiceMocks)
		wantErr    bool
		errMessage string
	}{
		{
			name: "error checking user IDs existence",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(false, errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "some user IDs do not exist",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(false, nil)
			},
			wantErr:    true,
			errMessage: "some user IDs do not exist",
		},
		{
			name: "error checking chat with this userIDs",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(false, errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "chat with the same users already exists",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(true, nil)
			},
			wantErr:    true,
			errMessage: "chat with the same users already exists",
		},
		{
			name: "failed to check chat existence",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(false, nil)
				m.chatRepo.On("ExistsID", chatID).Return(true, errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "chat not found",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(false, nil)
				m.chatRepo.On("ExistsID", chatID).Return(false, nil)
			},
			wantErr:    true,
			errMessage: "chat not found",
		},
		{
			name: "error while updating chat",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(false, nil)
				m.chatRepo.On("ExistsID", chatID).Return(true, nil)
				m.chatRepo.On("Update", chatID, createRequestExample.Name, createRequestExample.UserIDs).Return(errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "chat updated successfully",
			setup: func(m *chatServiceMocks) {
				m.userRepo.On("IDsExists", createRequestExample.UserIDs).Return(true, nil)
				m.chatRepo.On("ExistsSetUserIDs", createRequestExample.UserIDs).Return(false, nil)
				m.chatRepo.On("ExistsID", chatID).Return(true, nil)
				m.chatRepo.On("Update", chatID, createRequestExample.Name, createRequestExample.UserIDs).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := setupChatService()
			tt.setup(&mocks)

			err := mocks.chatSrv.Update(chatID, createRequestExample)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}
			mocks.userRepo.AssertExpectations(t)
		})
	}
}
func TestChatService_HandleWebSocket(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m *chatServiceMocks)
		wantErr    bool
		errMessage string
	}{
		{
			name: "error getting chat IDs",
			setup: func(m *chatServiceMocks) {
				m.chatRepo.On("GetChatIDsByUserID", userID).Return(nil, errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "error upgrading to WebSocket",
			setup: func(m *chatServiceMocks) {
				m.chatRepo.On("GetChatIDsByUserID", userID).Return([]uint{chatID}, nil)
				m.wsUpgrader.On("Upgrade", mock.Anything, mock.Anything, mock.Anything).Return(nil, errExample)
			},
			wantErr:    true,
			errMessage: "WebSocket upgrade failed",
		},
		{
			name: "successful websocket connection",
			setup: func(m *chatServiceMocks) {
				m.chatRepo.On("GetChatIDsByUserID", userID).Return([]uint{chatID}, nil)
				m.wsUpgrader.On("Upgrade", mock.Anything, mock.Anything, mock.Anything).Return(&websocket.Conn{}, nil)
				m.hub.On("RegisterClient", mock.Anything).Return()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := setupChatService()
			tt.setup(&mocks)

			// Используем httptest.ResponseRecorder и httptest.NewRequest для моков
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)

			err := mocks.chatSrv.HandleWebSocket(userID, recorder, request)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}

			mocks.chatRepo.AssertExpectations(t)
			mocks.wsUpgrader.AssertExpectations(t)
			mocks.hub.AssertExpectations(t)
		})
	}
}
