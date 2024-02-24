package model

type ChatRoom struct {
	ChatRoomID uint32 `gorm:"column:chat_room_id;primaryKey;autoIncrement:true;comment:チャットルームID;"`
	UserID1    uint32 `gorm:"column:user_id1;not null;comment:ユーザーID1;"`
	UserID2    uint32 `gorm:"column:user_id2;not null;comment:ユーザーID2;"`
}
