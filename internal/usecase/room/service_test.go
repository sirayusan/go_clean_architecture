package usecase

import (
	"context"
	"time"

	"business/internal/entity"
	"github.com/stretchr/testify/mock"
)

type MockChatMessageRepository struct {
	mock.Mock
}

func (m *MockChatMessageRepository) GetMessageList(chatID uint32) ([]entity.Message, error) {
	args := m.Called(chatID)
	return args.Get(0).([]entity.Message), args.Error(1)
}

func (m *MockChatMessageRepository) CreateMessage(chatRoomID uint32, msg string, SenderUserID uint32) (entity.ChatMessage, error) {
	args := m.Called(chatRoomID)
	return args.Get(0).(entity.ChatMessage), args.Error(1)
}

type MockWebSocket struct {
	mock.Mock
}

// WriteMessage は MockConn インターフェースのメソッドです。
func (m *MockWebSocket) WriteMessage(messageType int, data []byte) error {
	args := m.Called(messageType, data)
	return args.Error(0)
}

// ReadMessage は MockWebSocket インターフェースのメソッドです。
func (m *MockWebSocket) ReadMessage() (messageType int, p []byte, err error) {
	args := m.Called() // 引数なしで Called を呼び出します

	// messageType を取得し、int へキャストします
	if args.Get(0) != nil {
		messageType = args.Get(0).(int)
	}

	// p (ペイロード) を取得し、[]byte へキャストします
	if args.Get(1) != nil {
		p = args.Get(1).([]byte)
	}

	// err を取得し、error へキャストします（存在する場合）
	err, _ = args.Error(2).(error)

	return messageType, p, err
}

type MockWebSocketConnWrapper struct {
	Conn entity.WebSocketWrapper
}

// ここでWebSocketWrapperインターフェースを実装していることを確認します。
var _ entity.WebSocketWrapper = (*MockWebSocketConnWrapper)(nil)

func (w *MockWebSocketConnWrapper) WriteMessage(messageType int, data []byte) error {
	return w.Conn.WriteMessage(messageType, data)
}

func (w *MockWebSocketConnWrapper) ReadMessage() (messageType int, p []byte, err error) {
	return w.Conn.ReadMessage()
}

type MockRedisWrapper struct {
	mock.Mock
}

func (m *MockRedisWrapper) Publish(ctx context.Context, channel string, message interface{}) *entity.IntCmd {
	args := m.Called(ctx, channel, message)
	return args.Get(0).(*entity.IntCmd)
}

func (m *MockRedisWrapper) Subscribe(ctx context.Context, channels ...string) *entity.PubSub {
	args := m.Called(ctx, channels)
	return args.Get(0).(*entity.PubSub)
}

type MockWrapperTime struct {
	mock.Mock
}

func (m *MockWrapperTime) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time) // モックから返される時刻を取得
}

// テストケース
//func TestMessageUseCase_JoinRoom(t *testing.T) {
//	mockRepo := new(MockChatMessageRepository)
//	uc := MessageUseCase{
//		repo: mockRepo,
//		wt:   ct.CustomTime{}, // カスタムタイムのインスタンスを注入
//	}
//	_time := time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local)
//
//	// メッセージリストのモック設定
//	expectedMessages := []entity.Message{{
//		SenderUserID: 1,
//		UserName:     "aaa",
//		Messages:     "テスト1",
//		CreatedAt:    _time,
//	}}
//	mockRepo.On("GetMessageList", uint32(1)).Return(expectedMessages, nil)
//
//	ws := &MockWebSocket{}
//	wsw := &MockWebSocketConnWrapper{Conn: ws}
//
//	// WebSocketWrapperのWriteMessageをモック
//	ws.On("WriteMessage", mock.Anything, mock.Anything).Return(nil)
//
//	roomManager := make(map[uint32]*entity.ChatRoom)
//	err := uc.JoinRoom(uint32(1), wsw, roomManager)
//	assert.NoError(t, err)
//
//	// クライアントがルームに追加されたことを確認
//	room, exists := roomManager[uint32(1)]
//	assert.True(t, exists)
//	assert.Equal(t, 1, len(room.Clients))
//
//	// メッセージが送信されたことを確認
//	ws.AssertCalled(t, "WriteMessage", websocket.TextMessage, mock.Anything)
//}
//
//func TestPubSub(t *testing.T) {
//	mockRepo := new(MockChatMessageRepository)
//	mockRedis := new(MockRedisWrapper)
//	// モックのカスタムタイムを設定
//	mockTime := new(MockWrapperTime)
//	fixedTime := time.Date(2024, time.February, 28, 12, 0, 0, 0, time.UTC)
//	mockTime.On("Now").Return(fixedTime)
//	uc := MessageUseCase{
//		repo: mockRepo,
//		wt:   mockTime, // カスタムタイムのインスタンスを注入
//	}
//	//_time := time.Date(2023, time.June, 19, 0, 0, 0, 0, time.Local)
//
//	chatRoomID := uint32(1)
//	_msg := "test message"
//	ws := &MockWebSocket{}
//	wsw := &MockWebSocketConnWrapper{Conn: ws}
//
//	// メッセージリストのモック設定
//	expectedMessages := entity.ChatMessage{
//		ChatMessageID: 1,
//		ChatRoomID:    1,
//		Message:       "テスト",
//		SenderUserID:  1,
//		CreatedAt:     uc.wt.Now(),
//	}
//
//	// JSON形式のメッセージを生成
//	_json, err := json.Marshal(entity.Message{
//		UserName:  "",
//		Messages:  _msg,
//		CreatedAt: uc.wt.Now(), // カスタムタイムを使用して現在時刻を取得
//	})
//	if err != nil {
//		t.Fatalf("JSONマーシャリングに失敗しました: %v", err)
//	}
//
//	ws.On("ReadMessage").Return(1, []byte(_msg), nil).Twice() // 2回メッセージを読み込む
//	ws.On("ReadMessage").Return(0, nil, io.EOF).Once()        // その後、EOFエラーを返す
//	ws.On("WriteMessage", 1, _json).Return(nil).Maybe()
//	mockRepo.On("CreateMessage", uint32(1)).Return(expectedMessages, nil)
//	mockRedis.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(&entity.IntCmd{}, nil).Maybe()
//
//	roomManager := make(map[uint32]*entity.ChatRoom)
//
//	client := &entity.Client{Ws: wsw}
//	newRoom := entity.ChatRoom{}
//	newRoom.AddClient(client)
//	roomManager[chatRoomID] = &newRoom
//	e := echo.New()
//	req := httptest.NewRequest(echo.GET, "/", nil)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	fmt.Printf("%v", roomManager[chatRoomID])
//	// PubSub メソッドをテスト
//	uc.PubSub(c, wsw, roomManager, chatRoomID, mockRedis)
//
//	// アサーションを確認
//	mockRepo.AssertExpectations(t)
//	mockRedis.AssertExpectations(t)
//	ws.AssertExpectations(t)
//}
