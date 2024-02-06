package entity

// User -.
type User struct {
	UserID    uint32 `gorm:"column:user_id"`
	LastName  string `gorm:"column:last_name"`
	FirstName string `gorm:"column:first_name"`
}

type UserListResponse struct {
	UserList []User `json:"user_list"`
}
