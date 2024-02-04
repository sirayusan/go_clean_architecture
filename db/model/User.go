package model

import "time"

type User struct {
	UserID                    uint      `gorm:"column:user_id;primaryKey;autoIncrement:true;comment:ユーザーID;"`
	LastName                  string    `gorm:"column:last_name;size:100;not null;comment:氏;"`
	FirstName                 string    `gorm:"column:first_name;size:100;not null;comment:名;"`
	HiraganaLastName          string    `gorm:"column:hiragana_last_name;size:100;not null;comment:氏(かな);"`
	HiraganaFirstName         string    `gorm:"column:hiragana_first_name;size:100;not null;comment:名(かな);"`
	Email                     string    `gorm:"column:email;unique;size:256;not null;comment:メールアドレス;"`
	Password                  string    `gorm:"column:password;not null;comment:パスワード;"`
	IsPasswordResetInProgress bool      `gorm:"column:is_password_reset_in_progress;not null;default:false;comment:パスワード再設定中フラグ;"`
	CreatedUserID             uint      `gorm:"column:created_user_id;not null;default:0;comment:登録者ID;"`
	CreatedAt                 time.Time `gorm:"column:created_at;not null;comment:登録日時;"`
	UpdateUserID              uint      `gorm:"column:updated_user_id;not null;default:0;comment:更新者ID;"`
	UpdatedAt                 time.Time `gorm:"column:updated_at;not null;comment:更新日時;"`
}
